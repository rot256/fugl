package fugl

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

var RAND_ALPHABET = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")

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

func GetRandBytes(n int) []byte {
	b := make([]byte, n)
	i, err := rand.Read(b)
	if i != n || err != nil {
		panic(err)
	}
	return b
}

func GetRandStr(n int) string {
	str := make([]rune, n)
	for i, v := range GetRandBytes(n) {
		str[i] = RAND_ALPHABET[int(v)%len(RAND_ALPHABET)]
	}
	return string(str)
}
