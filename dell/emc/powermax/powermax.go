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

// Add common headers
func populateCommonHeaders(req *http.Request) {
	type httpHeader struct {
		header string
		value  string
	}

	headers := []httpHeader{
		{"Accept", "application/json"},
		{"Content-Type", "application/json"},
		{"X-EMC-REST-CLIENT", "true"},
	}
	for _, header := range headers {
		req.Header.Set(header.header, header.value)
	}
}

// New Init PowerMax Object
func New(server string, port string, username string, password string, symmid string) (*PowerMax, error) {
	if utils.Logger == nil {
		utils.InitLogger("", "")
	}

	var err error
	utils.Log("debug", fmt.Sprintf("server: %s, username: %s, password: %s, port: %s, symmid: %s", server, username, password, port, symmid))
	if server == "" || username == "" || password == "" || symmid == "" {
		return nil, errors.New("PowerMax server address, username, password and symmid must be specified")
	}

	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr, Jar: cookieJar}

	// Check if the provided parameters are correct
	req, err := http.NewRequest("GET", utils.URL("https", server, port, "/univmax/restapi/system/symmetrix", "/"+symmid), nil)
	req.SetBasicAuth(username, password)
	populateCommonHeaders(req)

	reqDetails, _ := httputil.DumpRequest(req, true)
	utils.Log("debug", fmt.Sprintf("Query Symmtric Request: %s", string(reqDetails)))
	resp, err := client.Do(req)
	if err != nil {
		utils.Log("error", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	respDetails, _ := httputil.DumpResponse(resp, true)
	fmt.Printf("%s", string(respDetails))

	if resp.StatusCode != 200 {
		message := fmt.Sprintf("Fail to query symmtric with symmid %s", symmid)
		utils.Log("error", message)
		respDetails, _ := httputil.DumpResponse(resp, true)
		utils.Log("debug", fmt.Sprintf("Query Symmtric Response: %s", string(respDetails)))
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
	if method == "" || URI == "" {
		return errors.New("method, or URI is missed")
	}

	url := utils.URL("https", pmax.server, pmax.port, URI)

	var req *http.Request
	if payload != nil {
		payloadJSON, _ := json.Marshal(payload)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(payloadJSON))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	populateCommonHeaders(req)

	reqDetails, _ := httputil.DumpRequest(req, true)
	utils.Log("debug", fmt.Sprintf("Request: %s", string(reqDetails)))

	resp, err := pmax.client.Do(req)
	if err != nil {
		utils.Log("error", err.Error())
		return err
	}
	defer resp.Body.Close()

	respDetails, _ := httputil.DumpResponse(resp, true)
	utils.Log("debug", string(respDetails))
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		respRaw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if result != nil {
			json.Unmarshal(respRaw, result)
		}

		return nil
	}
	utils.Log("error", "Fail to perform the request")
	return fmt.Errorf("Request Fails: %s", requestParams)
}
