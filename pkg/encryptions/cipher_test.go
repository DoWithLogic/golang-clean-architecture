package encryptions_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/encryptions"
)

func TestCipherEncryptDecrypt(t *testing.T) {
	// AES-256 key = 32 bytes
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}

	cipher, err := encryptions.NewAES256GCM(key)
	if err != nil {
		t.Fatal(err)
	}

	plaintext := []byte("This is a test message.")

	// Encrypt
	ciphertext, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Decrypt
	decrypted, err := cipher.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("decrypted text mismatch: expected %q, got %q", plaintext, decrypted)
	}
}

func TestDecryptInvalidCiphertextLength(t *testing.T) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}

	cipher, err := encryptions.NewAES256GCM(key)
	if err != nil {
		t.Fatal(err)
	}

	// Too short ciphertext (less than nonce size)
	invalidCiphertext := make([]byte, cipher.AEAD().NonceSize()-1)
	_, err = cipher.Decrypt(invalidCiphertext)

	if err == nil {
		t.Error("expected error for invalid ciphertext length, got none")
	} else if err.Error() != "ciphertext too short" {
		t.Errorf("expected 'ciphertext too short', got %q", err.Error())
	}
}
