package main

import (
	"errors"
	"io"
	"log"
	"os"
	"time"
)

const (
	TIME_FORMAT = "2006-01-02 15:04:05 MST"
)

type logWriter struct {
	output io.Writer
}

func (writer logWriter) Write(msg []byte) (int, error) {
	logged := []byte(time.Now().Format(TIME_FORMAT) + " ")
	logged = append(logged, msg...)
	return writer.output.Write(logged)
}

var (
	LogEnableDebug   bool
	LogEnableInfo    bool
	LogEnableWarning bool
	LogEnableError   bool
)

func initLogging() {
	// Expand logging level
	switch config.Logging.Level {
	case "debug":
		LogEnableDebug = true
		LogEnableInfo = true
		LogEnableWarning = true
		LogEnableError = true
	case "info":
		LogEnableDebug = false
		LogEnableInfo = true
		LogEnableWarning = true
		LogEnableError = true
	case "warning":
		LogEnableDebug = false
		LogEnableInfo = false
		LogEnableWarning = true
		LogEnableError = true
	case "error":
		LogEnableDebug = false
		LogEnableInfo = false
		LogEnableWarning = false
		LogEnableError = true
	default:
		logFatal("Config: log_level must be \"error\", \"warning\", \"info\" or \"debug\"")
	}

	// Multiplex output
	var output io.Writer = os.Stdout
	if config.Logging.File != "" {
		logFile, err := os.OpenFile(config.Logging.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Panicln("Failed to open log:", config.Logging.File, err)
		}
		output = io.MultiWriter(output, logFile)
	}
	log.SetFlags(0)
	log.SetOutput(logWriter{output})
}

func logPrefix(s interface{}, v []interface{}) {
	args := make([]interface{}, 0, len(v)+1)
	args = append(args, s)
	args = append(args, v...)
	log.Println(args...)
}

func logDebug(v ...interface{}) {
	if LogEnableDebug {
		logPrefix("[DEBUG] :", v)
	}
}

func logInfo(v ...interface{}) {
	if LogEnableInfo {
		logPrefix("[INFO] :", v)
	}
}

func logWarning(v ...interface{}) {
	if LogEnableWarning {
		logPrefix("[WARNING] :", v)
	}
}

func logError(v ...interface{}) {
	if LogEnableError {
		logPrefix("[ERROR] :", v)
	}
}

func logFatal(v ...interface{}) {
	logPrefix("[FATAL] :", v)
	panic(errors.New("The server experienced a critical error, see the log for details"))
}

func logCheck(err error) {
	if err != nil {
		logFatal(err)
	}
}
