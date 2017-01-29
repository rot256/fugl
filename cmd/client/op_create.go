package main

import (
	"fmt"
	"github.com/rot256/fugl"
	"io/ioutil"
	"os"
	"time"
)

func requiredFlagsCreate(flags Flags) {
	var msg string
	if flags.PrivateKey == "" {
		msg += requiredArgument(FlagNamePrivateKey)
	}
	if flags.Output == "" {
		msg += requiredArgument(FlagNameOutput)
	}
	if flags.Store == "" {
		msg += requiredArgument(FlagNameStore)
	}
	if flags.Expire == time.Duration(0) {
		msg += requiredArgument(FlagNameExpire)
	}
	if msg != "" {
		fmt.Fprintf(os.Stderr, "%s", msg)
		os.Exit(EXIT_INVALID_ARGUMENTS)
	}
}

func operationCreate(flags Flags) {
	// verify supplied flags
	requiredFlagsCreate(flags)

	// load latest proof from store
	oldProof, err := fugl.LoadLatestProof(flags.Store)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load latest proof: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}

	// load private key
	skData, err := ioutil.ReadFile(flags.PrivateKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read private key: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}
	sk, err := fugl.PGPLoadPrivateKey(skData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}

	// load message (optional)
	var message string
	if flags.Message != "" {
		tmp, err := ioutil.ReadFile(flags.Message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read message: %s", err.Error())
			os.Exit(EXIT_FILE_READ_ERROR)
		}
		message = string(tmp)
	}

	// create canary
	canary := fugl.Canary{
		Version:  fugl.CanaryVersion,
		Message:  message,
		Previous: fugl.HashString(oldProof),
		Nonce:    fugl.GetRandStr(fugl.CanaryNonceSize),
		Deadline: fugl.CanaryTime(time.Now().Add(flags.Expire)),
	}

	// sign canary, producing proof
	proof, err := fugl.SealProof(sk, canary)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to sign canary: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}

	// write to output
	err = ioutil.WriteFile(flags.Output, []byte(proof), 0555)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write proof to file: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}
	fmt.Println("Wrote new proof to:", flags.Output)
}
