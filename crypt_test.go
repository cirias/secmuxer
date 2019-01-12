package main

import (
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
}

func TestCryptWithWrongPassphrase(t *testing.T) {
	plaintext := "HelloWorld"
	passphrase := "secret"
	wrongPassphrase := "wrong"
	ciphertext, err := encrypt([]byte(plaintext), passphrase)
	if err != nil {
		t.Error(err)
	}

	decrypted, err := decrypt(ciphertext, wrongPassphrase)
	if err == nil {
		t.Errorf("Decrypted text does not match original:\n[%s]\n[%s]\n", plaintext, decrypted)
	}
}

func TestDeriveKey(t *testing.T) {
	s1 := make([]byte, saltSize)

	k1 := deriveKey("aaa", s1)
	k2 := deriveKey("aaa", s1)

	if *k1 != *k2 {
		t.Errorf("Same password and salt should generate same key")
	}

	k1 = deriveKey("aaa", s1)
	k2 = deriveKey("bbb", s1)

	if *k1 == *k2 {
		t.Errorf("Different password with same key should generate different key")
	}

	s2 := make([]byte, saltSize)
	s2[0] = 1
	k1 = deriveKey("aaa", s1)
	k2 = deriveKey("aaa", s2)

	if *k1 == *k2 {
		t.Errorf("Same password with different key should generate different key")
	}
}
