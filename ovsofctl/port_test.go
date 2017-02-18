package ovsofctl

import (
	"testing"
	"github.com/joatmon08/ovs_exporter/utils"
)



func TestGetPort(t *testing.T) {
	testPortInfo, err := utils.ReadTestData("port")
	if err != nil {
		t.Error(err)
	}
	port, err := GetPort(testPortInfo)
	if err != nil {
		t.Error(err)
	}
	expectedPort := OvsPort{
		ID: "1",
		Name: "87be70c563b24_l",
		Addr: "d2:e5:b4:c9:e9:86",
		Current: "10GB-FD COPPER",
		Speed: 10000,
		MaxSpeed: 0,
	}
	if expectedPort != *port {
		t.Errorf("Expected %v, got %v", expectedPort, *port)
	}
}