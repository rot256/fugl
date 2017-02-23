package main

import (
	"fmt"
	"github.com/rot256/fugl"
	"golang.org/x/crypto/openpgp"
	"io/ioutil"
	"time"
)

/* Validates and adds a canary to the store
 * and updates the state of the store
 */

func requiredFlagsVerify(flags Flags) {
	var opt FlagOpt
	opt.Required(FlagNamePublicKey, flags.PublicKey != "")
	opt.Required(FlagNameProof, flags.Proof != "")
	opt.Check()
}

func operationVerify(flags Flags) {
	requiredFlagsVerify(flags)

	// load new proof (input)
	proof, err := ioutil.ReadFile(flags.Proof)
	if err != nil {
		exitError(EXIT_FILE_READ_ERROR, "Failed to input proof: %s", err.Error())
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
		exitError(EXIT_FILE_READ_ERROR, "Failed to read public key: %s", err.Error())
	}

	// validate new proof
	canary, description, err := fugl.OpenProof(pk, string(proof))
	if err != nil {
		exitError(EXIT_INVALID_SIGNATURE, "Failed to validate signature on proof: %s", err.Error())
	}

	// verify fields
	err = fugl.CheckCanaryFormat(canary, time.Now())
	if err != nil {
		exitError(EXIT_INVALID_CANARY, "Failed to validate canary fields: %s", err.Error())
	}
	fmt.Println("Author:", canary.Author)
	fmt.Println("Expires:", canary.Expiry.String())
	fmt.Println("Description:\n" + description)
}
