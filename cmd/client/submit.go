package main

/*
import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/clearsign"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type CanaryTime time.Time

type Canary struct {
	Publish  CanaryTime `json:"publish"`  // Time of publication
	Version  int64      `json:"version"`  // Struct version
	Message  string     `json:"message"`  // Optional notification
	Previous string     `json:"previous"` // Hash of previous canary
	Deadline CanaryTime `json:"deadline"` // New deadline
	Nonce    string     `json:"nonce"`    // Random nonce
}

const (
	URL_CANARY_STATUS = "/canary/status"
	URL_CANARY_LATEST = "/canary/latest"
	URL_CANARY_SUBMIT = "/canary/submit"
	POST_FIELD        = "canary"
	CANARY_MAX_SIZE   = 1024 * 100
	CANARY_VERSION    = 0
	CANARY_PUBLISH    = time.Hour * 24
	CANARY_DEADLINE   = time.Hour * 24 * 2
)

func (t CanaryTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02"))
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
	date, err := time.Parse("2006-01-02", str[1:len(str)-1])
	if err != nil {
		return err
	}
	*t = CanaryTime(date)
	return nil
}

func CanaryParse(key *openpgp.Entity, proof []byte) (*Canary, error) {
	// Verify using public key
	block, rest := clearsign.Decode(proof)
	if len(rest) > 0 {
		return nil, errors.New("Proof contains junk data :(")
	}

	// Verify signature
	keyring := make(openpgp.EntityList, 1)
	keyring[0] = key
	content := bytes.NewReader(block.Bytes)
	_, err := openpgp.CheckDetachedSignature(keyring, content, block.ArmoredSignature.Body)
	if err != nil {
		return nil, errors.New("Signature on proof is invalid")
	}

	// Unserialize body of proof into canary
	var canary Canary
	err = json.Unmarshal(block.Bytes, &canary)
	return &canary, err
}

func CanaryFetchLatestRaw(address string) ([]byte, error) {
	// Make request
	url := address + URL_CANARY_LATEST
	resp, err := TorHTTPGet(url)
	if err != nil {
		return []byte{}, err
	}

	// Handle canaries not enabled
	if resp.StatusCode == http.StatusNoContent {
		return []byte{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Server returned status %s", resp.Status)
		return []byte{}, errors.New(msg)
	}

	// Read full proof
	defer resp.Body.Close()
	proof, err := ioutil.ReadAll(resp.Body)
	return proof, err
}

func CanaryFetchLatest(key *openpgp.Entity, address string) (*Canary, string, error) {
	proof, err := CanaryFetchLatestRaw(address)
	if err != nil {
		return nil, "", err
	} else if len(proof) == 0 {
		return nil, "", nil
	}
	hash := hex.EncodeToString(SecureHash(proof))

	// Parse
	canary, err := CanaryParse(key, proof)
	return canary, hash, err
}

func CanaryPush(address string, proof string) error {
	data := url.Values{}
	data.Add(POST_FIELD, proof)
	url := address + URL_CANARY_SUBMIT
	resp, err := TorHTTPPostForm(url, data)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New("Submission failed: " + resp.Status)
		}
		return errors.New("Submission failed: " + resp.Status + " with:\n" + string(msg))
	}
	return nil
}

func CanarySubmitNew(key *openpgp.Entity, address string, message string, enable bool) error {
	// Load old canary
	var canary Canary
	fmt.Println("Fetching newest canary from:", address)
	oldCanary, oldHash, err := CanaryFetchLatest(key, address)
	if err != nil {
		return errors.New("Failed to fetch newest canary: " + err.Error())
	}
	if oldCanary == nil && !enable {
		return errors.New("Inital canary, use enable flag!")
	}

	// Create new canary
	canary.Message = message
	canary.Previous = oldHash
	canary.Nonce = GetRandStr(32)
	canary.Version = CANARY_VERSION
	canary.Deadline = CanaryTime(time.Now().Add(CANARY_DEADLINE + CANARY_PUBLISH))
	canary.Publish = CanaryTime(time.Now().Add(CANARY_PUBLISH))

	// If the deadline is violated, publish earlier
	if oldCanary != nil && time.Time(oldCanary.Deadline).Before(time.Time(canary.Publish)) {
		canary.Publish = oldCanary.Deadline
	}

	// Produce proof
	ser, err := json.MarshalIndent(canary, "", "    ")
	proof, err := PGPSign(key, ser)
	fmt.Println("Created a new proof:")
	fmt.Println(proof)

	// Push new canary to server
	return CanaryPush(address, proof)
}
*/
