package ovsofctl

import (
	"github.com/joatmon08/ovs_exporter/utils"
	"strings"
	"github.com/Sirupsen/logrus"
	"regexp"
	"strconv"
)

type OpenvSwitch struct {
	N_tables     int
	N_buffers    int
	Capabilities []string
	Actions      []string
	Ports        map[string]*OvsPort
}

func (o *OpenvSwitch) parsePortsFromSwitchInfo(switchInfo string) {
	o.Ports = map[string]*OvsPort{}
	portIndices := regexp.MustCompile(` [\d\w].*\(.*\):`).FindAllStringIndex(switchInfo, -1)
	indices := []int{}
	for _, n := range portIndices {
		indices = append(indices, n[0])
	}
	logrus.Debugf("Parse ports, start indices at %v", indices)
	for i := range indices {
		portInfo := ""
		if i < len(indices) - 1 {
			portInfo = switchInfo[indices[i]:indices[i+1]]
		} else {
			portInfo = switchInfo[indices[i]:]
		}
		port, _ := GetPort(portInfo)
		o.Ports[port.Name] = port
	}
}

func castStringToIntField(o *OpenvSwitch, k string, v string) error {
	intvalue, err := strconv.Atoi(v)
	if err != nil {
		intvalue = 0
	}
	if err := utils.SetField(o, k, intvalue); err != nil {
		return err
	}
	return nil
}

func castStringToStringSliceField(o *OpenvSwitch, k string, v string) error {
	stringslice := strings.Split(v, " ")
	if err := utils.SetField(o, k, stringslice); err != nil {
		return err
	}
	return nil
}

func (o *OpenvSwitch) Fill(m map[string]interface{}) error {
	for k, v := range m {
		logrus.Debugf("Evaluating key %s with value %v", k, v)
		switch strings.ToLower(k) {
		case "n_tables":
			bytes := regexp.MustCompile("[0-9]+").FindAllString(v.(string), -1)
			if err := castStringToIntField(o, k, bytes[0]); err != nil {
				return err
			}
		case "n_buffers":
			if err := castStringToIntField(o, k, v.(string)); err != nil {
				return err
			}
		case "capabilities":
			if err := castStringToStringSliceField(o, k, v.(string)); err != nil {
				return err
			}
		case "actions":
			if err := castStringToStringSliceField(o, k, v.(string)); err != nil {
				return err
			}
		default:
			if err := utils.SetField(o, k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func GetSwitch(switchInfo string) (*OpenvSwitch, error) {
	logrus.Infof("ovs-ofctl show switch")
	sw := &OpenvSwitch{}
	switchInterface := utils.MapStringToInterface(*sw, switchInfo)
	logrus.Infof("Interface from output : %v", switchInterface)
	if err := sw.Fill(switchInterface); err != nil {
		return sw, err
	}
	sw.parsePortsFromSwitchInfo(switchInfo)
	return sw, nil
}