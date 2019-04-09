package main

import (
	// "fmt"
	"github.com/kckecheng/storagemetric/dell/emc/powermax"
	"github.com/kckecheng/storagemetric/utils"
	"time"
)

func main() {
	// var err error
	pmax, _ := powermax.New("10.228.234.200", "8443", "smc", "smc", "000197900151")

	dirs := pmax.GetFEDirectors()
	utils.PrettyPrint(dirs, "", "")

	ports := pmax.GetDirPorts("FA-1D")
	utils.PrettyPrint(ports, "", "")

	sgs := pmax.GetStorageGroups()
	utils.PrettyPrint(sgs, "", "")

	current_tm := time.Now()
	from_tm := current_tm.Add(-time.Second * 360)

	var sgmetric powermax.StorageGroupMetric
	sgmetric = pmax.GetStorageGroupMetric("vmw-automation-ci", from_tm, current_tm)
	utils.PrettyPrint(sgmetric, "", "")

	var arrmetric powermax.ArrayMetric
	arrmetric = pmax.GetArrayMetric(from_tm, current_tm)
	utils.PrettyPrint(arrmetric, "", "")
}
