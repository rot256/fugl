package main

import (
	"fmt"
	"github.com/rot256/fugl"
	"io/ioutil"
	"net/http"
	"os"
)

func createDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func createState() *ServerState {
	// read public key
	var state ServerState
	key, err := ioutil.ReadFile(config.Canary.KeyFile)
	if err != nil {
		logFatal("Unable to load public key from:", config.Canary.KeyFile)
	}
	state.canaryKeyArmor = string(key)
	state.canaryKey, err = fugl.PGPLoadPublicKey(key)
	if err != nil {
		logFatal("Unable to parse PGP key:", err)
	}

	// load latest proof
	err = createDir(config.Canary.Store)
	if err != nil {
		logFatal("Unable to create store:", err)
	}
	state.latestProof, err = fugl.LoadLatestProof(config.Canary.Store)
	if err != nil {
		logFatal("Failed to load latest proof")
	}

	// parse latest proof
	if state.latestProof != "" {
		state.latestCanary, err = fugl.OpenProof(state.canaryKey, state.latestProof)
		if err != nil {
			logFatal("Failed to load latest canary:", err.Error())
		}
	}
	state.storeDir = config.Canary.Store
	return &state
}

func buildHandler() http.Handler {
	state := createState()
	handler := http.NewServeMux()
	handler.Handle("/submit", &SubmitHandler{state: state})
	handler.Handle("/latest", &LatestHandler{state: state})
	handler.Handle("/status", &StatusHandler{state: state})
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
