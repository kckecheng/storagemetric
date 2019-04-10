package utils

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Logger global log object
var Logger *log.Logger

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
