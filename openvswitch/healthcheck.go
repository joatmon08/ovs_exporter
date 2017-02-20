package openvswitch

import (
	"strings"
	"net/rpc"
	"github.com/Sirupsen/logrus"
)

func GenerateNetworkAndHealthCheck(uri string) (string, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	dbs := []string{}
	network := "unix"
	if strings.Contains(uri, ":") {
		network = "tcp"
	}
	logrus.Info("%s, %s", network, uri)
	client, err := rpc.Dial(network, uri)
	if err != nil {
		return network, err
	}
	//call remote procedure with args
	err = client.Call("list_dbs", nil, &dbs)
	if err != nil {
		return network, err
	}
	return network, nil
}