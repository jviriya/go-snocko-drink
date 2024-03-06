package cryptographer

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/chacha20"
)

func EncryptData(data string, key []byte) (string, error) {
	dataByte := []byte(data)
	var nonce [chacha20.NonceSize]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", err
	}

	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce[:])
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(data))
	cipher.XORKeyStream(ciphertext, dataByte)

	ciphertext = append(nonce[:], ciphertext...)

	encodeToString := base64.StdEncoding.EncodeToString(ciphertext)

	return encodeToString, nil
}

func DecryptData(data string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < chacha20.NonceSize {
		return "", fmt.Errorf("ciphertext is too short")
	}

	var nonce [chacha20.NonceSize]byte
	copy(nonce[:], ciphertext[:chacha20.NonceSize])

	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce[:])
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(ciphertext)-chacha20.NonceSize)
	cipher.XORKeyStream(plaintext, ciphertext[chacha20.NonceSize:])

	return string(plaintext), nil
}

func GenerateKey() ([]byte, error) {
	keyLength := 16 // Chacha20 key size is 256 bits (32 bytes)
	key, err := generateRandomKey(keyLength)
	if err != nil {
		fmt.Println("Error generating random key:", err)
		return nil, err
	}
	fmt.Printf("Derived key: %x\n", key)
	return key, nil
}

func generateRandomKey(keyLength int) ([]byte, error) {
	key := make([]byte, keyLength)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
