package main

import (
	"net/url"
)

func createURL(address string, path string) (string, error) {
	url, err := url.Parse(address)
	if err != nil {
		return "", err
	}

	// behaves strange if omitted
	// create a ticket if you disagree
	if url.Host == "" {
		url.Host = url.Path
		url.Path = ""
	}

	// set default scheme and path
	if url.Path == "" {
		url.Path = path
	}
	if url.Scheme == "" {
		url.Scheme = "https"
	}
	return url.String(), nil
}
