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
