package utils

import (
	"testing"
	"net/http/httptest"
	"fmt"
	"net/http"
)

func TestNewOVSExporterClient(t *testing.T) {
	endpoint := "http://localhost:8080"
	client := NewOVSExporterClient(endpoint)
	if client.Endpoint != endpoint + "/metrics" {
		t.Errorf("expected %s, actual %s", endpoint + "/metrics", client.Endpoint)
	}
}

func TestGetExporterMetrics(t *testing.T) {
	sampleresponse := "openvswitch_interfaces_total 2\n"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, sampleresponse)
	}))
	defer ts.Close()
	client := NewOVSExporterClient(ts.URL)
	metrics, err := client.GetExporterMetrics()
	if err != nil {
		t.Errorf("error retrieving metrics %s", err)
	}
	if metrics["openvswitch_interfaces_total"] != "2" {
		t.Errorf("expected %s, actual %s", "2", metrics["openvswitch_interfaces_total"])
	}
}