package powermax

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"time"

	"github.com/kckecheng/storagemetric/utils"
)

// PowerMax PowerMax array object
type PowerMax struct {
	server   string
	port     string
	username string
	password string
	symmid   string
	client   http.Client
}

// Covert UTC timestamp(millisecond) to date
func timestampToDate(ms int64) time.Time {
	tm := time.Unix(ms/1000, 0)
	return tm
}

// Convert date to UTC timestamp(millisecond)
func dateToTimestamp(tm time.Time) int64 {
	return tm.Unix() * 1000
}

// Get current UTC timestamp
func currentTimestamp() int64 {
	return time.Now().Unix() * 1000
}

// Add common headers
func populateCommonHeaders(req *http.Request) {
	headers := map[string]string{
		"Accept":            "application/json",
		"Content-Type":      "application/json",
		"X-EMC-REST-CLIENT": "true",
	}
	utils.UpdateHttpRequestHeaders(req, headers)
}

// New Init PowerMax Object
func New(server string, port string, username string, password string, symmid string) (*PowerMax, error) {
	if utils.Logger == nil {
		utils.InitLogger("", "")
	}

	var err error
	utils.Log("debug", fmt.Sprintf("server: %s, username: %s, password: %s, port: %s, symmid: %s", server, username, password, port, symmid))
	if utils.EmptyStrExists(server, username, password, symmid) == true {
		return nil, errors.New("PowerMax server address, username, password and symmid must be specified")
	}

	client := utils.InitHttpClient()

	// Check if the provided parameters are correct
	req, err := utils.InitHttpRequest("GET", utils.URL("https", server, port, "/univmax/restapi/system/symmetrix", "/"+symmid), nil)
	req.SetBasicAuth(username, password)
	populateCommonHeaders(req)

	utils.DoHttpRequest(&client, req)
	if err != nil {
		utils.Log("error", err.Error())
		return nil, err
	}

	if resp.StatusCode != 200 {
		message := fmt.Sprintf("Fail to query symmtric with symmid %s", symmid)
		utils.Log("error", message)
		return nil, errors.New(message)
	}

	return &PowerMax{
		server:   server,
		username: username,
		password: password,
		port:     port,
		symmid:   symmid,
		client:   client,
	}, nil
}

// Request Send get/post/delete request
func (pmax *PowerMax) Request(method string, URI string, payload interface{}, result interface{}) error {
	requestParams := fmt.Sprintf("method: %s, URI: %s, payload: %#v", method, URI, payload)
	utils.Log("debug", requestParams)
	if utils.EmptyStrExists(method, URI) == true {
		return errors.New("method, or URI is missed")
	}

	url := utils.URL("https", pmax.server, pmax.port, URI)
	req, err := utils.InitHttpRequest(method, url, payload)

	req.SetBasicAuth(pmax.username, pmax.password)
	populateCommonHeaders(req)

	resp, err := utils.DoHttpRequest(&pmax.client, req)
	if err != nil {
		utils.Log("error", err.Error())
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		err := utils.GetHttpResponseJson(resp, result)
		return err
	}
}
