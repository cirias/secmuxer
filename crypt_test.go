package main

import (
	"fmt"
	"testing"
)

/*
func TestBcrypt(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("test"), 14)
	fmt.Printf("Bcrypt hash size %d\nHash: [%s]", len(hash), hash)
	if err != nil {
		t.Error(err)
	}
}
*/

func TestCrypt(t *testing.T) {
	plaintext := "HelloWorld"
	passphrase := "secret"
	ciphertext, err := encrypt([]byte(plaintext), passphrase)
	if err != nil {
		t.Error(err)
	}

	decrypted, err := decrypt(ciphertext, passphrase)
	if err != nil {
		t.Error(err)
	}
	if plaintext != string(decrypted) {
		t.Errorf("Decrypted text does not match original:\n[%s]\n[%s]\n", plaintext, decrypted)
	}

	fmt.Printf("Successful\n")
}
