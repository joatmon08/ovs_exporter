package utils

import (
	"net/http"
	"io/ioutil"
	"strings"
)

type OVSExporterClient struct {
	Endpoint string
	Client   *http.Client
}

func NewOVSExporterClient(endpoint string) *OVSExporterClient {
	metricsEndpoint := endpoint + "/metrics"
	return &OVSExporterClient{
		Endpoint: metricsEndpoint,
		Client: &http.Client{},
	}
}

func (c *OVSExporterClient) GetExporterMetrics() (map[string]string, error) {
	metrics := map[string]string{}
	request, err := http.NewRequest("GET", c.Endpoint, nil)
	if err != nil {
		return metrics, err
	}
	//request.Header.Set("Content-Type", "text/plain")
	response, err := c.Client.Do(request)
	if err != nil {
		return metrics, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return metrics, err
	}
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if !strings.Contains(line, "# ") && len(line) > 0 {
			metric := strings.Fields(line)
			metrics[metric[0]] = metric[1]
			if err != nil {
				return metrics, err
			}
		}
	}
	return metrics, err
}