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
	var opt FlagOpt
	opt.Required(FlagNamePrivateKey, flags.PrivateKey != "")
	opt.Required(FlagNameManifest, flags.Manifest != "")
	opt.Required(FlagNameProof, flags.Proof != "")
	opt.Check()
}

func operationCreate(flags Flags) {
	// verify supplied flags
	requiredFlagsCreate(flags)

	// load manifest
	manifest, err := ParseManifest(flags.Manifest)
	if err != nil {
		exitError(EXIT_FILE_READ_ERROR, "Failed to load manifest %s", err.Error())
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

	// create canary
	now := time.Now()
	expire := now.Add(time.Duration(manifest.Delta) * time.Second)
	canary := fugl.Canary{
		Version:  fugl.CanaryVersion,
		Author:   manifest.Author,
		Creation: fugl.CanaryTime(now),
		Expiry:   fugl.CanaryTime(expire),
		Nonce:    fugl.GetRandStr(fugl.CanaryNonceSize),
		Final:    manifest.Final,
	}

	// sign canary, producing proof
	proof, err := fugl.SealProof(sk, canary, manifest.Description)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to sign canary: %s", err.Error())
		os.Exit(EXIT_FILE_READ_ERROR)
	}

	// write to output
	err = ioutil.WriteFile(flags.Proof, []byte(proof), 0644)
	if err != nil {
		exitError(EXIT_FILE_WRITE_ERROR, "Failed to write proof to file: %s", err.Error())
	}
	fmt.Println("Saved new proof to:", flags.Proof)
	if manifest.Final {
		fmt.Println("WARNING: This canary is final!")
	}
}
