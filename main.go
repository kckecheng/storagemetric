package main

import (
	// "fmt"
	"github.com/kckecheng/storagemetric/dell/emc/powermax"
	"github.com/kckecheng/storagemetric/utils"
)

func main() {
	// var err error
	pmax, _ := powermax.New("10.228.234.200", "8443", "smc", "smc", "000197900151")

	dirs := pmax.GetFEDirectors()
	utils.PrettyPrint(dirs, "", "")

	ports := pmax.GetDirPorts("FA-1D")
	utils.PrettyPrint(ports, "", "")
}
