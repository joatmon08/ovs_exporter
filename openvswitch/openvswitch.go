package openvswitch

import (
	"github.com/socketplane/libovsdb"
	"github.com/Sirupsen/logrus"
	"errors"
)

const STATISTICS = "statistics"
const NAME = "name"
const UUID = "_uuid"

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
	logrus.Debugf("%v, Reply from OVSDB, %v", op, reply)
	for _, r := range reply {
		return r.Rows
	}
	return rows
}