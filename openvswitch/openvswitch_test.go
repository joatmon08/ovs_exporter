package openvswitch

import (
	"github.com/socketplane/libovsdb"
	"testing"
	"github.com/joatmon08/ovs_exporter/utils"
	"encoding/json"
	"reflect"
)

func setup(t *testing.T) *libovsdb.OvsdbClient {
	client, err := libovsdb.Connect("127.0.0.1", 6640)
	if err != nil {
		t.Error(err)
	}
	return client
}

func TestCheckHealth(t *testing.T) {
	client := setup(t)
	dbs, err := CheckHealth(client)
	if err != nil {
		t.Error(err)
	}
	if len(dbs) != 1 {
		t.Errorf("Expected %d, got %d", 1, len(dbs))
	}
}

func TestGetTotalFromTable(t *testing.T) {
	client := setup(t)
	rows := GetRowsFromTable(client, "Bridge")
	if len(rows) != 0 {
		t.Errorf("Expected %d, got %d", 0, len(rows))
	}
}

func TestParseStatisticsFromData(t *testing.T) {
	var test []map[string]interface{}
	expected := map[string]float64{
		"collisions": 0,
		"rx_bytes": 1026,
		"rx_crc_err": 0,
		"rx_dropped": 0,
		"rx_errors": 0,
		"rx_frame_err": 0,
		"rx_over_err": 0,
		"rx_packets": 0,
		"tx_bytes": 1096,
		"tx_dropped": 0,
		"tx_errors": 0,
		"tx_packets": 14,
	}
	stats, err := utils.ReadTestDataToBytes("statistics.json")
	if err != nil {
		t.Error(err)
	}
	if err := json.Unmarshal(stats, &test); err != nil {
		t.Error(err)
	}
	result, err := ParseStatisticsFromInterfaces(test)
	if err != nil {
		t.Error(err)
	}
	if len(result) != 6 {
		t.Errorf("Expected %d, got %d", 6, len(result))
	}
	for _, r := range result {
		if r.Name == "1c1988eb903b4_l" {
			if r.UUID != "aa713415-8566-458b-b8ef-58e550af8a91" {
				t.Errorf("Expected %s, got %s", "aa713415-8566-458b-b8ef-58e550af8a91", r.UUID)
			}
			if reflect.DeepEqual(r.Statistics, expected) {
				t.Errorf("Expected %v, got %v", expected, r.Statistics)
			}
		}
	}
}

func TestParsePortsFromBridges(t *testing.T) {
	var test []map[string]interface{}
	expected := Bridge{
		Name: "ovsbr0",
		Ports: []Port{
			{UUID: "a5956ae0-25fd-46b2-a881-a6f63c8014d9"},
			{UUID: "dfb05de8-617b-4127-91aa-e1f3bfd7ab60"},
		},
	}
	stats, err := utils.ReadTestDataToBytes("bridges.json")
	if err != nil {
		t.Error(err)
	}
	if err := json.Unmarshal(stats, &test); err != nil {
		t.Error(err)
	}
	result, err := ParsePortsFromBridges(test)
	if err != nil {
		t.Error(err)
	}
	if len(result) != 3 {
		t.Errorf("Expected %d, got %d", 3, len(result))
	}
	for _, r := range result {
		if r.Name == "ovsbr0" {
			if !reflect.DeepEqual(expected, r) {
				t.Errorf("Expected %v, got %v", expected, r)
			}
		}
	}
}
