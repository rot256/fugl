package main

import (
	"errors"
	"fmt"
	"github.com/rot256/fugl"
	"golang.org/x/crypto/openpgp"
	"io/ioutil"
	"os"
)

/* Validates and adds a canary to the store
 * and updates the state of the store
 */

func requiredFlagsAdd(flags Flags) {
	var opt FlagOpt
	opt.Required(FlagNamePublicKey, flags.PublicKey != "")
	opt.Required(FlagNameInput, flags.Input != "")
	opt.Required(FlagNameStore, flags.Store != "")
	opt.Check()
}

func operationAdd(flags Flags) {
	requiredFlagsAdd(flags)

	// load latest proof from store
	oldProof, err := fugl.LoadLatestProof(flags.Store)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load latest proof: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}

	// load new proof (input)
	newProof, err := ioutil.ReadFile(flags.Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to input proof: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}

	// load public key
	pk, err := func() (*openpgp.Entity, error) {
		skData, err := ioutil.ReadFile(flags.PublicKey)
		if err != nil {
			return nil, err
		}
		return fugl.PGPLoadPublicKey(skData)
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read public key: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}

	// validate new proof
	newCanary, err := func(proof string, oldProof string, pk *openpgp.Entity) (*fugl.Canary, error) {
		// parse old proof
		oldCanary, err :=
			func(proof string, pk *openpgp.Entity) (*fugl.Canary, error) {
				if proof == "" {
					return nil, nil
				}
				return fugl.OpenProof(pk, oldProof)
			}(string(oldProof), pk)
		if err != nil {
			return nil, errors.New("Failed to read old proof from store")
		}

		// parse and verify signature
		canary, err := fugl.OpenProof(pk, proof)
		if err != nil {
			return nil, err
		}

		// check version field
		if canary.Version != fugl.CanaryVersion {
			return nil, errors.New("Invalid canary version field")
		}

		// validate against old proof
		if oldCanary != nil {
			if !oldCanary.Deadline.Time().After(canary.Deadline.Time()) {
				return nil, errors.New("New canary deadline must be after old")
			}
			if fugl.HashString(oldProof) != canary.Previous {
				return nil, errors.New("Previous hash does not match latest proof in store")
			}
		}
		return canary, nil
	}(string(newProof), string(oldProof), pk)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to validate new proof: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}

	// save new proof to store
	err = fugl.SaveToDirectory(string(newProof), flags.Store, newCanary.Deadline.Time())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save proof to store: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}
	if newCanary.Message != "" {
		fmt.Println(newCanary.Message)
	}
	fmt.Println("New deadline:", newCanary.Deadline.Time().Format(fugl.CanaryTimeFormat))
}
