package ovsofctl

import (
	"github.com/joatmon08/ovs_exporter/utils"
	"testing"
	"reflect"
)

func TestGetSwitch(t *testing.T) {
	testSwitchInfo, err := utils.ReadTestDataToString("switch")
	if err != nil {
		t.Error(err)
	}
	sw, err := GetSwitchDetails(testSwitchInfo)
	if err != nil {
		t.Error(err)
	}
	expectedSwitch := OpenvSwitch{
		N_tables: 254,
		N_buffers: 256,
		Capabilities: []string{
			"FLOW_STATS","TABLE_STATS",
			"PORT_STATS", "QUEUE_STATS", "ARP_MATCH_IP"},
		Actions: []string{
			"output", "enqueue",
			"set_vlan_vid", "set_vlan_pcp",
			"strip_vlan", "mod_dl_src",
			"mod_dl_dst", "mod_nw_src",
			"mod_nw_dst", "mod_nw_tos",
			"mod_tp_src", "mod_tp_dst"},
		Ports: map[string]*OvsPort{
			"87be70c563b24_l": {
				ID: "1",
				Name: "87be70c563b24_l",
				Addr: "d2:e5:b4:c9:e9:86",
				Current: "10GB-FD COPPER",
				Speed: 10000,
				MaxSpeed: 0,
			},
			"1c1988eb903b4_l": {
				ID: "2",
				Name: "1c1988eb903b4_l",
				Addr: "a2:c1:a6:7a:f5:80",
				Current: "10GB-FD COPPER",
				Speed: 10000,
				MaxSpeed: 0,
			},
			"ovs-br1": {
				ID: "LOCAL",
				Name: "ovs-br1",
				Addr: "02:ab:c7:77:71:40",
				Speed: 0,
				MaxSpeed: 0,
			},
		},
	}
	if expectedSwitch.N_buffers != sw.N_buffers {
		t.Errorf("Expected %v, got %v", expectedSwitch.N_buffers, sw.N_buffers)
	}
	if expectedSwitch.N_tables != sw.N_tables {
		t.Errorf("Expected %v, got %v", expectedSwitch.N_tables, sw.N_tables)
	}
	if len(expectedSwitch.Capabilities) != len(sw.Capabilities) {
		t.Errorf("Expected %d, got %d", len(expectedSwitch.Capabilities), len(sw.Capabilities))
	}
	if len(expectedSwitch.Actions) != len(sw.Actions) {
		t.Errorf("Expected %d, got %d", len(expectedSwitch.Actions), len(sw.Actions))
	}
	if len(expectedSwitch.Ports) != len(sw.Ports) {
		t.Errorf("Expected %d, got %d", len(expectedSwitch.Ports), len(sw.Ports))
	}
	if !reflect.DeepEqual(expectedSwitch.Ports, sw.Ports) {
		t.Errorf("Expected %d, got %d", expectedSwitch.Ports, sw.Ports)
	}
}