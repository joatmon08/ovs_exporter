package openvswitch

import (
	"strings"
	"net"
	"github.com/Sirupsen/logrus"
)

func GenerateNetworkAndHealthCheck(uri string) (string, error) {
	var err error
	network := "unix"
	if strings.Contains(uri, ":") {
		network = "tcp"
	}

	conn, err := net.Dial(network, uri)
	if err != nil {
		logrus.Errorf("Connection error:%s", err)
		return network, err
	}
	defer conn.Close()
	return network, nil
}