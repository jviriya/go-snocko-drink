package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func CalculateHMAC(message, secretKey string) string {
	key := []byte(secretKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
