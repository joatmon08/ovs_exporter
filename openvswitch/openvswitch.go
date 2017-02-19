package openvswitch

import "github.com/socketplane/libovsdb"

func CheckHealth(o *libovsdb.OvsdbClient) ([]string, error) {
	return o.ListDbs()
}

func GetTotalBridges(o *libovsdb.OvsdbClient) int {
	op := libovsdb.Operation{
		Op:        "select",
		Table:     "Open_vSwitch",
		Where:     []interface{}{},
	}
	reply, _ := o.Transact("Open_vSwitch", op)
	for _, r := range reply {
		return r.Count
	}
	return 0
}

