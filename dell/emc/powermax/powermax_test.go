package powermax

import (
	"flag"
	"fmt"
	"testing"
	"time"
)

var server = flag.String("server", "", "PowerMax Unisphere IP/FQDN")
var port = flag.String("port", "8443", "PowerMax Unisphere port, 8443 as default")
var username = flag.String("username", "", "PowerMax user name")
var password = flag.String("password", "", "PowerMax user password")
var symmid = flag.String("symmid", "", "PowerMax symmetrix id")
var interval = flag.Int("interval", 10, "Interval in seconds to collect metric, 10 as default")

func FailIfError(t *testing.T, err error) {
	if err != nil {
		t.Log(fmt.Sprintf("%s", err.Error()))
		t.FailNow()
	}
}

func TestPowerMax(t *testing.T) {
	if *server == "" || *username == "" || *password == "" || *symmid == "" {
		t.Log(fmt.Sprintf("server, username, and password must be specified as -args -server <IP/FQDN> -username <user> -password <password> -symmid <symmid>"))
		t.FailNow()
	}

	var err error
	pmax, err := New(*server, *port, *username, *password, *symmid)
	FailIfError(t, err)

	current_tm := time.Now()
	from_tm := current_tm.Add(-time.Second * time.Duration(*interval))
	arrmetric := pmax.GetArrayMetric(from_tm, current_tm)
	t.Log(arrmetric)
}
