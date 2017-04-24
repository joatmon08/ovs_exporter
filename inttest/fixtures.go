package inttest

import (
	"github.com/Sirupsen/logrus"
	"github.com/joatmon08/ovs_exporter/utils"
	"time"
	"errors"
	"testing"
	"github.com/docker/docker/client"
	"io/ioutil"
	"context"
	"github.com/docker/docker/api/types"
)

const (
	TCP = "tcp"
	UNIX = "unix"
	CREATE_WAIT_TIME = 2 * time.Second
	EXEC_WAIT_TIME = 5 * time.Second
	INTTEST_NETWORK = "ovs_exporter_inttest_network"
	INTTEST_NETWORK_CIDR = "172.19.0.0"
	OPENVSWITCH_IP = "172.19.0.2"
	OPENVSWITCH_PORT = ":6640"
	EXPORTER_PORT = ":9177"
	OVS_CONTAINER_IMAGE = "socketplane/openvswitch:latest"
	OPENVSWITCH_JSON = "openvswitch"
	EXPORTER_JSON = "ovs_exporter"
	BRIDGE_ID = "br0"
	PORT_ID = "eth0"
	IP = "192.168.128.5"
	OVS_STATE = "openvswitch_up"
	OVS_INTERFACES = "openvswitch_interfaces_total"
	OVS_PORTS = "openvswitch_ports_total"
)

var (
	BridgeMetric = "openvswitch_interfaces_statistics{name=\"" + BRIDGE_ID + "\",stat=\"rx_bytes\"}"
	AddBridge = "ovs-vsctl add-br " + BRIDGE_ID
	SetDatapath = "ovs-vsctl set bridge " + BRIDGE_ID + " datapath_type=netdev"
	AddPort = "ovs-vsctl add-port " + BRIDGE_ID + " " + PORT_ID
	CreateBridge = AddBridge + " && " + SetDatapath + " && " + AddPort
	ConfigureBridge = "ifconfig " + BRIDGE_ID + " " + IP
	OVSUNIXCommand = "app -listen-port " + EXPORTER_PORT
	OVSTCPCommand = OVSUNIXCommand + " -uri " + OPENVSWITCH_IP + OPENVSWITCH_PORT
)

type testSetupObject struct {
	ovsConnectionMode      string
	containerExecCmd       string
	ovsContainerID         string
	ovsExporterContainerID string
	networkID              string
	metrics                map[string]string
}

func createContainers(exporterCmd string) (ovsContainerID string, ovsExporterContainerID string) {
	var err error
	//err := pullOVSImage()
	//if err != nil {
	//	panic(err)
	//}
	ovsArgs := &utils.OptionalContainerArgs{
		Network: INTTEST_NETWORK,
	}
	if exporterCmd == OVSUNIXCommand {
		ovsArgs.HostBinds = []string{
			"/tmp/openvswitch:/usr/local/var/run/openvswitch",
		}
	}
	ovsContainerID, err = utils.CreateContainer(OPENVSWITCH_JSON, ovsArgs)
	if err != nil {
		panic(err)
	}
	err = utils.StartContainer(ovsContainerID)
	if err != nil {
		panic(err)
	}
	logrus.Debugf("created ovs container %s", ovsContainerID)
	exporterArgs := &utils.OptionalContainerArgs{
		Network: INTTEST_NETWORK,
		Cmd: exporterCmd,
	}
	if exporterCmd == OVSUNIXCommand {
		exporterArgs.HostBinds = []string{
			"/tmp/openvswitch:/var/run/openvswitch",
		}
	}
	ovsExporterContainerID, err = utils.CreateContainer(EXPORTER_JSON, exporterArgs)
	if err != nil {
		panic(err)
	}
	err = utils.StartContainer(ovsExporterContainerID)
	if err != nil {
		panic(err)
	}
	logrus.Debugf("created ovs exporter container %s", ovsExporterContainerID)
	time.Sleep(CREATE_WAIT_TIME)
	return ovsContainerID, ovsExporterContainerID
}

func pullOVSImage() error {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		logrus.Errorf("docker client was not created")
		panic(err)
	}
	reader, err := dockerClient.ImagePull(context.Background(), OVS_CONTAINER_IMAGE, types.ImagePullOptions{})
	_, err = ioutil.ReadAll(reader)
	reader.Close()
	if err != nil {
		logrus.Errorf("ovs image was not pulled")
		return err
	}
	return nil
}

func RetrieveMetrics(testSetup *testSetupObject) (error) {
	ovsClient := utils.NewOVSExporterClient("http://localhost:9177")
	metrics, err := ovsClient.GetExporterMetrics()
	if err != nil {
		return err
	}
	if len(metrics) == 0 {
		return errors.New("no metrics, metrics map is empty")
	}
	testSetup.metrics = metrics
	return nil
}

func Setup(t *testing.T, testSetup *testSetupObject) (*testSetupObject) {
	var ovsEntrypoint string
	networkID, err := utils.CreateNetwork(INTTEST_NETWORK, INTTEST_NETWORK_CIDR)
	if err != nil {
		t.Error(err)
	}
	testSetup.networkID = networkID
	switch connection := testSetup.ovsConnectionMode; connection {
	case TCP:
		ovsEntrypoint = OVSTCPCommand
	case UNIX:
		ovsEntrypoint = OVSUNIXCommand
	default:
		t.Error("Specify unix or tcp mode for OVS container")
	}
	ovs, exporter := createContainers(ovsEntrypoint)
	testSetup.ovsExporterContainerID = exporter
	testSetup.ovsContainerID = ovs
	if testSetup.containerExecCmd == "" {
		return testSetup
	}
	commands := []string{utils.SHELL, utils.COMMAND_OPTION, testSetup.containerExecCmd}
	if err := utils.ExecuteContainer(ovs, commands); err != nil {
		t.Error(err)
	}
	time.Sleep(EXEC_WAIT_TIME)
	return testSetup
}

func Teardown(ovsContainerID string, ovsExporterContainerID string, networkID string) {
	if err := utils.DeleteContainer(ovsExporterContainerID); err != nil {
		logrus.Error(err)
	}
	if err := utils.DeleteContainer(ovsContainerID); err != nil {
		logrus.Error(err)
	}
	if err := utils.DeleteNetwork(networkID); err != nil {
		logrus.Error(err)
	}
}


