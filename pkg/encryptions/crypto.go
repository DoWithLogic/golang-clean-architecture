package encryptions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type Crypto struct {
	key string
}

func NewCrypto(key string) *Crypto {
	return &Crypto{key: key}
}

// EncodeSHA1HMACBase64 : encrypt to SHA1HMAC input key, data String. Output to String in Base64 format
func (c *Crypto) EncodeSHA1HMACBase64(data ...string) string {
	return c.EncodeBASE64(c.ComputeSHA1HMAC(data...))
}

// EncodeSHA1HMAC : encrypt to SHA1HMAC input key, data String. Output to String in Base16/Hex format
func (c *Crypto) EncodeSHA1HMAC(data ...string) string {
	return fmt.Sprintf("%x", c.ComputeSHA1HMAC(data...))
}

// ComputeSHA1HMAC : encrypt to SHA1HMAC input key, data String. Output to String
func (c *Crypto) ComputeSHA1HMAC(data ...string) string {
	h := hmac.New(sha1.New, []byte(c.key))
	for _, v := range data {
		io.WriteString(h, v)
	}
	return string(h.Sum(nil))
}

func (c *Crypto) EncodeSHA256HMACBase64(data ...string) string {
	return c.EncodeBASE64(c.ComputeSHA256HMAC(data...))
}

func (c *Crypto) EncodeSHA256HMAC(data ...string) string {
	return fmt.Sprintf("%x", c.ComputeSHA256HMAC(data...))
}

func (c *Crypto) ComputeSHA256HMAC(data ...string) string {
	h := hmac.New(sha256.New, []byte(c.key))
	for _, v := range data {
		io.WriteString(h, v)
	}
	return string(h.Sum(nil))
}

func (c *Crypto) EncodeSHA512HMACBase64(data ...string) string {
	return c.EncodeBASE64(c.ComputeSHA512HMAC(data...))
}

func (c *Crypto) EncodeSHA512HMAC(data ...string) string {
	return fmt.Sprintf("%x", c.ComputeSHA512HMAC(data...))
}

func (c *Crypto) ComputeSHA512HMAC(data ...string) string {
	h := hmac.New(sha512.New, []byte(c.key))
	for _, v := range data {
		io.WriteString(h, v)
	}
	return string(h.Sum(nil))
}

// EncodeMD5 : encrypt to MD5 input string, output to string
func (c *Crypto) EncodeMD5(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *Crypto) EncodeMD5Base64(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	// return EncodeBASE64(hex.EncodeToString(h.Sum(nil)))
	return base64.StdEncoding.EncodeToString((h.Sum(nil)))
}

// EncodeBASE64 : Encrypt to Base64. Input string, output string
func (c *Crypto) EncodeBASE64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

// DecodeBASE64 : Decrypt Base64. Input string, output string
func (c *Crypto) DecodeBASE64(text string) (string, error) {
	byt, err := base64.StdEncoding.DecodeString(text)
	return string(byt), err
}

// EncodeBASE64URL : Encrypt to Base64URL. Input string, output text
func (c *Crypto) EncodeBASE64URL(text string) string {
	return base64.URLEncoding.EncodeToString([]byte(text))
}

// EncodeAES256GCM encrypts the given plaintext string using AES-256-GCM.
// keyHex should be a 64-character hex string (32 bytes key).
func (c *Crypto) EncodeAES256GCM(keyHex, plaintext string) (string, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", err
	}
	if len(key) != 32 {
		return "", errors.New("AES key must be 32 bytes (256 bits)")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}

// EncodeSHA256: Encrypt to SHA256. input string, output text
func (c *Crypto) EncodeSHA256(text string) string {
	h := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", h)
}

// EncodeSHA512 Encrypt to SHA512. input string, output text
func (c *Crypto) EncodeSHA512(text string) string {
	h := sha512.Sum512([]byte(text))
	return fmt.Sprintf("%x", h)
}
