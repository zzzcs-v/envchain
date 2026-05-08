package encrypt

import (
	"strings"
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	e := NewEncryptor("supersecret")
	plaintext := "MY_SECRET_VALUE"

	encoded, err := e.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if encoded == plaintext {
		t.Fatal("expected ciphertext to differ from plaintext")
	}

	decoded, err := e.Decrypt(encoded)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decoded != plaintext {
		t.Errorf("expected %q, got %q", plaintext, decoded)
	}
}

func TestEncrypt_ProducesUniqueOutputs(t *testing.T) {
	e := NewEncryptor("passphrase")
	a, _ := e.Encrypt("value")
	b, _ := e.Encrypt("value")
	if a == b {
		t.Error("expected different ciphertexts due to random nonce")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	e1 := NewEncryptor("correct")
	e2 := NewEncryptor("wrong")

	encoded, err := e1.Encrypt("secret")
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = e2.Decrypt(encoded)
	if err == nil {
		t.Fatal("expected decryption to fail with wrong passphrase")
	}
	if !strings.Contains(err.Error(), "wrong passphrase") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	e := NewEncryptor("key")
	_, err := e.Decrypt("!!!not-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64 input")
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	e := NewEncryptor("key")
	// base64 of a 3-byte slice — shorter than AES-GCM nonce (12 bytes)
	_, err := e.Decrypt("AAAA")
	if err == nil {
		t.Fatal("expected error for too-short ciphertext")
	}
}
