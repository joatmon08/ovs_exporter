package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/socketplane/libovsdb"
	"flag"
	"github.com/Sirupsen/logrus"
	"net/http"
	"github.com/joatmon08/ovs_exporter/openvswitch"
)

const (
	namespace = "openvswitch" // For Prometheus metrics.
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of Open vSwitch successful.",
		nil, nil,
	)
	dbs = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "dbs_total"),
		"How many Open vSwitch dbs on this node.",
		nil, nil,
	)
	bridges = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "bridges_total"),
		"How many Open vSwitch bridges on this node.",
		nil, nil,
	)
	interfaces = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "interfaces_total"),
		"How many Open vSwitch interfaces on this node.",
		nil, nil,
	)
)

type Exporter struct {
	URI        string
	client     *libovsdb.OvsdbClient
	up         *prometheus.Desc
	dbs        *prometheus.Desc
	bridges    *prometheus.Desc
	interfaces *prometheus.Desc
}

func NewExporter(uri string) (*Exporter, error) {
	client, err := libovsdb.ConnectWithUnixSocket(uri)
	if err != nil {
		return nil, err
	}
	return &Exporter{
		URI: uri,
		up: up,
		dbs: dbs,
		client: client,
		bridges: bridges,
		interfaces: interfaces,
	}, nil
}

func (e *Exporter) Describe(ch chan <- *prometheus.Desc) {
	ch <- up
	ch <- dbs
	ch <- bridges
	ch <- interfaces
}

func (e *Exporter) Collect(ch chan <- prometheus.Metric) {
	databases, err := openvswitch.CheckHealth(e.client)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
		logrus.Errorf("Query error is %v", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)
	ch <- prometheus.MustNewConstMetric(
		dbs, prometheus.GaugeValue, float64(len(databases)),
	)
	total_bridges := openvswitch.GetTotalFromTable(e.client, "Bridge")
	ch <- prometheus.MustNewConstMetric(
		bridges, prometheus.GaugeValue, float64(len(total_bridges)),
	)
	total_interfaces := openvswitch.GetTotalFromTable(e.client, "Interface")
	ch <- prometheus.MustNewConstMetric(
		interfaces, prometheus.GaugeValue, float64(len(total_interfaces)),
	)
}

func main() {
	var (
		uri = flag.String("uri", "/var/run/openvswitch/db.sock", "URI to connect to Open vSwitch")
		listenAddress = flag.String("listen-address", ":9107", "Address to listen on for web interface and telemetry.")
		metricsPath = flag.String("telemetry-path", "/metrics", "Path under which to expose metrics.")
	)
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)
	exporter, err := NewExporter(*uri)
	if err != nil {
		logrus.Fatalln(err)
	}
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Open vSwitch Exporter</title></head>
             <body>
             <h1>Open vSwitch Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	logrus.Infof("Listening on %s", *listenAddress)
	logrus.Fatal(http.ListenAndServe(*listenAddress, nil))
}