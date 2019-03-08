package unity

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"os"
	"path/filepath"
)

type Unity struct {
	server   string
	username string
	password string
	token    string
	client   http.Client
	logger   *log.Logger
}

func getLogLevel(logStr string) log.Level {
	//Valid string: panic fatal error warning info debug trace
	level, err := log.ParseLevel(logStr)
	if err != nil {
		return log.DebugLevel
	} else {
		return level
	}
}

func initLogger(filename string, levelStr string) *log.Logger {
	var defaultLogFile string = "storagemetric.log"
	var level log.Level

	logger := log.New()

	if filename == "" {
		filename = filepath.Join(os.TempDir(), defaultLogFile)
	}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(f)
	}

	envLogStr, ok := os.LookupEnv("STORAGEMETRIC_LOGLEVEL")
	if ok {
		level = getLogLevel(envLogStr)
	} else {
		level = getLogLevel(levelStr)
	}
	logger.SetLevel(level)

	if level == log.DebugLevel {
		logger.SetReportCaller(true)
	}

	return logger
}

func URL(server, URI string) string {
	return "https://" + server + URI
}

func New(server string, username string, password string) (*Unity, error) {
	logger := initLogger("", "")

	var err error
	logger.Debug(fmt.Sprintf("server: %s, username: %s, password: %s", server, username, password))
	if server == "" || username == "" || password == "" {
		return nil, errors.New("Unity server address, username, and password must be specified.")
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
	logger.Debug(fmt.Sprintf("Login Request: %s", string(reqDetails)))
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		message := fmt.Sprintf("Login %s with %s/%s failed.", server, username, password)
		logger.Error(message)
		respDetails, _ := httputil.DumpResponse(resp, true)
		logger.Debug(fmt.Sprintf("Login Response: %s", string(respDetails)))
		return nil, errors.New(message)
	}

	token := resp.Header.Get("Emc-Csrf-Token")

	return &Unity{
		server:   server,
		username: username,
		password: password,
		token:    token,
		client:   client,
		logger:   logger,
	}, nil
}

func (unity *Unity) Request(method string, URI string, fields string, filter string, payload interface{}, result interface{}) error {
	requestParams := fmt.Sprintf("method: %s, URI: %s, fields: %s, filter: %s, payload: %#v", method, URI, fields, filter, payload)
	unity.Log("debug", requestParams)
	if method == "" || URI == "" {
		return errors.New("method, or URI is missed.")
	}

	url := "https://" + unity.server + URI

	var req *http.Request
	if payload != nil {
		payloadJson, _ := json.Marshal(payload)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(payloadJson))
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
	unity.Log("debug", fmt.Sprintf("Request: %s", string(reqDetails)))

	resp, err := unity.client.Do(req)
	if err != nil {
		unity.Log("error", err.Error())
		return err
	}
	defer resp.Body.Close()

	respDetails, _ := httputil.DumpResponse(resp, true)
	unity.Log("debug", string(respDetails))
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		respRaw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if result != nil {
			json.Unmarshal(respRaw, result)
		}

		return nil
	} else {
		unity.Log("error", "Fail to perform the request")
		return fmt.Errorf("Request Fails: %s", requestParams)
	}
}

func (unity *Unity) Log(levelStr string, message string) {
	level := getLogLevel(levelStr)
	switch level {
	case log.PanicLevel:
		unity.logger.Panic(message)
	case log.FatalLevel:
		unity.logger.Fatal(message)
	case log.ErrorLevel:
		unity.logger.Error(message)
	case log.WarnLevel:
		unity.logger.Warn(message)
	case log.InfoLevel:
		unity.logger.Info(message)
	case log.DebugLevel:
		unity.logger.Debug(message)
	case log.TraceLevel:
		unity.logger.Trace(message)
	}
}

func (unity *Unity) Destroy() error {
	unity.Log("debug", "Logout")
	err := unity.Request("POST", "/api/types/loginSessionInfo/action/logout", "", "", nil, nil)
	return err
}
