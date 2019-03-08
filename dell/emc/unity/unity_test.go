package unity

import (
	"flag"
	"fmt"
	"testing"
	"time"
)

var server = flag.String("server", "", "Unity IP/FQDN")
var username = flag.String("username", "admin", "Unity user name")
var password = flag.String("password", "", "Unity user password")
var interval = flag.Int("interval", 10, "Interval in seconds to collect metric")

func FailIfError(t *testing.T, err error) {
	if err != nil {
		t.Log(fmt.Sprintf("%s", err.Error()))
		t.FailNow()
	}
}

func TestUnity(t *testing.T) {
	if *server == "" || *username == "" || *password == "" {
		t.Log(fmt.Sprintf("server, username, and password must be specified as -args -server <IP/FQDN> -username <user> -password <password>"))
		t.FailNow()
	}

	var err error
	unityBox, err := New(*server, *username, *password)
	FailIfError(t, err)

	var pathSet [2][]string
	pathSet[0] = []string{
		"sp.*.cpu.summary.busyTicks",
		"sp.*.cpu.summary.idleTicks",
	}
	pathSet[1] = []string{
		"sp.*.memory.summary.freeBytes",
		"sp.*.memory.summary.totalBytes",
		"sp.*.memory.summary.totalUsedBytes",
	}

	for _, paths := range pathSet {
		var ret Metric
		var message string
		var err error

		message = fmt.Sprintf("---\nQuery Metric Paths: %+v\n---", paths)
		t.Log(message)

		id, err := unityBox.NewMetricRealTimeQuery(paths, *interval)
		FailIfError(t, err)
		existed := unityBox.MetricRealTimeQueryExisted(id)
		if existed == false {
			t.Log(fmt.Sprintf("Query id %d does not exists", id))
			t.FailNow()
		}
		// Sleep interval * seconds, otherwise, the first record may be empty
		time.Sleep(time.Duration(*interval) * time.Second)

		err = unityBox.GetMetricQueryResult(id, &ret)
		FailIfError(t, err)
		message = fmt.Sprintf("---\nMetric Paths: %+v\nResult:\n%+v\n---", paths, ret)
		t.Log(message)

		err = unityBox.DeleteMetricRealTimeQuery(id)
		FailIfError(t, err)
	}

	var histRet Metric
	historicalPaths := []string{
		"sp.*.cpu.summary.utilization",
	}
	for _, path := range historicalPaths {
		err := unityBox.GetHistoricalMetric(path, &histRet)
		FailIfError(t, err)
		message := fmt.Sprintf("---\nMetric Path: %+v\nResult:\n%+v\n---", path, histRet)
		t.Log(message)
	}

	err = unityBox.Destroy()
	FailIfError(t, err)
}
