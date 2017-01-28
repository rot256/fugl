package main

import (
	"fmt"
	"github.com/rot256/fugl"
	"os"
	"time"
)

func requiredFlagsCreate(flags Flags) {
	var msg string
	if flags.PrivateKey == "" {
		msg += requiredArgument(FlagNamePrivateKey)
	}
	if flags.PublicKey == "" {
		msg += requiredArgument(FlagNamePublicKey)
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
	proof, err := fugl.LoadLatestProof(flags.Store)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load latest proof: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}
	if flags.Debug {
		fmt.Println(proof) // debug
	}

	// load private key

	// load message (optional)

	// create canary
	var canary fugl.Canary
	canary.Version = fugl.CanaryVersion
	canary.Message = ""
	canary.Previous = fugl.HashString(proof)
	canary.Nonce = fugl.GetRandStr(fugl.CanaryNonceSize)
	canary.Deadline = fugl.CanaryTime(time.Now().Add(flags.Expire))
	if flags.Debug {
		fmt.Println(canary)
	}

	// sign canary, produce proof

	// write to output
}
