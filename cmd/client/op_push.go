package main

import (
	"fmt"
	"github.com/rot256/fugl"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func requiredFlagsPush(flags Flags) {
	var opt FlagOpt
	opt.Required(FlagNameAddress, flags.Address != "")
	opt.Required(FlagNameInput, flags.Input != "")
	opt.Optional(FlagNameProxy, flags.Proxy != "")
	opt.Check()
}

func operationPush(flags Flags) {
	requiredFlagsPush(flags)

	// create (proxied) client
	client, err := func(addr string) (*http.Client, error) {
		if addr == "" {
			return &http.Client{}, nil
		}
		dialer, err := proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
		httpTransport := &http.Transport{
			Dial: dialer.Dial,
		}
		return &http.Client{Transport: httpTransport}, err
	}(flags.Proxy)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed connect to proxy: %s", err.Error())
		os.Exit(EXIT_BAD_PROXY)
	}

	// read input file
	content, err := ioutil.ReadFile(flags.Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read input file: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}

	// post
	form := url.Values{}
	form.Add(fugl.SERVER_SUBMIT_FIELD_NAME, string(content))
	resp, err := client.PostForm(flags.Address, form)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to remote server: %s", err.Error())
		os.Exit(EXIT_CONNECTION_FAILURE)
	}

	// check for server error message
	if resp.StatusCode != http.StatusNoContent {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Submission failed: %s", resp.Status)
			os.Exit(EXIT_CONNECTION_FAILURE)
		}
		fmt.Fprintf(os.Stderr, "Submission failed %s with: '%s'", resp.Status, string(msg))
		os.Exit(EXIT_CONNECTION_FAILURE)
	}
}
