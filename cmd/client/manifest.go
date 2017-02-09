package main

import (
	"github.com/BurntSushi/toml"
)

type Manifest struct {
	Author      string   `toml:"author"`      // supposed author of canary
	Delta       int64    `toml:"delta"`       // time in seconds
	Promises    []string `toml:"promises"`    // list of promises (for machines)
	Description string   `toml:"description"` // content of human readable portion
}

func ParseManifest(path string) (Manifest, error) {
	var manifest Manifest
	_, err := toml.DecodeFile(path, &manifest)
	return manifest, err
}
