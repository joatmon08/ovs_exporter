package ovsofctl

import (
	"github.com/joatmon08/ovs_exporter/utils"
	"strings"
	"strconv"
	"regexp"
	"errors"
	"github.com/Sirupsen/logrus"
)

type OvsPort struct {
	ID       string
	Name     string
	Addr     string
	Current  string
	Speed    int
	MaxSpeed int
}

func getSpeeds(data interface{}) (int, int, error) {
	var err error
	now := 0
	max := 0
	speedInfo, ok := data.(string)
	if !ok {
		return 0, 0, errors.New("speed is not string type")
	}
	speeds := regexp.MustCompile(`([0-9].*?)Mbps`).FindAllStringSubmatch(speedInfo, -1)
	if now, err = strconv.Atoi(strings.TrimSpace(speeds[0][1])); err != nil {
		return now, max, err
	}
	if max, err = strconv.Atoi(strings.TrimSpace(speeds[1][1])); err != nil {
		return now, max, err
	}
	return now, max, nil
}

func (p *OvsPort) Fill(m map[string]interface{}) error {
	for k, v := range m {
		if strings.ToLower(k) == "speed" {
			now, max, err := getSpeeds(v)
			if err != nil {
				return err
			}
			if err = utils.SetField(p, k, now); err != nil {
				return err
			}
			if err := utils.SetField(p, "Max" + k, max); err != nil {
				return err
			}
		} else {
			err := utils.SetField(p, k, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetPort(portInfo string) (*OvsPort, error) {
	header := strings.Split(portInfo, "\n")[0]
	id := regexp.MustCompile(`[A-Z0-9]+`).FindStringSubmatch(header)[0]
	name := regexp.MustCompile(`\((.*?)\)`).FindStringSubmatch(header)[1]
	logrus.Infof("bridge: id %s, name %s", id, name)
	port := &OvsPort{
		ID: id,
		Name: name,
	}
	portInterface := utils.MapStringToInterface(*port, portInfo)
	if err := port.Fill(portInterface); err != nil {
		return port, err
	}
	return port, nil
}
