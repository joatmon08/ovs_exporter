package inttest

import (
	"testing"
	"github.com/joatmon08/ovs_exporter/utils"
	"strconv"
)

func TestTCPOpenvSwitchBridgeWithTraffic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	t.Logf("OVS Bridge %s should have traffic on it", BRIDGE_ID)
	ovsCommand := CreateBridge + " && " + ConfigureBridge
	setup := Setup(t, ovsCommand)
	if err := RetrieveMetrics(setup); err != nil {
		t.Error(err)
	}
	if actual, _ := strconv.Atoi(setup.metrics[OVS_STATE]); actual != 1 {
		t.Errorf("expected ovs state to be 1, actual %d", actual)
	}
	if actual, _ := strconv.Atoi(setup.metrics[OVS_INTERFACES]); actual != 2 {
		t.Errorf("expected ovs interfaces total to be 2, actual %d", actual)
	}
	if actual, _ := strconv.Atoi(setup.metrics[OVS_PORTS]); actual != 2 {
		t.Errorf("expected ovs ports total to be 2, actual %d", actual)
	}
	if actual, _ := strconv.Atoi(setup.metrics[BridgeMetric]); actual == 0 {
		t.Errorf("expected greater than 0, actual %d", actual)
	}
	t.Logf("metric %s has %s", BridgeMetric, setup.metrics[BridgeMetric])
	Teardown(setup.ovsContainerID, setup.ovsExporterContainerID)
}

func TestTCPOpenvSwitchBridgeWithoutTraffic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	t.Logf("OVS Bridge %s should not have traffic on it", BRIDGE_ID)
	ovsCommand := CreateBridge
	setup := Setup(t, ovsCommand)
	if err := RetrieveMetrics(setup); err != nil {
		t.Error(err)
	}
	if actual, _ := strconv.Atoi(setup.metrics[OVS_STATE]); actual != 1 {
		t.Errorf("expected ovs state to be 1, actual %d", actual)
	}
	if actual, _ := strconv.Atoi(setup.metrics[OVS_INTERFACES]); actual != 2 {
		t.Errorf("expected ovs interfaces total to be 2, actual %d", actual)
	}
	if actual, _ := strconv.Atoi(setup.metrics[OVS_PORTS]); actual != 2 {
		t.Errorf("expected ovs ports total to be 2, actual %d", actual)
	}
	if actual, _ := strconv.Atoi(setup.metrics[BridgeMetric]); actual != 0 {
		t.Errorf("expected %d, actual %d", 0, actual)
	}
	t.Logf("metric %s is %s", BridgeMetric, setup.metrics[BridgeMetric])
	Teardown(setup.ovsContainerID, setup.ovsExporterContainerID)
}

func TestTCPOpenvSwitchWithoutBridge(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	t.Log("OVS should not have a bridge or interfaces")
	setup := Setup(t, "")
	if err := RetrieveMetrics(setup); err != nil {
		t.Error(err)
	}
	if actual, _ := strconv.Atoi(setup.metrics[OVS_STATE]); actual != 1 {
		t.Errorf("expected ovs state to be 1, actual %d", actual)
	}
	if actual, _ := strconv.Atoi(setup.metrics[OVS_INTERFACES]); actual != 0 {
		t.Errorf("expected ovs interfaces total to be 0, actual %d", actual)
	}
	if actual, _ := strconv.Atoi(setup.metrics[OVS_PORTS]); actual != 0 {
		t.Errorf("expected ovs ports total to be 0, actual %d", actual)
	}
	Teardown(setup.ovsContainerID, setup.ovsExporterContainerID)
}

func TestTCPOpenvSwitchDown(t *testing.T) {
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