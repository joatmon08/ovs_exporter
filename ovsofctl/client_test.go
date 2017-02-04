package ovsofctl

import (
	"testing"
	"github.com/joatmon08/ovs_exporter/command"
	"github.com/Sirupsen/logrus"
)

func TestConnectLocal(t *testing.T) {
	logrus.SetLevel(logrus.InfoLevel)
	output, err := ConnectLocal(command.TestRunner{FailTest:false})
	if err != nil {
		t.Error("Should grab OVS version information")
	}
	if output.ofVersion != "2.5.0" {
		t.Errorf("Should contain 2.5.0, got %s", output.ofVersion)
	}
}

func TestConnectLocalFail(t *testing.T) {
	_, err := ConnectLocal(command.TestRunner{FailTest:true})
	if err == nil {
		t.Error("Should not find ovs-ofctl")
	}
}
