package main

import (
	"fmt"
	"github.com/rot256/fugl"
	"io/ioutil"
	"net/http"
	"os"
)

func requiredFlagsPull(flags Flags) {
	var opt FlagOpt
	opt.Required(FlagNameAddress, flags.Address != "")
	opt.Required(FlagNameProof, flags.Proof != "")
	opt.Optional(FlagNameProxy, flags.Proxy != "")
	opt.Check()
}

func operationPull(flags Flags) {
	requiredFlagsPull(flags)

	// create final url
	addr, err := createURL(flags.Address, fugl.SERVER_LATEST_PATH)
	if err != nil {
		exitError(EXIT_INVALID_ADDRESS, "Failed to parse address %s", err.Error())
	}

	// create (proxied) client
	client, err := CreateHttpClient(flags.Proxy)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed connect to proxy: %s", err.Error())
		os.Exit(EXIT_BAD_PROXY)
	}

	// do http request
	resp, err := client.Get(addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed connect to address: %s", err.Error())
		os.Exit(EXIT_CONNECTION_FAILURE)
	}
	defer resp.Body.Close()

	// check if body is present
	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("No canary available")
		return
	}
	if resp.StatusCode != http.StatusOK {
		exitError(EXIT_HTTP_UNEXPECTED_STATUS, "Server returned unexpected status code: %s", resp.Status)
	}

	// read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		exitError(EXIT_CONNECTION_FAILURE, "Failed to read response body: %s", err.Error())
	}

	// write to file
	err = ioutil.WriteFile(flags.Proof, body, 0644)
	if err != nil {
		exitError(EXIT_FILE_WRITE_ERROR, "Failed to write proof to file: %s", err.Error())
	}
	fmt.Println("Saved to:", flags.Proof)
}
