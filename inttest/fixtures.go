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
	CREATE_WAIT_TIME = 2 * time.Second
	EXEC_WAIT_TIME = 5 * time.Second
	SHELL = "/bin/sh"
	COMMAND = "-c"
	OVS_CONTAINER_IMAGE = "socketplane/openvswitch:latest"
)

type testSetupObject struct {
	ovsContainerID         string
	ovsExporterContainerID string
	metrics                map[string]string
}

func createContainers() (ovsContainerID string, ovsExporterContainerID string) {
	var err error
	//err := pullOVSImage()
	//if err != nil {
	//	panic(err)
	//}
	ovsContainerID, err = utils.CreateContainer("openvswitch")
	if err != nil {
		panic(err)
	}
	err = utils.StartContainer(ovsContainerID)
	if err != nil {
		panic(err)
	}
	logrus.Debugf("created ovs container %s", ovsContainerID)
	ovsExporterContainerID, err = utils.CreateContainer("ovs_exporter")
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

func Setup(t *testing.T, cmd string) (*testSetupObject) {
	ovs, exporter := createContainers()
	testSetup := testSetupObject{
		ovsContainerID: ovs,
		ovsExporterContainerID: exporter,
	}
	if cmd == "" {
		return &testSetup
	}
	commands := []string{SHELL, COMMAND, cmd}
	if err := utils.ExecuteContainer(ovs, commands); err != nil {
		t.Error(err)
	}
	time.Sleep(EXEC_WAIT_TIME)
	return &testSetup
}

func Teardown(ovsContainerID string, ovsExporterContainerID string) {
	if err := utils.DeleteContainer(ovsExporterContainerID); err != nil {
		logrus.Error(err)
	}
	if err := utils.DeleteContainer(ovsContainerID); err != nil {
		logrus.Error(err)
	}
}


