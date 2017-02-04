package command

import (
	"errors"
	"strings"
)

type TestRunner struct {
	FailTest bool
}

var ovs_version = []byte("ovs-ofctl (Open vSwitch) 2.5.0\n" +
	"Compiled Mar 18 2016 15:00:11\n" +
	"OpenFlow versions 0x1:0x4")

func (r TestRunner) Run(cmd string, args ...string) ([]byte, error) {
	if strings.Contains(args[0], "version") {
		if r.FailTest {
			return nil, errors.New("bash: ovs-ofctl: command not found")
		} else {
			return ovs_version, nil
		}

	}
	return []byte{}, nil
}