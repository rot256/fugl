package fugl

import (
	"bytes"
	"errors"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/clearsign"
	"golang.org/x/crypto/openpgp/packet"
)

func PGPLoadPrivateKey(key []byte) (*openpgp.Entity, error) {
	block, err := armor.Decode(bytes.NewReader([]byte(key)))
	if err != nil {
		return nil, err
	} else if block.Type != openpgp.PrivateKeyType {
		return nil, errors.New("Not a OpenPGP public key")
	}
	return openpgp.ReadEntity(packet.NewReader(block.Body))
}

func PGPLoadPublicKey(key []byte) (*openpgp.Entity, error) {
	block, err := armor.Decode(bytes.NewReader([]byte(key)))
	if err != nil {
		return nil, err
	} else if block.Type != openpgp.PublicKeyType {
		return nil, errors.New("Not a OpenPGP private key")
	}
	return openpgp.ReadEntity(packet.NewReader(block.Body))
}

func PGPSign(entity *openpgp.Entity, message []byte) (string, error) {
	if entity.PrivateKey == nil {
		return "", errors.New("invalid private key")
	}

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
	if entity == nil {
		return nil, errors.New("invalid public key")
	}

	// parse clear signature
	block, rest := clearsign.Decode(signature)
	if len(rest) > 0 {
		return nil, errors.New("Proof contains junk")
	}
	if block == nil {
		return nil, errors.New("Unable to read pgp block")
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
