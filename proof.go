package fugl

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"strings"
)

func OpenProof(entity *openpgp.Entity, proof string) (*Canary, string, error) {
	// parse and verify signature
	block, err := PGPVerify(entity, []byte(proof))
	if err != nil {
		return nil, "", err
	}

	// scan for seperator
	start := 0
	lines := strings.Split(string(block.Bytes), "\n")
	for ; start < len(lines); start++ {
		if strings.TrimRight(lines[start], "\n\r") == CANARY_SEPERATOR {
			break
		}
	}
	if start == len(lines) {
		return nil, "", errors.New("Unable to find canary seperator")
	}

	// eat seperator and empty lines
	des := strings.Join(lines[:start-1], "\n")
	for start = start + 1; start < len(lines); start++ {
		if strings.TrimRight(lines[start], "\n\r") != "" {
			break
		}
	}

	// load JSON structure
	var canary Canary
	ser := strings.Join(lines[start:], "\n")
	err = json.Unmarshal([]byte(ser), &canary)
	if err != nil {
		return nil, "", errors.New("Unable to parse json structure")
	}
	return &canary, des, nil
}

func SealProof(entity *openpgp.Entity, canary Canary, description string) (string, error) {
	// serialize canary
	ser, err := json.MarshalIndent(canary, "", "    ")
	if err != nil {
		return "", err
	}

	// add serperator and sign
	var inner string
	if description == "" {
		inner = fmt.Sprintf("%s\n%s", CANARY_SEPERATOR, string(ser))
	} else {
		inner = fmt.Sprintf("%s\n%s\n%s", description, CANARY_SEPERATOR, string(ser))
	}
	return PGPSign(entity, []byte(inner))
}
