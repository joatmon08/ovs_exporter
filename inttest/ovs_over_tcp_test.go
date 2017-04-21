package inttest

import (
	"testing"
	"github.com/joatmon08/ovs_exporter/utils"
	"strconv"
)

const (
	BRIDGE_ID = "br0"
	PORT_ID = "eth0"
	IP = "192.168.128.5"
	OVS_STATE = "openvswitch_up"
)

var (
	bridge_metric = "openvswitch_interfaces_statistics{name=\"" + BRIDGE_ID + "\",stat=\"rx_bytes\"}"
	add_bridge = "ovs-vsctl add-br " + BRIDGE_ID
	set_datapath = "ovs-vsctl set bridge " + BRIDGE_ID + " datapath_type=netdev"
	add_port = "ovs-vsctl add-port " + BRIDGE_ID + " " + PORT_ID
	create_bridge = add_bridge + " && " + set_datapath + " && " + add_port
	configure_bridge = "ifconfig " + BRIDGE_ID + " " + IP
)

func TestOpenvSwitchBridgeWithTraffic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	t.Logf("OVS Bridge %s should have traffic on it", BRIDGE_ID)
	ovsCommand := create_bridge + " && " + configure_bridge
	setup := Setup(t, ovsCommand)
	if err := RetrieveMetrics(setup); err != nil {
		t.Error(err)
	}
	if actual, _ := strconv.Atoi(setup.metrics[OVS_STATE]); actual != 1 {
		t.Errorf("expected ovs state to be 1, actual %d", actual)
	}
	if actual, _ := strconv.Atoi(setup.metrics[bridge_metric]); actual == 0 {
		t.Errorf("expected greater than 0, actual %d", actual)
	}
	t.Logf("metric %s has %s", bridge_metric, setup.metrics[bridge_metric])
	Teardown(setup.ovsContainerID, setup.ovsExporterContainerID)
}

func TestOpenvSwitchBridgeWithoutTraffic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	t.Logf("OVS Bridge %s should not have traffic on it", BRIDGE_ID)
	ovsCommand := create_bridge
	setup := Setup(t, ovsCommand)
	if err := RetrieveMetrics(setup); err != nil {
		t.Error(err)
	}
	if actual, _ := strconv.Atoi(setup.metrics[OVS_STATE]); actual != 1 {
		t.Errorf("expected ovs state to be 1, actual %d", actual)
	}
	if actual, _ := strconv.Atoi(setup.metrics[bridge_metric]); actual != 0 {
		t.Errorf("expected %d, actual %d", 0, actual)
	}
	t.Logf("metric %s is %s", bridge_metric, setup.metrics[bridge_metric])
	Teardown(setup.ovsContainerID, setup.ovsExporterContainerID)
}

func TestOpenvSwitchDown(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	t.Log("OVS state should be down")
	setup := Setup(t, "")
	if err := utils.StopContainer(setup.ovsContainerID); err != nil {
		t.Error(err)
	}
	if err := RetrieveMetrics(setup); err != nil {
		t.Error(err)
	}
	if actual, _ := strconv.Atoi(setup.metrics[OVS_STATE]); actual != 0 {
		t.Errorf("expected ovs state to be 0, actual %d", actual)
	}
	Teardown(setup.ovsContainerID, setup.ovsExporterContainerID)
}