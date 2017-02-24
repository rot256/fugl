package fugl

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"testing"
)

var (
	pair keypair
)

func init() {
	p, err := newKeypair()
	if err != nil {
		panic(fmt.Sprintf("error creating pgp keys = %v", err))
	}

	if len(p.public) == 0 || len(p.private) == 0 {
		panic(fmt.Sprint("public (len=%d) or private (len=%d) keys don't contain data", len(p.public), len(p.private)))
	}

	pair = p
}

func TestPGP__LoadKeypair(t *testing.T) {
	// private
	priv, err := PGPLoadPrivateKey([]byte(pair.private))
	if err != nil {
		t.Fatalf("error reading private key, err=%v", err)
	}
	if priv == nil {
		t.Fatalf("private key is read as nil")
	}

	// public
	pub, err := PGPLoadPublicKey([]byte(pair.public))
	if err != nil {
		t.Fatalf("error reading public key, err=%v", err)
	}
	if pub == nil {
		t.Fatalf("public key is read as nil")
	}

}

func TestPGP__LoadInvalidKeypair(t *testing.T) {
	// private
	priv, err := PGPLoadPrivateKey([]byte(""))
	if err == nil {
		t.Fatal("no error reading invalid key")
	}
	if priv != nil {
		t.Fatal("given a non-nil private key")
	}

	// public
	pub, err := PGPLoadPublicKey([]byte(""))
	if err == nil {
		t.Fatal("no error reading invalid pub key")
	}
	if pub != nil {
		t.Fatal("given a non-nil public key")
	}
}

func TestPGP__SignAndVerify(t *testing.T) {
	// ignore errors, they should be picked up by __LoadKeypair
	pub, _ := PGPLoadPublicKey([]byte(pair.public))
	priv, _ := PGPLoadPrivateKey([]byte(pair.private))

	// sign and verify a message
	message := []byte("this is a test message")
	sig, err := PGPSign(priv, message)
	if err != nil {
		t.Fatalf("error signing message, err=%v", err)
	}
	if len(sig) == 0 {
		t.Fatal("empty signature created")
	}

	block, err := PGPVerify(pub, []byte(sig))
	if err != nil {
		t.Fatalf("error verifying signature, err=%v", err)
	}
	if block == nil || len(block.Bytes) == 0 {
		t.Fatal("verify block is empty")
	}
}

func TestPGP__InvalidSignAndVerify(t *testing.T) {
	pub, _ := PGPLoadPublicKey([]byte(pair.public))

	// sign and verify a message
	message := []byte("this is a test message")
	sig, err := PGPSign(pub, message)
	if err == nil {
		t.Fatal("expected an error when signing with wrong key")
	}
	if len(sig) != 0 {
		t.Fatal("should not get signature on error")
	}

	block, err := PGPVerify(pub, []byte(""))
	if err == nil {
		t.Fatal("expected an error verifying an invalid signature")
	}
	if block != nil {
		t.Fatal("should not get a block on invalid verify")
	}
}

// helper functions to generate keys

type keypair struct {
	public, private string
}
func (p *keypair) empty() bool {
	return len(p.public) == 0 || len(p.private) == 0
}

func newKeypair() (keypair, error) {
	p := keypair{}

	entity, err := openpgp.NewEntity("", "", "", nil)
	entity.Subkeys = entity.Subkeys[:0]
	if len(entity.Identities) != 1 {
		return p, nil
	}

	// Serialize private key
	var secretArmor bytes.Buffer
	secArmIn, err := armor.Encode(&secretArmor, openpgp.PrivateKeyType, nil)
	if err != nil {
		return p, err
	}
	err = entity.SerializePrivate(secArmIn, nil)
	if err != nil {
		return p, err
	}
	err = secArmIn.Close()
	if err != nil {
		return p, err
	}
	p.private = secretArmor.String()

	// Serialize public key
	var publicArmor bytes.Buffer
	pubArmIn, err := armor.Encode(&publicArmor, openpgp.PublicKeyType, nil)
	if err != nil {
		return p, err
	}
	err = entity.Serialize(pubArmIn)
	if err != nil {
		return p, err
	}
	err = pubArmIn.Close()
	if err != nil {
		return p, err
	}
	p.public = publicArmor.String()

	return p, nil
}
