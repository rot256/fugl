package main

import (
	"golang.org/x/net/proxy"
	"net/http"
)

func CreateHttpClient(proxyAddr string) (*http.Client, error) {
	if proxyAddr == "" {
		return &http.Client{}, nil
	}
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	httpTransport := &http.Transport{
		Dial: dialer.Dial,
	}
	return &http.Client{Transport: httpTransport}, err
}
