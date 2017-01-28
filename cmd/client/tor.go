package main

/*
import (
	"fmt"
	"golang.org/x/net/proxy"
	"log"
	"net/http"
	"net/url"
)

// Proxy dialer
func TorNewHTTPClient() http.Client {
	// insecure connection
	if FlagInsecure {
		return http.Client{}
	}

	// connect though tor
	dialer, err := proxy.SOCKS5("tcp", FlagProxyAddress, nil, proxy.Direct)
	if err != nil {
		log.Panicln("Could not connect to proxy", err)
	}
	httpTransport := &http.Transport{
		Dial: dialer.Dial,
	}
	return http.Client{Transport: httpTransport}
}

// Make a get request though the tor proxy
func TorHTTPGet(url string) (*http.Response, error) {
	if FlagDebug {
		fmt.Println("[DEBUG]: Get " + url)
	}
	client := TorNewHTTPClient()
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// Post a form though tor
func TorHTTPPostForm(url string, data url.Values) (*http.Response, error) {
	if FlagDebug {
		fmt.Println("[DEBUG]: Post " + url)
	}
	client := TorNewHTTPClient()
	resp, err := client.PostForm(url, data)
	return resp, err
}
*/
