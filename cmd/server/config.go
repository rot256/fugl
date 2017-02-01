package main

import (
	"github.com/BurntSushi/toml"
	"time"
)

type ConfigLogging struct {
	File  string `toml:"file"`
	Level string `toml:"level"`
}

type ConfigServer struct {
	Port             uint16
	Address          string
	TimeoutRead      time.Duration
	TimeoutWrite     time.Duration
	CertFile         string `toml:"cert_file"`     // tls: certificate
	KeyFile          string `toml:"key_file"`      // tls: private key
	EnableViewSubmit bool   `toml:"enable_submit"` // enable submit view
	EnableViewStatus bool   `toml:"enable_status"` // enable status view
	EnableViewLatest bool   `toml:"enable_latest"` // enable latest view
	EnableViewGetKey bool   `toml:"enable_getkey"` // enable get key view
}

type ConfigCanary struct {
	OnFailure string `toml:"on_failure"` // command on failure
	KeyFile   string `toml:"key_file"`   // load key from this file
	Store     string `toml:"store"`      // directory for storing canaries
}

type Config struct {
	Logging ConfigLogging `toml:"logging"` // log settings
	Server  ConfigServer  `toml:"server"`  // http server settings
	Canary  ConfigCanary  `toml:"canary"`  // canary settings
}

func loadConfig() (Config, error) {
	var config Config
	_, err := toml.DecodeFile(*FlagConfigPath, &config)
	return config, err
}
