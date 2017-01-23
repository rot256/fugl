package fugl

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/clearsign"
	"golang.org/x/crypto/openpgp/packet"
)

func HashString(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func LoadPublicKey(key string) (*openpgp.Entity, error) {
	block, err := armor.Decode(bytes.NewReader([]byte(key)))
	if err != nil {
		return nil, err
	} else if block.Type != openpgp.PublicKeyType {
		return nil, errors.New("Not a OpenPGP public key")
	}
	return openpgp.ReadEntity(packet.NewReader(block.Body))
}

func VerifyProof(publicKey *openpgp.Entity, proof string) (*Canary, error) {
	// parse clear signature
	block, rest := clearsign.Decode([]byte(proof))
	if len(rest) > 0 {
		return nil, errors.New("Proof contains junk")
	}

	// verify signature
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
