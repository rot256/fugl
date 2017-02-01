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

func createState(config Config) *ServerState {
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

func buildHandler(config Config) (http.Handler, *ServerState) {
	state := createState(config)
	handler := http.NewServeMux()
	if config.Server.EnableViewSubmit {
		logInfo("Enable view: Submit")
		handler.Handle(fugl.SERVER_SUBMIT_PATH, &SubmitHandler{state: state})
	}
	if config.Server.EnableViewLatest {
		logInfo("Enable view: Latest")
		handler.Handle(fugl.SERVER_LATEST_PATH, &LatestHandler{state: state})
	}
	if config.Server.EnableViewStatus {
		logInfo("Enable view: Status")
		handler.Handle(fugl.SERVER_STATUS_PATH, &StatusHandler{state: state})
	}
	if config.Server.EnableViewGetKey {
		logInfo("Enable view: GetKey")
		handler.Handle(fugl.SERVER_GETKEY_PATH, &GetKeyHandler{state: state})
	}
	return handler, state
}

func main() {
	// initalize logger
	config, err := loadConfig()
	if err != nil {
		logFatal("Unable to load config")
	}
	initLogging(config)

	// build handler and server state
	handler, state := buildHandler(config)
	go actionRunner(config.Canary.OnFailure, state)

	// build server
	bind := fmt.Sprintf(
		"%s:%d",
		config.Server.Address,
		config.Server.Port)
	server := &http.Server{
		ReadTimeout:    config.Server.TimeoutRead,
		WriteTimeout:   config.Server.TimeoutWrite,
		Addr:           bind,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
	}

	// run http(s) server
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
