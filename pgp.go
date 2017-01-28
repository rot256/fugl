package fugl

import (
	"bytes"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/clearsign"
)

/* Helper function for OpenPGP */

func OpenProof(entity *openpgp.Entity, proof string) (*Canary, error) {
	// check signature
	block, err := PGPVerify(entity, []byte(proof))
	if err != nil {
		return nil, err
	}

	// load inner JSON structure
	var canary Canary
	err = json.Unmarshal(block.Bytes, &canary)
	if err != nil {
		return nil, errors.New("Unable to parse inner canary structure")
	}
	return &canary, nil
}

func SealProof(entity *openpgp.Entity, canary Canary) (string, error) {
	// serialize and sign
	inner, err := json.MarshalIndent(canary, "", "    ")
	if err != nil {
		return "", err
	}
	return PGPSign(entity, inner)
}

func PGPSign(entity *openpgp.Entity, message []byte) (string, error) {
	// create signature writer
	var outSig bytes.Buffer
	writer, err := clearsign.Encode(&outSig, entity.PrivateKey, nil)
	if err != nil {
		return "", err
	}

	// sign entire message and flush
	_, err = writer.Write(message)
	if err != nil {
		return "", err
	}
	err = writer.Close()
	if err != nil {
		return "", err
	}
	return outSig.String(), err
}

func PGPVerify(entity *openpgp.Entity, signature []byte) (*clearsign.Block, error) {
	// parse clear signature
	block, rest := clearsign.Decode([]byte(signature))
	if len(rest) > 0 {
		return nil, errors.New("Proof contains junk")
	}

	// verify signature
	keyring := make(openpgp.EntityList, 1)
	keyring[0] = entity
	content := bytes.NewReader(block.Bytes)
	_, err := openpgp.CheckDetachedSignature(keyring, content, block.ArmoredSignature.Body)
	if err != nil {
		return nil, errors.New("Invalid signature")
	}
	return block, nil
}

// todo: remove?
func PGPNewKey() (secretStr string, publicStr string, err error) {
	entity, err := openpgp.NewEntity("", "", "", nil) // opt: config
	entity.Subkeys = entity.Subkeys[:0]
	if len(entity.Identities) != 1 {
		return "", "", errors.New("Multiple identities for entity")
	}

	// Serialize private key
	var secretArmor bytes.Buffer
	secArmIn, err := armor.Encode(&secretArmor, openpgp.PrivateKeyType, nil)
	if err != nil {
		return
	}
	err = entity.SerializePrivate(secArmIn, nil) // opt: config
	if err != nil {
		return
	}
	err = secArmIn.Close()
	if err != nil {
		return
	}
	secretStr = secretArmor.String()

	// Serialize public key
	var publicArmor bytes.Buffer
	pubArmIn, err := armor.Encode(&publicArmor, openpgp.PublicKeyType, nil)
	if err != nil {
		return
	}
	err = entity.Serialize(pubArmIn)
	if err != nil {
		return
	}
	err = pubArmIn.Close()
	if err != nil {
		return
	}
	publicStr = publicArmor.String()
	return
}
