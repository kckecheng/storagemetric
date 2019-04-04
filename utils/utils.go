package utils

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Logger global log object
var Logger *log.Logger

// PrettyPrint Print json to console elegantly, which is helpful for degbugging
func PrettyPrint(data interface{}, pre string, post string) {
	dataJSON, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Printf("Data cannot be reformatted, the original data is:\n %#v\n", data)
	} else {
		if pre != "" {
			fmt.Printf("%s\n", pre)
		}
		fmt.Printf("%s\n", dataJSON)
		if post != "" {
			fmt.Printf("%s\n", post)
		}
	}
}

// GetLoginOptions Get login information
func GetLoginOptions() (string, string, string, error) {
	var server, username, password string
	flag.StringVar(&server, "server", "", "Server address")
	flag.StringVar(&username, "username", "admin", "Server login username, admin as default")
	flag.StringVar(&password, "password", "", "Server login password")
	flag.Parse()

	if server == "" || username == "" || password == "" {
		flag.PrintDefaults()
		return "", "", "", errors.New("server, username, and password must all be specified")
	}
	return server, username, password, nil
}

//  URL Compose a URL based on protocol, server address, URI, etc.
func URL(proto string, fqdn string, port string, uri ...string) string {
	var url string

	if proto == "" {
		url = fqdn
	} else {
		url = proto + "://" + fqdn
	}

	if port != "" {
		url += ":" + port
	}

	for _, uri := range uri {
		url += uri
	}

	return url
}

// GetLogLevel Get the log level based on input string
// Valid string: panic fatal error warning info debug trace
func GetLogLevel(logStr string) log.Level {
	level, err := log.ParseLevel(logStr)
	if err != nil {
		return log.DebugLevel
	}
	return level
}

// InitLogger Init global log object
func InitLogger(filename string, levelStr string) {
	var defaultLogFile = "storagemetric.log"
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
		level = GetLogLevel(envLogStr)
	} else {
		level = GetLogLevel(levelStr)
	}
	logger.SetLevel(level)

	// if level == log.DebugLevel {
	//   logger.SetReportCaller(true)
	// }
	Logger = logger
}

// Log log message based on log level
func Log(levelStr string, message string) {
	level := GetLogLevel(levelStr)
	switch level {
	case log.PanicLevel:
		Logger.Panic(message)
	case log.FatalLevel:
		Logger.Fatal(message)
	case log.ErrorLevel:
		Logger.Error(message)
	case log.WarnLevel:
		Logger.Warn(message)
	case log.InfoLevel:
		Logger.Info(message)
	case log.DebugLevel:
		Logger.Debug(message)
	case log.TraceLevel:
		Logger.Trace(message)
	}
}
