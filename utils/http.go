package utils

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
)

func InitHttpClient() http.Client {
	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr, Jar: cookieJar}
	return client
}

func InitHttpRequest(method string, url string, payload interface{}) (*http.Request, error) {
	requestInfo := fmt.Sprintf("method: %s, url: %s", method, url)
	Log("debug", requestInfo)

	if EmptyStrExists(method, url) == true {
		return nil, errors.New("method, or url is missed")
	}

	var req *http.Request
	if payload != nil {
		payloadJSON, _ := json.Marshal(payload)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(payloadJSON))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	return req, nil
}

func UpdateHttpRequestHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		Log("debug", fmt.Sprintf("Update header: %s -> %s", k, v))
		req.Header.Set(k, v)
	}
}

func UpdateHttpRequestParams(req *http.Request, pairs map[string]string) {
	params := req.URL.Query()
	for k, v := range pairs {
		Log("debug", fmt.Sprintf("Update query parameters: %s -> %s", k, v))
		params.Add(k, v)
	}
	req.URL.RawQuery = params.Encode()
}

func DoHttpRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	reqDetails, _ := httputil.DumpRequest(req, true)
	Log("debug", fmt.Sprintf("Request Details:\n%s", string(reqDetails)))

	resp, err := client.Do(req)
	respDetails, _ := httputil.DumpResponse(resp, true)
	Log("debug", fmt.Sprintf("Response Details:\n%s", string(respDetails)))
	return resp, err
}

func GetHttpResponseJson(resp *http.Response, result interface{}) error {
	var err error

	defer resp.Body.Close()
	respRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log("error", fmt.Sprintf("Fail to read response body due to %s", err.Error()))
		return err
	}

	err = json.Unmarshal(respRaw, result)
	if err != nil {
		Log("error", fmt.Sprintf("Fail to decode json from response body due to %s", err.Error()))
		return err
	}
	return nil
}
