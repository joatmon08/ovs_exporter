package openvswitch

import (
	"github.com/socketplane/libovsdb"
	"github.com/Sirupsen/logrus"
)

func CheckHealth(o *libovsdb.OvsdbClient) ([]string, error) {
	return o.ListDbs()
}

func GetTotalFromTable(o *libovsdb.OvsdbClient, table string) []map[string]interface{} {
	rows := []map[string]interface{}{}
	op := libovsdb.Operation{
		Op:        "select",
		Table:     table,
		Where:     []interface{}{},
	}
	reply, _ := o.Transact("Open_vSwitch", op)
	logrus.Debugf("Reply from OVSDB, %v", reply)
	for _, r := range reply {
		return r.Rows
	}
	return rows
}