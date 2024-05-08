package app_crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Cipher represents a cryptographic cipher using the AES (Advanced Encryption Standard) algorithm in CBC (Cipher Block Chaining) mode.
type Cipher struct {
	block cipher.Block // The underlying block cipher.
}

// NewES256 creates a new Cipher instance with the specified key for AES encryption.
// Parameters:
//   - key: The encryption key.
//
// Returns:
//   - *Cipher: A pointer to the newly created Cipher instance.
//   - error: An error if the cipher creation fails.
func NewES256(key []byte) (*Cipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &Cipher{block: block}, nil
}

// Decrypt decrypts the given ciphertext using the AES-CBC decryption algorithm.
// Parameters:
//   - ciphertext: The ciphertext to decrypt.
//
// Returns:
//   - []byte: The decrypted plaintext.
//   - error: An error if decryption fails.
func (c *Cipher) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(c.block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Remove padding
	padding := int(ciphertext[len(ciphertext)-1])
	return ciphertext[:len(ciphertext)-padding], nil
}

// Encrypt encrypts the given plaintext using the AES-CBC encryption algorithm.
// Parameters:
//   - plaintext: The plaintext to encrypt.
//
// Returns:
//   - []byte: The encrypted ciphertext.
//   - error: An error if encryption fails.
func (c *Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	paddedPlaintext := append(plaintext, bytes.Repeat([]byte{byte(padding)}, padding)...)

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(paddedPlaintext))
	copy(ciphertext[:aes.BlockSize], iv)

	mode := cipher.NewCBCEncrypter(c.block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], paddedPlaintext)

	return ciphertext, nil
}
