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

// New Init Unity Object
func New(server string, username string, password string) (*Unity, error) {
	if utils.Logger == nil {
		utils.InitLogger("", "")
	}

	var err error
	utils.Log("debug", fmt.Sprintf("server: %s, username: %s, password: %s", server, username, password))
	if utils.EmptyStrExists(server, username, password) == true {
		return nil, errors.New("Unity server address, username, and password must all be specified")
	}

	client := utils.InitHttpClient()

	req, err := utils.InitHttpRequest("GET", utils.URL("https", server, "", "/api/types/loginSessionInfo/instances"), nil)
	req.SetBasicAuth(username, password)
	req.Header.Set("X-EMC-REST-CLIENT", "true")

	resp, err := utils.DoHttpRequest(&client, req)
	if err != nil {
		utils.Log("error", err.Error())
		return nil, err
	}

	if resp.StatusCode != 200 {
		message := fmt.Sprintf("Login %s with %s/%s failed.", server, username, password)
		utils.Log("error", message)
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
	if utils.EmptyStrExists(method, URI) == true {
		return errors.New("method, or URI is missed")
	}

	url := utils.URL("https", unity.server, "", URI)
	req, err := utils.InitHttpRequest(method, url, payload)

	// Add token in header if POST/DELETE
	if method == "POST" || method == "DELETE" {
		req.Header.Set("EMC-CSRF-TOKEN", unity.token)
	}

	commonHeaders := map[string]string{
		"Accept":            "application/json",
		"Content-Type":      "application/json",
		"X-EMC-REST-CLIENT": "true",
	}
	utils.UpdateHttpRequestHeaders(req, commonHeaders)

	if method != "DELETE" {
		utils.UpdateHttpRequestParams(req, map[string]string{"compact": "true"})
	}
	if fields != "" {
		utils.UpdateHttpRequestParams(req, map[string]string{"fields": fields})
	}
	if filter != "" {
		utils.UpdateHttpRequestParams(req, map[string]string{"filter": filter})
	}

	resp, err := utils.DoHttpRequest(&unity.client, req)
	if err != nil {
		utils.Log("error", err.Error())
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		err := utils.GetHttpResponseJson(resp, result)
		return err
	}
	return errors.New("Invalid status code")
}

// Destroy logout Unity
func (unity *Unity) Destroy() error {
	utils.Log("debug", "Logout")
	err := unity.Request("POST", "/api/types/loginSessionInfo/action/logout", "", "", nil, nil)
	return err
}
