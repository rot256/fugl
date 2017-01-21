package main

import (
	"fmt"
	"net/http"
)

const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

func buildHandler() http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("/submit", HandleSubmit)
	handler.HandleFunc("/latest", HandleLatest)
	handler.HandleFunc("/status", HandleStatus)
	return handler
}

func main() {
	// build server
	bind := fmt.Sprintf(
		"%s:%d",
		config.Server.Address,
		config.Server.Port)
	server := &http.Server{
		ReadTimeout:    config.Server.TimeoutRead,
		WriteTimeout:   config.Server.TimeoutWrite,
		Addr:           bind,
		Handler:        buildHandler(),
		MaxHeaderBytes: 1 << 20,
	}

	// run http(s) server
	var err error
	if config.Server.CertFile == "" || config.Server.KeyFile == "" {
		logInfo("Starting HTTP server on:", bind)
		err = server.ListenAndServe()
	} else {
		logInfo("Starting HTTPS server on:", bind)
		err = server.ListenAndServeTLS(
			config.Server.CertFile,
			config.Server.KeyFile)
	}
	logFatal("Server terminated with:", err)
}
