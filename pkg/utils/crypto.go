package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func Encrypt(text, secretKey, salt string) string {
	plaintext := []byte(text)
	key := []byte(secretKey)
	saltBytes := []byte(salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	cipherText := make([]byte, aes.BlockSize+len(plaintext))

	// Use the salt in the initialization vector
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(append(saltBytes, cipherText...))
}

func Decrypt(encryptedText, secretKey, salt string) string {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		panic(err)
	}

	key := []byte(secretKey)
	saltBytes := []byte(salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext[len(saltBytes):])
}
