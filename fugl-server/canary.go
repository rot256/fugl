package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/clearsign"
	"golang.org/x/crypto/openpgp/packet"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

/* contains helper functions and initalization routines for handlers
 *
 */

type CanaryTime time.Time

const (
	CanaryVersion      = 0
	CanaryTimeFormat   = "2006-01-02"
	ProofFileExtension = ".sig"
	ProofFileName      = "proof-%s-%s" + ProofFileExtension
)

var (
	CanaryKey      *openpgp.Entity
	CanaryKeyArmor string
	LatestProof    string
	LatestCanary   *Canary
)

func (t CanaryTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(CanaryTimeFormat))
	return []byte(stamp), nil
}

func (t *CanaryTime) UnmarshalJSON(val []byte) error {
	str := string(val)
	if len(str) < 2 {
		return errors.New("Time field too short")
	}
	if str[0] != '"' || str[len(str)-1] != '"' {
		return errors.New("Time must be json string type")
	}
	date, err := time.Parse(CanaryTimeFormat, str[1:len(str)-1])
	if err != nil {
		return err
	}
	*t = CanaryTime(date)
	return nil
}

func PGPLoadPublicKey(key string) (*openpgp.Entity, error) {
	block, err := armor.Decode(bytes.NewReader([]byte(key)))
	if err != nil {
		return nil, err
	} else if block.Type != openpgp.PublicKeyType {
		return nil, errors.New("Not a PGP public key")
	}
	return openpgp.ReadEntity(packet.NewReader(block.Body))
}

func CreateDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func LoadLatestProof(dir string) (string, error) {
	// find newest proof
	var proofFile string
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		// check if proof file
		if file.IsDir() {
			logWarning("Directory found in store")
			continue
		}
		if !strings.HasSuffix(file.Name(), ProofFileExtension) {
			logWarning("Non-proof file in store:", file.Name())
			continue
		}
		proofFile = file.Name()
		logDebug("Store holds:", file.Name())
	}

	// read proof
	if proofFile != "" {
		logInfo("Loading proof from:", proofFile)
		proof, err := ioutil.ReadFile(path.Join(dir, proofFile))
		return string(proof), err
	}
	return "", nil
}

func VerifyProof(publicKey *openpgp.Entity, proof string) (*Canary, error) {
	// Parse clear signature
	block, rest := clearsign.Decode([]byte(proof))
	if len(rest) > 0 {
		return nil, errors.New("Proof contains junk")
	}

	// Verify signature
	keyring := make(openpgp.EntityList, 1)
	keyring[0] = publicKey
	content := bytes.NewReader(block.Bytes)
	_, err := openpgp.CheckDetachedSignature(keyring, content, block.ArmoredSignature.Body)
	if err != nil {
		return nil, errors.New("Invalid signature")
	}

	// load inner JSON structure
	var canary Canary
	err = json.Unmarshal(block.Bytes, &canary)
	if err != nil {
		return nil, errors.New("Unable to parse inner canary structure")
	}
	return &canary, nil
}

func SaveToStore(proof string, when time.Time) error {
	hash := Sha256StringToHex(proof)
	date := time.Time(when).Format(CanaryTimeFormat)
	fileName := fmt.Sprintf(ProofFileName, date, hash)
	filePath := path.Join(config.Canary.Store, fileName)
	return ioutil.WriteFile(filePath, []byte(proof), 0600)
}

func initCanary() {
	// read public key
	key, err := ioutil.ReadFile(config.Canary.KeyFile)
	if err != nil {
		logFatal("Unable to load public key from:", config.Canary.KeyFile)
	}
	CanaryKeyArmor = string(key)
	CanaryKey, err = PGPLoadPublicKey(CanaryKeyArmor)
	if err != nil {
		logFatal("Unable to parse PGP key:", err)
	}

	// load latest canary
	CreateDir(config.Canary.Store)
	LatestProof, err = LoadLatestProof(config.Canary.Store)
	if err != nil {
		logFatal("Failed to load latest proof")
	}
	logDebug("Loaded proof:", LatestProof)
	if LatestProof != "" {
		LatestCanary, err = VerifyProof(CanaryKey, LatestProof)
		if err != nil {
			logFatal("Failed to load latest canary:", err.Error())
		}
	}
}
