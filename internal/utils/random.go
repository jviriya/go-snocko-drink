package utils

import (
	"math/rand"
	"time"
	"unsafe"
)

var src = rand.NewSource(time.Now().UnixNano())

const (
	userBytes     = "abcdefghijklmnopqrstuvwxyz1234567890"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandomNonSensitive(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(userBytes) {
			b[i] = userBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

const (
	passwordBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	passwordIdxBits = 6                      // 6 bits to represent a password index
	passwordIdxMask = 1<<passwordIdxBits - 1 // All 1-bits, as many as passwordIdxBits
	passwordIdxMax  = 63 / passwordIdxBits   // # of password indices fitting in 63 bits
)

func RandomSensitive(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), passwordIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), passwordIdxMax
		}
		if idx := int(cache & passwordIdxMask); idx < len(passwordBytes) {
			b[i] = passwordBytes[idx]
			i--
		}
		cache >>= passwordIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

const (
	uidBytes   = "1234567890"
	uidIdxBits = 4
	uidIdxMask = 1<<uidIdxBits - 1
	uidIdxMax  = 15 / uidIdxBits
)

func RandomUID(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), uidIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), uidIdxMax
		}
		if idx := int(cache & uidIdxMask); idx < len(uidBytes) {
			b[i] = uidBytes[idx]
			i--
		}
		cache >>= uidIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
