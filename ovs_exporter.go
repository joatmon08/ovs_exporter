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
	ports = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ports_total"),
		"How many Open vSwitch ports on this node.",
		nil, nil,
	)
	bridges_num_ports = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "bridges_ports",
		Help:      "Number of ports attached to bridges",
	},
		[]string{"name"},
	)
	interfaces_stats = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "interfaces_statistics",
		Help:      "Open vSwitch interface statistics",
	},
		[]string{"name", "stat"},
	)
)

type Exporter struct {
	URI               string
	client            *libovsdb.OvsdbClient
	up                *prometheus.Desc
	dbs               *prometheus.Desc
	bridges           *prometheus.Desc
	ports             *prometheus.Desc
	interfaces        *prometheus.Desc
	bridges_num_ports *prometheus.GaugeVec
	interfaces_stats  *prometheus.CounterVec
}

func NewExporter(uri string) (*Exporter, error) {
	return &Exporter{
		URI: uri,
		up: up,
		dbs: dbs,
		client: &libovsdb.OvsdbClient{},
		bridges: bridges,
		ports: ports,
		interfaces: interfaces,
		bridges_num_ports: bridges_num_ports,
		interfaces_stats: interfaces_stats,
	}, nil
}

func (e *Exporter) Describe(ch chan <- *prometheus.Desc) {
	ch <- up
	ch <- dbs
	ch <- bridges
	ch <- ports
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

func (e *Exporter) connect() error {
	network, err := openvswitch.GenerateNetworkAndHealthCheck(e.URI)
	if err != nil {
		return err
	}
	e.client, err = libovsdb.ConnectUsingProtocol(network, e.URI)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"event": "cannot connect to ovsdb",
		}).Error(err)
	}
	return err
}

func (e *Exporter) Collect(ch chan <- prometheus.Metric) {
	if err := e.connect(); err != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
		return
	}
	defer e.client.Disconnect()
	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)

	databases, err := openvswitch.GetDatabases(e.client)
	if err != nil {
		logrus.Error(err)
	}
	ch <- prometheus.MustNewConstMetric(
		dbs, prometheus.GaugeValue, float64(len(databases)),
	)
	total_bridges := openvswitch.GetRowsFromTable(e.client, "Bridge")
	ch <- prometheus.MustNewConstMetric(
		bridges, prometheus.GaugeValue, float64(len(total_bridges)),
	)
	total_ports := openvswitch.GetRowsFromTable(e.client, "Port")
	ch <- prometheus.MustNewConstMetric(
		ports, prometheus.GaugeValue, float64(len(total_ports)),
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


func init() {
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logrus.SetFormatter(formatter)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	var (
		uri = flag.String("uri", "/var/run/openvswitch/db.sock", "URI to connect to Open vSwitch")
		listenAddress = flag.String("listen-port", ":9177", "Address to listen on for web interface and telemetry.")
		metricsPath = flag.String("metrics-path", "/metrics", "Path under which to expose metrics.")
	)
	flag.Parse()

	exporter, err := NewExporter(*uri)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"uri": *uri,
			"event": "starting exporter",
		}).Fatal(err)
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

	logrus.WithFields(logrus.Fields{
		"port": *listenAddress,
		"path": *metricsPath,
		"event": "listening",
	}).Info("prometheus started")

	logrus.WithFields(logrus.Fields{
		"port": *listenAddress,
		"path": *metricsPath,
		"event": "web server error",
	}).Fatal(http.ListenAndServe(*listenAddress, nil))
}