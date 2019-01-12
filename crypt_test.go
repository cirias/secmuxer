package main

import (
	"fmt"
	"testing"
)

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
