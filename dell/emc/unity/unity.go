package unity

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

// Unity Unity array object
type Unity struct {
	server   string
	username string
	password string
	token    string
	client   http.Client
}

// URL join server base address with a URI and return the final full URL
func URL(server, URI string) string {
	return "https://" + server + URI
}

// New Init Unity Object
func New(server string, username string, password string) (*Unity, error) {
	if utils.Logger == nil {
		utils.InitLogger("", "")
	}

	var err error
	utils.Log("debug", fmt.Sprintf("server: %s, username: %s, password: %s", server, username, password))
	if server == "" || username == "" || password == "" {
		return nil, errors.New("Unity server address, username, and password must be specified")
	}

	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr, Jar: cookieJar}

	req, err := http.NewRequest("GET", "https://"+server+"/api/types/loginSessionInfo/instances", nil)
	req.SetBasicAuth(username, password)
	req.Header.Set("X-EMC-REST-CLIENT", "true")

	reqDetails, _ := httputil.DumpRequest(req, true)
	utils.Log("debug", fmt.Sprintf("Login Request: %s", string(reqDetails)))
	resp, err := client.Do(req)
	if err != nil {
		utils.Log("error", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		message := fmt.Sprintf("Login %s with %s/%s failed.", server, username, password)
		utils.Log("error", message)
		respDetails, _ := httputil.DumpResponse(resp, true)
		utils.Log("debug", fmt.Sprintf("Login Response: %s", string(respDetails)))
		return nil, errors.New(message)
	}

	token := resp.Header.Get("Emc-Csrf-Token")

	return &Unity{
		server:   server,
		username: username,
		password: password,
		token:    token,
		client:   client,
	}, nil
}

// Request Send get/post/delete request
func (unity *Unity) Request(method string, URI string, fields string, filter string, payload interface{}, result interface{}) error {
	requestParams := fmt.Sprintf("method: %s, URI: %s, fields: %s, filter: %s, payload: %#v", method, URI, fields, filter, payload)
	utils.Log("debug", requestParams)
	if method == "" || URI == "" {
		return errors.New("method, or URI is missed")
	}

	url := "https://" + unity.server + URI

	var req *http.Request
	if payload != nil {
		payloadJSON, _ := json.Marshal(payload)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(payloadJSON))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	if method == "POST" || method == "DELETE" {
		req.Header.Set("EMC-CSRF-TOKEN", unity.token)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-EMC-REST-CLIENT", "true")

	params := req.URL.Query()
	if method != "DELETE" {
		params.Add("compact", "true")
	}
	if fields != "" {
		params.Add("fields", fields)
	}
	if filter != "" {
		params.Add("filter", filter)
	}
	req.URL.RawQuery = params.Encode()

	reqDetails, _ := httputil.DumpRequest(req, true)
	utils.Log("debug", fmt.Sprintf("Request: %s", string(reqDetails)))

	resp, err := unity.client.Do(req)
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

// Destroy logout Unity
func (unity *Unity) Destroy() error {
	utils.Log("debug", "Logout")
	err := unity.Request("POST", "/api/types/loginSessionInfo/action/logout", "", "", nil, nil)
	return err
}
