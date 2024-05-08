package app_crypto_test

import (
	"testing"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/app_crypto"
)

func TestCrypto_EncodeSHA1HMACBase64(t *testing.T) {
	c := app_crypto.NewCrypto("secretKey")
	result := c.EncodeSHA1HMACBase64("data1", "data2")

	// You need to replace the expectedValue with the actual HMAC SHA1 Base64 value based on your secret key and input data.
	expectedValue := "BC9mO9N3TM9DeXyopI7eFXJ78pM="
	if result != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, result)
	}
}

func TestCrypto_EncodeSHA256HMAC(t *testing.T) {
	c := app_crypto.NewCrypto("secretKey")
	result := c.EncodeSHA256HMAC("data1", "data2")

	// You need to replace the expectedValue with the actual HMAC SHA256 value based on your secret key and input data.
	expectedValue := "9cdb1c43d56c85451dfb630f66748778380ed7e7d114fab185f1b82adb1b658d"
	if result != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, result)
	}
}

func TestCrypto_EncodeSHA512HMACBase64(t *testing.T) {
	c := app_crypto.NewCrypto("secretKey")
	result := c.EncodeSHA512HMACBase64("data1", "data2")

	// You need to replace the expectedValue with the actual HMAC SHA512 Base64 value based on your secret key and input data.
	expectedValue := "28PiVYGlVn8XQrfF0rC9LXRmZcG1VYTK+akgZn3LBuxXRcwDGdaTKB05KfhehZXV1gO43Ie+s4NfHu17Bw16zA=="
	if result != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, result)
	}
}

func TestCrypto_EncodeMD5(t *testing.T) {
	c := app_crypto.NewCrypto("secretKey")
	result := c.EncodeMD5("text")

	// You need to replace the expectedValue with the actual MD5 value based on your input text.
	expectedValue := "1cb251ec0d568de6a929b520c4aed8d1"
	if result != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, result)
	}
}

func TestCrypto_EncodeMD5Base64(t *testing.T) {
	c := app_crypto.NewCrypto("secretKey")
	result := c.EncodeMD5Base64("text")

	// You need to replace the expectedValue with the actual MD5 Base64 value based on your input text.
	expectedValue := "HLJR7A1WjeapKbUgxK7Y0Q=="
	if result != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, result)
	}
}

func TestCrypto_EncodeBASE64(t *testing.T) {
	c := app_crypto.NewCrypto("secretKey")
	result := c.EncodeBASE64("text")

	// You need to replace the expectedValue with the actual Base64 value based on your input text.
	expectedValue := "dGV4dA=="
	if result != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, result)
	}
}

func TestCrypto_DecodeBASE64(t *testing.T) {
	c := app_crypto.NewCrypto("secretKey")
	encodedText := c.EncodeBASE64("text")
	result, err := c.DecodeBASE64(encodedText)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if result != "text" {
		t.Errorf("Expected %s, got %s", "text", result)
	}
}

func TestCrypto_EncodeBASE64URL(t *testing.T) {
	c := app_crypto.NewCrypto("secretKey")
	result := c.EncodeBASE64URL("text")

	// You need to replace the expectedValue with the actual Base64URL value based on your input text.
	expectedValue := "dGV4dA=="
	if result != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, result)
	}
}

func TestCrypto_EncodeSHA256(t *testing.T) {
	c := app_crypto.NewCrypto("secretKey")
	result := c.EncodeSHA256("text")

	// You need to replace the expectedValue with the actual SHA256 value based on your input text.
	expectedValue := "982d9e3eb996f559e633f4d194def3761d909f5a3b647d1a851fead67c32c9d1"
	if result != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, result)
	}
}

func TestCrypto_EncodeSHA512(t *testing.T) {
	c := app_crypto.NewCrypto("secretKey")
	result := c.EncodeSHA512("text")

	// You need to replace the expectedValue with the actual SHA512 value based on your input text.
	expectedValue := "eaf2c12742cb8c161bcbd84b032b9bb98999a23282542672ca01cc6edd268f7dce9987ad6b2bc79305634f89d90b90102bcd59a57e7135b8e3ceb93c0597117b"
	if result != expectedValue {
		t.Errorf("Expected %s, got %s", expectedValue, result)
	}
}
