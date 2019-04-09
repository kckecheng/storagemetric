package powermax

import (
	"time"
	// "fmt"
	// "github.com/kckecheng/storagemetric/utils"
)

type StorageGroupMetric struct {
	HostReads         float64 `json:"HostReads"`
	HostWrites        float64 `json:"HostWrites"`
	HostMBReads       float64 `json:"HostMBReads"`
	HostMBWritten     float64 `json:"HostMBWritten"`
	ReadResponseTime  float64 `json:"ReadResponseTime"`
	WriteResponseTime float64 `json:"WriteResponseTime"`
	ResponseTime      float64 `json:"ResponseTime"`
	AvgIOSize         float64 `json:"AvgIOSize"`
	AvgReadSize       float64 `json:"AvgReadSize"`
	AvgWriteSize      float64 `json:"AvgWriteSize"`
	Timestamp         int64   `json:"timestamp"`
}

type ArrayMetric struct {
	HostIOs           float64 `json:"HostIOs"`
	HostReads         float64 `json:"HostReads"`
	HostWrites        float64 `json:"HostWrites"`
	HostMBReads       float64 `json:"HostMBReads"`
	HostMBWritten     float64 `json:"HostMBWritten"`
	FEReadReqs        float64 `json:"FEReadReqs"`
	FEWriteReqs       float64 `json:"FEWriteReqs"`
	ReadResponseTime  float64 `json:"ReadResponseTime"`
	WriteResponseTime float64 `json:"WriteResponseTime"`
	FEUtilization     float64 `json:"FEUtilization"`
	Timestamp         int64   `json:"timestamp"`
}

// GetFEDirectors List available FE directors
func (pmax *PowerMax) GetFEDirectors() []string {
	var dirs []string

	payload := struct {
		SymmetrixId string `json:"symmetrixId"`
	}{pmax.symmid}

	fedirs := struct {
		FEDirectorInfo []struct {
			DirectorId         string `json:"directorId"`
			FirstAvailableDate int64  `json:"firstAvailableDate"`
			LastAvailableDate  int64  `json:"lastAvailableDate"`
		} `json:"feDirectorInfo"`
	}{}
	pmax.Request("POST", "/univmax/restapi/performance/FEDirector/keys", payload, &fedirs)
	for _, dir := range fedirs.FEDirectorInfo {
		dirs = append(dirs, dir.DirectorId)
	}

	return dirs
}

func (pmax *PowerMax) GetDirPorts(dir string) []string {
	var ports []string

	payload := struct {
		SymmetrixId string `json:"symmetrixId"`
		DirectorId  string `json:"directorId"`
	}{pmax.symmid, dir}

	dirports := struct {
		FePortInfo []struct {
			PortId             string `json:"portId"`
			FirstAvailableDate string `json:"firstAvailableDate"`
			LastAvailableDate  string `json:"lastAvailableDate"`
		} `json:"fePortInfo"`
	}{}
	pmax.Request("POST", "/univmax/restapi/performance/FEPort/keys", payload, &dirports)
	for _, port := range dirports.FePortInfo {
		ports = append(ports, port.PortId)
	}

	return ports
}

func (pmax *PowerMax) GetStorageGroups() []string {
	var sgs []string

	payload := struct {
		SymmetrixId string `json:"symmetrixId"`
	}{pmax.symmid}

	sgroups := struct {
		StorageGroupInfo []struct {
			StorageGroupId     string `json:"storageGroupId"`
			FirstAvailableDate int64  `json:"firstAvailableDate"`
			LastAvailableDate  int64  `json:"lastAvailableDate"`
		} `json:"storageGroupInfo"`
	}{}
	pmax.Request("POST", "/univmax/restapi/performance/StorageGroup/keys", payload, &sgroups)
	for _, sg := range sgroups.StorageGroupInfo {
		sgs = append(sgs, sg.StorageGroupId)
	}

	return sgs
}

func (pmax *PowerMax) GetStorageGroupMetric(sg string, from time.Time, to time.Time) StorageGroupMetric {
	payload := struct {
		SymmetrixId    string   `json:"symmetrixId"`
		StorageGroupId string   `json:"storageGroupId"`
		DataFormat     string   `json:"dataFormat"`
		StartDate      int64    `json:"startDate"`
		EndDate        int64    `json:"endDate"`
		Metrics        []string `json:"metrics"`
	}{
		pmax.symmid,
		sg,
		"Average",
		dateToTimestamp(from),
		dateToTimestamp(to),
		[]string{"HostReads", "HostWrites", "HostMBReads", "HostMBWritten", "ResponseTime", "ReadResponseTime", "WriteResponseTime", "AvgIOSize", "AvgReadSize", "AvgWriteSize"},
	}
	result := struct {
		ResultList struct {
			Result []StorageGroupMetric `json:"result"`
			From   int64                `json:"from"`
			To     int64                `json:"to"`
		} `json:"resultList"`
		Id             string `json:"id"`
		Count          int64  `json:"count"`
		ExpirationTime int64  `json:"expirationTime"`
		MaxPageSize    int64  `json:"maxPageSize"`
		WarningMessage string `json:"warningMessage"`
	}{}
	pmax.Request("POST", "/univmax/restapi/performance/StorageGroup/metrics", payload, &result)

	metrics := result.ResultList.Result
	// Only return the latest result if exists
	if len(metrics) == 0 {
		return StorageGroupMetric{}
	}
	return metrics[len(metrics)-1]
}

func (pmax *PowerMax) GetArrayMetric(from time.Time, to time.Time) ArrayMetric {
	payload := struct {
		SymmetrixId string   `json:"symmetrixId"`
		DataFormat  string   `json:"dataFormat"`
		StartDate   int64    `json:"startDate"`
		EndDate     int64    `json:"endDate"`
		Metrics     []string `json:"metrics"`
	}{
		pmax.symmid,
		"Average",
		dateToTimestamp(from),
		dateToTimestamp(to),
		[]string{"HostIOs", "HostReads", "HostWrites", "HostMBReads", "HostMBWritten", "FEReadReqs", "FEWriteReqs", "ReadResponseTime", "WriteResponseTime", "FEUtilization"},
	}
	result := struct {
		ResultList struct {
			Result []ArrayMetric `json:"result"`
			From   int64         `json:"from"`
			To     int64         `json:"to"`
		} `json:"resultList"`
		Id             string `json:"id"`
		Count          int64  `json:"count"`
		ExpirationTime int64  `json:"expirationTime"`
		MaxPageSize    int64  `json:"maxPageSize"`
		WarningMessage string `json:"warningMessage"`
	}{}
	pmax.Request("POST", "/univmax/restapi/performance/Array/metrics", payload, &result)

	metrics := result.ResultList.Result
	// Only return the latest result if exists
	if len(metrics) == 0 {
		return ArrayMetric{}
	}
	return metrics[len(metrics)-1]
}
