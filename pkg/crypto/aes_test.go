package crypto

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := []byte("01234567890123456789012345678901") // 32 bytes
	plaintext := "secret_api_key_123"

	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	if encrypted == plaintext {
		t.Fatal("Encrypted text should not match plaintext")
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Expected %s, got %s", plaintext, decrypted)
	}
}

func TestDecryptInvalidKey(t *testing.T) {
	key1 := []byte("01234567890123456789012345678901")
	key2 := []byte("12345678901234567890123456789012")
	plaintext := "secret"

	encrypted, _ := Encrypt(plaintext, key1)
	_, err := Decrypt(encrypted, key2)

	if err == nil {
		t.Fatal("Decryption should fail with wrong key")
	}
}
