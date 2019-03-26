package unity

import (
	"fmt"

	"github.com/kckecheng/storagemetric/utils"
)

// NewMetricRealTimeQuery Create a real time query
func (unity *Unity) NewMetricRealTimeQuery(paths []string, interval int) (int, error) {
	payload := struct {
		Paths    []string `json:"paths"`
		Interval int      `json:"interval"`
	}{paths, interval}
	var ret MetricRealTimeQuery
	err := unity.Request("POST", "/api/types/metricRealTimeQuery/instances", "", "", payload, &ret)

	if err != nil {
		utils.Log("error", "Fail to create a new metric real time query")
		return 0, err
	}
	id := ret.Content.Id
	utils.Log("debug", fmt.Sprintf("New Metric Real Time Query ID: %d", id))
	return id, nil
}

// MetricRealTimeQueryExisted Check if a real time query with the specified ID exists
func (unity *Unity) MetricRealTimeQueryExisted(id int) bool {
	utils.Log("debug", fmt.Sprintf("Check if the specified query id %d exists", id))
	err := unity.Request("GET", fmt.Sprintf("/api/instances/metricRealTimeQuery/%d", id), "", "", "", nil)
	if err != nil {
		utils.Log("warning", fmt.Sprintf("The query id %d does not exist", id))
		return false
	}
	utils.Log("info", fmt.Sprintf("The query id %d exists", id))
	return true
}

// DeleteMetricRealTimeQuery Delete a real time query based on its ID
func (unity *Unity) DeleteMetricRealTimeQuery(id int) error {
	utils.Log("debug", fmt.Sprintf("Delete metric real time query %d", id))
	uri := fmt.Sprintf("/api/instances/metricRealTimeQuery/%d", id)
	err := unity.Request("DELETE", uri, "", "", nil, nil)
	return err
}

// GetMetricQueryResult Get metric results
// Pitfall: Metric will be empty if it is retrivded without waiting for at least a query interval after creating the query
func (unity *Unity) GetMetricQueryResult(id int, result *Metric) error {
	utils.Log("debug", fmt.Sprintf("Get metric result with id %d", id))
	filter := fmt.Sprintf("queryId eq %d", id)
	err := unity.Request("GET", "/api/types/metricQueryResult/instances", "", filter, nil, result)
	return err
}

// GetHistoricalMetric Query historial metrics
func (unity *Unity) GetHistoricalMetric(path string, result *Metric) error {
	utils.Log("debug", fmt.Sprintf("Get historical metric data with path %s", path))
	filter := fmt.Sprintf("path eq \"%s\"", path)
	err := unity.Request("GET", "/api/types/metricValue/instances", "", filter, nil, result)
	return err
}
