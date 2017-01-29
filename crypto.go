package fugl

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

var RAND_ALPHABET = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_")

func HashString(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
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
