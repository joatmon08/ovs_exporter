package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/socketplane/libovsdb"
	"flag"
	"github.com/Sirupsen/logrus"
	"net/http"
	"github.com/joatmon08/ovs_exporter/openvswitch"
	"strings"
	"strconv"
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
	bridges_num_ports = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "bridges_num_ports",
		Help:      "Number of ports attached to bridges",
	},
		[]string{"name"},
	)
	interfaces_stats = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "interfaces",
		Help:      "Open vSwitch interface statistics",
	},
		[]string{"name", "stat"},
	)
)

type Exporter struct {
	URI              string
	client           *libovsdb.OvsdbClient
	up               *prometheus.Desc
	dbs              *prometheus.Desc
	bridges          *prometheus.Desc
	interfaces       *prometheus.Desc
	bridges_num_ports    *prometheus.GaugeVec
	interfaces_stats *prometheus.CounterVec
}

func NewExporter(uri string) (*Exporter, error) {
	var client *libovsdb.OvsdbClient
	var err error
	if strings.Contains(uri, ":") {
		ip := strings.Split(uri, ":")
		port, err := strconv.Atoi(ip[1])
		if err != nil {
			port = 6640
		}
		client, err = libovsdb.Connect(ip[0], port)
	} else {
		client, err = libovsdb.ConnectWithUnixSocket(uri)
	}
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
		bridges_num_ports: bridges_num_ports,
		interfaces_stats: interfaces_stats,
	}, nil
}

func (e *Exporter) Describe(ch chan <- *prometheus.Desc) {
	ch <- up
	ch <- dbs
	ch <- bridges
	ch <- interfaces
	e.bridges_num_ports.Describe(ch)
	e.interfaces_stats.Describe(ch)
}

func (e *Exporter) collectPortsForBridges(rows []map[string]interface{}) {
	e.bridges_num_ports.Reset()
	bridges, err := openvswitch.ParsePortsFromBridges(rows)
	if err != nil {
		return
	}
	for _, bridge := range bridges {
		e.bridges_num_ports.WithLabelValues(bridge.Name).Set(float64(len(bridge.Ports)))
	}
}

func (e *Exporter) collectInterfacesStats(rows []map[string]interface{}) {
	e.interfaces_stats.Reset()
	interfaces, err := openvswitch.ParseStatisticsFromInterfaces(rows)
	if err != nil {
		return
	}
	for _, iface := range interfaces {
		for stat_name, num := range iface.Statistics {
			e.interfaces_stats.WithLabelValues(iface.Name, stat_name).Add(num)
			e.interfaces_stats.WithLabelValues(iface.Name, stat_name).Add(num)
		}
	}
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
	total_bridges := openvswitch.GetRowsFromTable(e.client, "Bridge")
	ch <- prometheus.MustNewConstMetric(
		bridges, prometheus.GaugeValue, float64(len(total_bridges)),
	)
	total_interfaces := openvswitch.GetRowsFromTable(e.client, "Interface")
	ch <- prometheus.MustNewConstMetric(
		interfaces, prometheus.GaugeValue, float64(len(total_interfaces)),
	)
	e.collectPortsForBridges(total_bridges)
	e.bridges_num_ports.Collect(ch)
	e.collectInterfacesStats(total_interfaces)
	e.interfaces_stats.Collect(ch)
}

func main() {
	var (
		uri = flag.String("uri", "/var/run/openvswitch/db.sock", "URI to connect to Open vSwitch")
		listenAddress = flag.String("listen-port", ":9107", "Address to listen on for web interface and telemetry.")
		metricsPath = flag.String("metrics-path", "/metrics", "Path under which to expose metrics.")
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