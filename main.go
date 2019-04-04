package main

import (
	// "fmt"
	"github.com/kckecheng/storagemetric/dell/emc/powermax"
	"github.com/kckecheng/storagemetric/utils"
)

func main() {
	// var err error
	pmax, _ := powermax.New("10.228.234.200", "8443", "smc", "smc", "000197900151")

	payload := struct {
		StartDate   int64    `json:"startDate"`
		EndDate     int64    `json:"endDate"`
		SymmetrixId string   `json:"symmetrixId"`
		DataFormat  string   `json:"dataFormat"`
		Metrics     []string `json:"metrics"`
	}{
		1257894000000,
		1554359805000,
		"000197900151",
		"Average",
		[]string{"OverallHealthScore"},
	}

	type arrayMetricsResultList struct {
		MaxPageSize    int64  `json:"maxPageSize"`
		ExpirationTime int64  `json:"expirationTime"`
		Count          int64  `json:"count"`
		Id             string `json:"id"`
		ResultList     struct {
			From   int64 `json:"from"`
			To     int64 `json:"to"`
			Result []struct {
				OverallHealthScore float64 `json:"OverallHealthScore"`
				Timestamp          int64   `json:"timestamp"`
			} `json:"result"`
		} `json:resultList`
	}
	var ametrics arrayMetricsResultList
	pmax.Request("POST", "/univmax/restapi/performance/Array/metrics", payload, &ametrics)
	utils.PrettyPrint(ametrics, "", "")
}
