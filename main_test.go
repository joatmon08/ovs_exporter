package main

import (
	"os"
	"testing"
	"github.com/docker/docker/client"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"io/ioutil"
	"github.com/docker/docker/api/types/network"
	"github.com/joatmon08/ovs_exporter/utils"
	"flag"
	"time"
)

const (
	OVS_CONTAINER_IMAGE = "socketplane/openvswitch:latest"
)

var (
	inttests = flag.Bool("intTests", false, "run integration tests")
	setupwait = 2 * time.Second
)

type ContainerConfig struct {
	HostConfig    *container.HostConfig
	Config        *container.Config
	NetworkConfig *network.NetworkingConfig
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

func Setup() (ovsContainerID string, ovsExporterContainerID string) {
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
	logrus.Infof("created ovs container %s", ovsContainerID)
	ovsExporterContainerID, err = utils.CreateContainer("ovs_exporter")
	if err != nil {
		panic(err)
	}
	err = utils.StartContainer(ovsExporterContainerID)
	if err != nil {
		panic(err)
	}
	logrus.Infof("created ovs exporter container %s", ovsExporterContainerID)
	time.Sleep(setupwait)
	return ovsContainerID, ovsExporterContainerID
}

func Shutdown(ovsContainerID string, ovsExporterContainerID string) {
	err := utils.DeleteContainer(ovsExporterContainerID)
	if err != nil {
		logrus.Error(err)
	}
	err = utils.DeleteContainer(ovsContainerID)
	if err != nil {
		logrus.Error(err)
	}
}

func CleanupTesting() {
	ids, err := utils.GetAllContainerIDs()
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	logrus.Info("deleting all existing containers")
	for _, id := range ids {
		err := utils.DeleteContainer(id)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	//CleanupTesting()
	result := m.Run()
	os.Exit(result)
}
