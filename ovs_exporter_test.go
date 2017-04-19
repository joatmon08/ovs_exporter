package main

import (
	"testing"
	"github.com/joatmon08/ovs_exporter/utils"
	"strconv"
	"time"
	"errors"
)

const (
	WAITTIME = 5 * time.Second
	SHELL = "/bin/sh"
	COMMAND = "-c"
	BRIDGE_ID = "br0"
	PORT_ID = "eth0"
	IP = "192.168.128.5"
	METRIC = "openvswitch_interfaces_statistics{name=\"" + BRIDGE_ID + "\",stat=\"rx_bytes\"}"
	ADD_BR = "ovs-vsctl add-br " + BRIDGE_ID
	SET_DATAPATH = "ovs-vsctl set bridge " + BRIDGE_ID + " datapath_type=netdev"
	ADD_PORT = "ovs-vsctl add-port " + BRIDGE_ID + " " + PORT_ID
	CONFIG_BR = "ifconfig " + BRIDGE_ID + " " + IP
)

func setupAndGetMetrics(cmd string) (string, string, map[string]string, error) {
	ovs, exporter := Setup()
	commands := []string{ SHELL, COMMAND, cmd }
	err := utils.ExecuteContainer(ovs, commands)
	if err != nil {
		return ovs, exporter, nil, err
	}
	time.Sleep(WAITTIME)
	ovsClient := utils.NewOVSExporterClient("http://localhost:9177")
	metrics, err := ovsClient.GetExporterMetrics()
	if err != nil {
		return ovs, exporter, metrics, err
	}
	if len(metrics) == 0 {
		return ovs, exporter, metrics, errors.New("no metrics, metrics map is empty")
	}
	if err != nil {
		return ovs, exporter, metrics, err
	}
	return ovs, exporter, metrics, nil
}

func TestOpenvSwitchBridgeWithTraffic(t *testing.T) {
	ovs, exporter, metrics, err := setupAndGetMetrics(ADD_BR + " && " + SET_DATAPATH + " && " +
		ADD_PORT + " && " + CONFIG_BR)
	if err != nil {
		t.Error(err)
	}
	if actual, _ := strconv.Atoi(metrics[METRIC]); actual == 0 {
		t.Errorf("expected greater than 0, actual %d", actual)
	}
	t.Logf("metric %s has %s", METRIC, metrics[METRIC])
	Shutdown(ovs, exporter)
}

func TestOpenvSwitchBridgeWithoutTraffic(t *testing.T) {
	ovs, exporter, metrics, err := setupAndGetMetrics(ADD_BR + " && " + SET_DATAPATH + " && " + ADD_PORT)
	if err != nil {
		t.Error(err)
	}
	if actual, _ := strconv.Atoi(metrics[METRIC]); actual != 0 {
		t.Errorf("expected %d, actual %d", 0, actual)
	}
	t.Logf("metric %s has %s", METRIC, metrics[METRIC])
	Shutdown(ovs, exporter)
}