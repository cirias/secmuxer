package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const saltSize = 16 // bytes

// Openssl's way of derive key from passphrase is considered too weak, see:
// https://security.stackexchange.com/questions/29106/openssl-recover-key-and-iv-by-passphrase
// Use pbkdf2 instead see:
// https://github.com/riverrun/comeonin/wiki/Choosing-the-password-hashing-algorithm
func deriveKey(passphrase string, salt []byte) *[32]byte {
	var key [32]byte
	dk := pbkdf2.Key([]byte(passphrase), salt, 4096, len(key), sha1.New)
	copy(key[:], dk)
	return &key
}

func encrypt(plaintext []byte, passphrase string) (cipherStr string, err error) {
	salt := make([]byte, saltSize)
	_, err = io.ReadFull(rand.Reader, salt)
	if err != nil {
		return "", err
	}

	key := deriveKey(passphrase, salt)

	encrypted, err := encryptWithKey(plaintext, key)
	if err != nil {
		return "", err
	}

	cipherBytes := make([]byte, 0, saltSize+len(encrypted))
	cipherBytes = append(cipherBytes, salt...)
	cipherBytes = append(cipherBytes, encrypted...)

	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

func decrypt(cipherStr string, passphrase string) (plaintext []byte, err error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cipherStr)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < saltSize {
		return nil, fmt.Errorf("Malformed encrypted data")
	}

	salt := ciphertext[:saltSize]
	key := deriveKey(passphrase, salt)
	return decryptWithKey(ciphertext[saltSize:], key)
}

// Encrypt and decrypt with key are copied from https://github.com/gtank/cryptopasta
func encryptWithKey(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decryptWithKey(ciphertext []byte, key *[32]byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}
