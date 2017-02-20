package openvswitch

import (
	"github.com/socketplane/libovsdb"
	"github.com/Sirupsen/logrus"
	"errors"
	"strings"
)

const STATISTICS = "statistics"
const NAME = "name"
const UUID = "_uuid"
const PORTS = "ports"

type Bridge struct {
	Name  string
	Ports []Port
}

type Port struct {
	UUID string
	Name string
	Interfaces []Interface
}

type Interface struct {
	UUID string
	Name string
	Statistics map[string]float64
}

func CheckHealth(o *libovsdb.OvsdbClient) ([]string, error) {
	return o.ListDbs()
}

func ParsePortsFromBridges(rows []map[string]interface{}) ([]Bridge, error) {
	result := []Bridge{}
	for _, row := range rows {
		name := row[NAME].(string)
		fieldDetails := row[PORTS].([]interface{})
		ports := []Port{}
		if fieldDetails[0].(string) == "set" {
			info := fieldDetails[1].([]interface{})
			for _, entry := range info {
				e := entry.([]interface{})
				port := Port{UUID: e[1].(string)}
				ports = append(ports, port)
			}
		} else {
			uuid := fieldDetails[1].(string)
			port := Port{UUID: uuid}
			ports = append(ports, port)
		}
		result = append(result, Bridge{Name:name, Ports:ports})
	}
	logrus.WithFields(logrus.Fields{
		"event": "parsed ports from bridges",
		"rows": len(result),
	}).Info("retrieved bridge ports")
	return result, nil
}

func ParseStatisticsFromInterfaces(rows []map[string]interface{}) ([]Interface, error) {
	result := []Interface{}
	for _, row := range rows {
		uuid := row[UUID].([]interface{})[1].(string)
		name := row[NAME].(string)
		fieldDetails := row[STATISTICS].([]interface{})
		if fieldDetails[0].(string) != "map" {
			return nil, errors.New("field " + STATISTICS + " in OVSDB is not of map type")
		}
		info := fieldDetails[1].([]interface{})
		statMap := map[string]float64{}
		for _, entry := range info {
			e := entry.([]interface{})
			statMap[e[0].(string)] = e[1].(float64)
		}
		iface := Interface{
			UUID: uuid,
			Name: name,
			Statistics: statMap,
		}
		result = append(result, iface)
	}
	logrus.WithFields(logrus.Fields{
		"event": "parsed statistics from interfaces",
		"rows": len(result),
	}).Info("retrieved statistics")
	return result, nil
}

func GetRowsFromTable(o *libovsdb.OvsdbClient, table string) []map[string]interface{} {
	rows := []map[string]interface{}{}
	op := []libovsdb.Operation{{
		Op:        "select",
		Table:     table,
		Where:     []interface{}{},
	}}
	reply, _ := o.Transact("Open_vSwitch", op...)
	if strings.Contains(reply[0].Error, "error") {
		logrus.WithFields(logrus.Fields{
			"event": "ovsdb error",
			"table": table,
			"operation": op,
		}).Error(reply[0].Error)
	}
	for _, r := range reply {
		return r.Rows
	}
	logrus.WithFields(logrus.Fields{
		"event": "ovsdb transact",
		"table": table,
		"rows": len(rows),
	}).Info("retrieved ovsdb rows")
	return rows
}