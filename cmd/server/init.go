package main

import (
	"flag"
	"log"
	"os"
)

var FlagConfigPath = flag.String("config", "config.toml", "path to config file")

func init() {
	flag.Parse()
	log.SetFlags(0)
	log.SetOutput(logWriter{os.Stdout})
}
