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

const ENDPOINT = "unix:///var/run/docker.sock"

func ReadInContainerTemplate(name string) (dockerclient.CreateContainerOptions, error) {
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

func DeleteContainer(id string) error {
	client, err := dockerclient.NewClient(ENDPOINT)
	if err != nil {
		return err
	}
	options := dockerclient.RemoveContainerOptions{
		ID: id,
		RemoveVolumes: true,
		Force: true,
	}
	logrus.Debugf("Removing container %s", id)
	if err := client.RemoveContainer(options); err != nil {
		return err
	}
	return nil
}

func StartContainer(id string) error {
	client, err := dockerclient.NewClient(ENDPOINT)
	if err != nil {
		return err
	}
	logrus.Debugf("Starting container %s", id)
	if err := client.StartContainer(id, nil); err != nil {
		return err
	}
	return nil
}

func StopContainer(id string) error {
	client, err := dockerclient.NewClient(ENDPOINT)
	if err != nil {
		return err
	}
	logrus.Debugf("Stopping container %s", id)
	if err := client.StopContainer(id, 0); err != nil {
		return err
	}
	return nil
}

func CreateContainer(description string) (string, error) {
	containerOptions, err := ReadInContainerTemplate(description)
	if err != nil {
		return "", err
	}
	client, err := dockerclient.NewClient(ENDPOINT)
	if err != nil {
		return "", err
	}
	container, err := client.CreateContainer(containerOptions)
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

func GetAllContainerIDs() ([]string, error) {
	var containerIDs []string
	client, err := dockerclient.NewClient(ENDPOINT)
	if err != nil {
		return containerIDs, err
	}
	containers, err := client.ListContainers(dockerclient.ListContainersOptions{All: true})
	if err != nil {
		return containerIDs, err
	}
	for _, container := range containers {
		containerIDs = append(containerIDs, container.ID)
	}
	return containerIDs, nil
}