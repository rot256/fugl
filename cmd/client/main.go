package main

import (
	"fmt"
	"os"
)

const (
	SCRYPT_N  = 2 << 20
	SCRYPT_R  = 8
	SCRYPT_P  = 1
	SALT_SIZE = 16
)

func main() {
	// parse flags
	flags := parseFlags()

	// handle different operations
	switch flags.Operation {
	case "create":
		operationCreate(flags)
	case "verify":
		operationVerify(flags)
	case "push":
		operationPush(flags)
	case "pull":
		operationPull(flags)
	case "":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Invalid operation: %s\n", flags.Operation)
		os.Exit(EXIT_INVALID_OPERATION)
	}
}
