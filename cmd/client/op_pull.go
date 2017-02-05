package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func requiredFlagsPull(flags Flags) {
	var opt FlagOpt
	opt.Required(FlagNameAddress, flags.Address != "")
	opt.Required(FlagNameOutput, flags.Output != "")
	opt.Optional(FlagNameProxy, flags.Proxy != "")
	opt.Check()
}

func operationPull(flags Flags) {
	requiredFlagsPull(flags)

	// create (proxied) client
	client, err := CreateHttpClient(flags.Proxy)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed connect to proxy: %s", err.Error())
		os.Exit(EXIT_BAD_PROXY)
	}

	// do http request
	resp, err := client.Get(flags.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed connect to address: %s", err.Error())
		os.Exit(EXIT_CONNECTION_FAILURE)
	}
	defer resp.Body.Close()

	// read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read response body: %s", err.Error())
		os.Exit(EXIT_CONNECTION_FAILURE)
	}

	// write to file
	err = ioutil.WriteFile(flags.Output, body, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write proof to file: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}
	fmt.Println("Saved to:", flags.Output)
}
