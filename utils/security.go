package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashSha256(stringToHash string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(stringToHash))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash, nil
}

func BytesToHex(array []byte) (string, error) {
	return hex.EncodeToString(array), nil
}
