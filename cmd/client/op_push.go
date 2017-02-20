package main

import (
	"fmt"
	"github.com/rot256/fugl"
	"io/ioutil"
	"net/http"
	"net/url"
)

func requiredFlagsPush(flags Flags) {
	var opt FlagOpt
	opt.Required(FlagNameAddress, flags.Address != "")
	opt.Required(FlagNameProof, flags.Proof != "")
	opt.Optional(FlagNameProxy, flags.Proxy != "")
	opt.Check()
}

func operationPush(flags Flags) {
	requiredFlagsPush(flags)

	// create final url
	addr, err := createURL(flags.Address, fugl.SERVER_SUBMIT_PATH)
	if err != nil {
		exitError(EXIT_INVALID_ADDRESS, "Failed to parse address %s", err.Error())
	}

	// create (proxied) client
	client, err := CreateHttpClient(flags.Proxy)
	if err != nil {
		exitError(EXIT_BAD_PROXY, "Failed connect to proxy: %s", err.Error())
	}

	// read input file
	content, err := ioutil.ReadFile(flags.Proof)
	if err != nil {
		exitError(EXIT_FILE_READ_ERROR, "Failed to read input file: %s", err.Error())
	}

	// post
	form := url.Values{}
	form.Add(fugl.SERVER_SUBMIT_FIELD_NAME, string(content))
	resp, err := client.PostForm(addr, form)
	if err != nil {
		exitError(EXIT_CONNECTION_FAILURE, "Failed to connect to remote server: %s", err.Error())
	}

	// check for server error message
	if resp.StatusCode != http.StatusNoContent {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			exitError(EXIT_CONNECTION_FAILURE, "Submission failed: %s", resp.Status)
		}
		exitError(EXIT_CONNECTION_FAILURE, "Submission failed %s with: '%s'", resp.Status, string(msg))
	}
	fmt.Println("Successfully pushed new proof to server")
}
