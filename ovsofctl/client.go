package ovsofctl

import (
	"sync"
	"errors"
	"strings"
	"github.com/joatmon08/ovs_exporter/command"
)

const ovsofctl_cmd = "ovs-ofctl"

type OvsOFClient struct {
	runner    *command.Runner
	ofVersion string
	mutex     *sync.Mutex
}

func newOvsOFClient(commander *command.Runner, ofVersion string) *OvsOFClient {
	of := &OvsOFClient{
		runner:    commander,
		ofVersion: ofVersion,
		mutex: &sync.Mutex{},
	}
	return of
}

func grabLine(index int, out []byte) (string, error) {
	lines := strings.Split(string(out), "\n")
	if index > len(lines) {
		return "", errors.New("Invalid index")
	}
	return lines[index], nil
}

// Connect and returns an OvsOFClient
func ConnectLocal(c command.Runner) (*OvsOFClient, error) {
	output, err := c.Run(ovsofctl_cmd, "--version")
	if err != nil {
		return nil, err
	}
	versionLine, err := grabLine(0, output)
	if err != nil {
		return nil, err
	}
	version := strings.TrimSpace(strings.Split(versionLine, ") ")[1])
	of := newOvsOFClient(&c, version)
	return of, nil
}