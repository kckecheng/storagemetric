package powermax

import (
// "fmt"
// "github.com/kckecheng/storagemetric/utils"
)

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
