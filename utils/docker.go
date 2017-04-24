package utils

import (
	"io/ioutil"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	dockerclient "github.com/fsouza/go-dockerclient"
	"errors"
	"strconv"
	"bytes"
)

const (
	ENDPOINT = "unix:///var/run/docker.sock"
	SHELL = "/bin/sh"
	COMMAND_OPTION = "-c"
)

type OptionalContainerArgs struct {
	Network   string
	Cmd       string
	HostBinds []string
}

func readInContainerTemplate(name string) (dockerclient.CreateContainerOptions, error) {
	var container dockerclient.CreateContainerOptions
	file, err := ioutil.ReadFile("testdata/" + name + ".json")
	if err != nil {
		return container, err
	}

	if err := json.Unmarshal(file, &container); err != nil {
		return container, err
	}
	return container, nil
}

func readInNetworkTemplate(name string) (dockerclient.CreateNetworkOptions, error) {
	var network dockerclient.CreateNetworkOptions
	file, err := ioutil.ReadFile("testdata/" + name + ".json")
	if err != nil {
		return network, err
	}
	if err := json.Unmarshal(file, &network); err != nil {
		return network, err
	}
	return network, nil
}

func DeleteContainer(id string) error {
	containerClient, err := dockerclient.NewClient(ENDPOINT)
	if err != nil {
		return err
	}
	options := dockerclient.RemoveContainerOptions{
		ID: id,
		RemoveVolumes: true,
		Force: true,
	}
	logrus.Debugf("Removing container %s", id)
	if err := containerClient.RemoveContainer(options); err != nil {
		return err
	}
	return nil
}

func StartContainer(id string) error {
	containerClient, err := dockerclient.NewClient(ENDPOINT)
	if err != nil {
		return err
	}
	logrus.Debugf("Starting container %s", id)
	if err := containerClient.StartContainer(id, nil); err != nil {
		return err
	}
	return nil
}

func StopContainer(id string) error {
	containerClient, err := dockerclient.NewClient(ENDPOINT)
	if err != nil {
		return err
	}
	logrus.Debugf("Stopping container %s", id)
	if err := containerClient.StopContainer(id, 0); err != nil {
		return err
	}
	return nil
}

func CreateContainer(description string, optionalArgs *OptionalContainerArgs) (string, error) {
	containerOptions, err := readInContainerTemplate(description)
	if err != nil {
		return "", err
	}
	if optionalArgs.Network != "" {
		containerOptions.HostConfig.NetworkMode = optionalArgs.Network
	}
	if optionalArgs.Cmd != "" {
		containerOptions.Config.Cmd = []string{
			SHELL,
			COMMAND_OPTION,
			optionalArgs.Cmd,
		}
	}
	if optionalArgs.HostBinds != nil {
		containerOptions.HostConfig.Binds = optionalArgs.HostBinds
	}
	logrus.Infof("%s, %s", description, optionalArgs)
	containerClient, err := dockerclient.NewClient(ENDPOINT)
	if err != nil {
		return "", err
	}
	container, err := containerClient.CreateContainer(containerOptions)
	if err != nil {
		return "", err
	}
	return container.ID, nil
}

func ExecuteContainer(containerID string, commands []string) (error) {
	client, err := dockerclient.NewClient(ENDPOINT)
	if err != nil {
		return err
	}
	execInstance, err := client.CreateExec(dockerclient.CreateExecOptions{
		Container: containerID,
		AttachStdin: false,
		AttachStdout: true,
		AttachStderr: true,
		Tty: false,
		Cmd: commands,
	})
	logrus.Debugf("container %s, exec instance %s", containerID, execInstance.ID)
	if err != nil {
		return err
	}
	var stdout bytes.Buffer
	err = client.StartExec(execInstance.ID, dockerclient.StartExecOptions{
		OutputStream: &stdout,
		Detach: false,
		Tty: false,
		RawTerminal: true,
	})
	if err := client.StartExec(execInstance.ID, dockerclient.StartExecOptions{
		OutputStream: &stdout,
		Detach: false,
		Tty: false,
		RawTerminal: true,
	}); err != nil {
		return err
	}
	execResult, err := client.InspectExec(execInstance.ID)
	if err != nil {
		return err
	}
	if execResult.ExitCode != 0 {
		logrus.Errorf("docker exec failed with exit code %d, %s", execResult.ExitCode, stdout.String())
		return errors.New(strconv.Itoa(execResult.ExitCode))
	}
	return nil
}

func CreateNetwork(description string, cidr string) (string, error) {
	options, err := readInNetworkTemplate(description)
	if err != nil {
		return "", err
	}
	if cidr != "" {
		options.IPAM.Config = []dockerclient.IPAMConfig{
			{Subnet: cidr + "/16" },
		}
	}
	client, err := dockerclient.NewClient(ENDPOINT)
	if err != nil {
		return "", err
	}
	logrus.Debugf("Creating network %s", options.Name)
	network, err := client.CreateNetwork(options)
	if err != nil {
		return "", err
	}
	return network.ID, nil
}

func DeleteNetwork(id string) error {
	client, err := dockerclient.NewClient(ENDPOINT)
	if err != nil {
		return err
	}
	logrus.Debugf("Removing network %s", id)
	err = client.RemoveNetwork(id)
	if err != nil {
		return err
	}
	return nil
}