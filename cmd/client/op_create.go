package main

import (
	"errors"
	"fmt"
	"github.com/rot256/fugl"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/ssh/terminal"
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
	sk, err := func() (*openpgp.Entity, error) {
		skData, err := ioutil.ReadFile(flags.PrivateKey)
		if err != nil {
			return nil, err
		}
		sk, err := fugl.PGPLoadPrivateKey(skData)
		if err != nil {
			return nil, err
		}
		if sk.PrivateKey.Encrypted {
			fmt.Println("Private key encrypted, please enter passphrase:")
			passwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return nil, err
			}
			err = sk.PrivateKey.Decrypt(passwd)
			if err != nil {
				return nil, errors.New("Failed to decrypt key")
			}
		}
		return sk, err
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read private key: %s", err.Error())
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
	err = ioutil.WriteFile(flags.Output, []byte(proof), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write proof to file: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}
	fmt.Println("Wrote new proof to:", flags.Output)
}
