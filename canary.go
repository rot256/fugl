package fugl

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

/* contains helper functions and initalization routines for handlers
 */

const (
	CanaryVersion       = 0
	CanaryTimeFormat    = time.RFC3339
	ProofFileTimeFormat = "20060102150405" // must be a valid filename and sortable
	CanaryNonceSize     = 32
	ProofFileExtension  = ".proof"
	ProofFileName       = "proof-%s-%s" + ProofFileExtension
)

func CheckCanary(new *Canary, old *Canary, now time.Time) error {
	err := CheckCanaryFormat(new, now)
	if err != nil {
		return err
	}
	return CheckCanaryPrevious(new, old)
}

func CheckCanaryFormat(canary *Canary, now time.Time) error {
	if canary.Version != CanaryVersion {
		return errors.New("Unsupported canary version")
	}
	if len(canary.Nonce) != CanaryNonceSize {
		return errors.New(fmt.Sprintf("Nonce must be %d characters long", CanaryNonceSize))
	}
	if now.After(canary.Expiry.Time()) {
		return errors.New("Canary has expired")
	}
	if canary.Author == "" {
		return errors.New("Author field is empty")
	}
	if canary.Creation.Time().After(now) {
		return errors.New("Creation time cannot be in the future (canary not valid yet)")
	}
	return nil
}

func CheckCanaryPrevious(new *Canary, old *Canary) error {
	if old == nil {
		return nil
	}
	if old.Final {
		return errors.New("Current canary is final")
	}
	if old.Creation.Time().After(new.Creation.Time()) {
		return errors.New("Current canary has creation after new")
	}
	if !new.Expiry.Time().After(old.Expiry.Time()) {
		return errors.New("New canary does not have expiry time after old")
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
			return "", errors.New("Directory found in store")
		}
		if !strings.HasSuffix(file.Name(), ProofFileExtension) {
			return "", errors.New("Non-proof file in store: " + file.Name())
		}
		proofFile = file.Name()
	}

	// read proof
	if proofFile != "" {
		proof, err := ioutil.ReadFile(path.Join(dir, proofFile))
		return string(proof), err
	}
	return "", nil
}

func SaveToDirectory(proof string, dir string, when time.Time) error {
	hash := HashString(proof)
	date := time.Time(when).Format(ProofFileTimeFormat)
	fileName := fmt.Sprintf(ProofFileName, date, hash)
	filePath := path.Join(dir, fileName)
	return ioutil.WriteFile(filePath, []byte(proof), 0600)
}
