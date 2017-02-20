package openvswitch

import (
	"github.com/socketplane/libovsdb"
	"github.com/Sirupsen/logrus"
	"errors"
)

const STATISTICS = "statistics"
const NAME = "name"

func CheckHealth(o *libovsdb.OvsdbClient) ([]string, error) {
	return o.ListDbs()
}

func ParseStatisticsFromData(rows []map[string]interface{}) (map[string]map[string]float64, error) {
	result := map[string]map[string]float64{}
	for _, row := range rows {
		name := row[NAME].(string)
		fieldDetails := row[STATISTICS].([]interface{})
		if fieldDetails[0].(string) != "map" {
			return nil, errors.New("field " + STATISTICS + " in OVSDB is not of map type")
		}
		info := fieldDetails[1].([]interface{})
		entryMap := map[string]float64{}
		for _, entry := range info {
			e := entry.([]interface{})
			entryMap[e[0].(string)] = e[1].(float64)
		}
		result[name] = entryMap
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