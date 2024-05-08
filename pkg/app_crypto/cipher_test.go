package app_crypto_test

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_crypto"
)

func TestCipherEncryptDecrypt(t *testing.T) {
	key := make([]byte, aes.BlockSize)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}

	cipher, err := app_crypto.NewES256(key)
	if err != nil {
		t.Fatal(err)
	}

	plaintext := []byte("This is a test message.")

	// Test encryption
	ciphertext, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}

	// Test decryption
	decryptedText, err := cipher.Decrypt(ciphertext)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure the decrypted text matches the original plaintext
	if !bytes.Equal(decryptedText, plaintext) {
		t.Errorf("Decrypted text doesn't match original plaintext. Expected %s, got %s", plaintext, decryptedText)
	}
}

func TestInvalidCiphertextLength(t *testing.T) {
	key := make([]byte, aes.BlockSize)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}

	cipher, err := app_crypto.NewES256(key)
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to decrypt invalid ciphertext length
	invalidCiphertext := make([]byte, aes.BlockSize-1)
	_, err = cipher.Decrypt(invalidCiphertext)

	if err == nil {
		t.Error("Expected error for invalid ciphertext length, but got none.")
	} else {
		expectedErrorMsg := "ciphertext too short"
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMsg, err.Error())
		}
	}
}

func TestInvalidBlockSize(t *testing.T) {
	key := make([]byte, aes.BlockSize)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}

	cipher, err := app_crypto.NewES256(key)
	if err != nil {
		t.Fatal(err)
	}

	// Attempt to decrypt ciphertext with invalid block size
	invalidCiphertext := make([]byte, aes.BlockSize+1)
	_, err = cipher.Decrypt(invalidCiphertext)

	if err == nil {
		t.Error("Expected error for invalid block size, but got none.")
	} else {
		expectedErrorMsg := "ciphertext is not a multiple of the block size"
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMsg, err.Error())
		}
	}
}
