package main

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

const (
	DefaultProtocol = "https"
)

func createURL(address string, p string) (string, error) {
	// If we're missing the uri scheme, prepend it
	if !strings.HasPrefix(address, "http") {
		address = fmt.Sprintf("%s://%s", DefaultProtocol, address)
	}

	url, err := url.Parse(address)
	if err != nil {
		return "", err
	}

	url.Path = path.Join(url.Path, p)

	return url.String(), nil
}
