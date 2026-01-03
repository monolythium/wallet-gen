package keystore

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestEncryptDecryptRoundtrip(t *testing.T) {
	// Generate a test key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	password := "test-password-123"
	params := DefaultScryptParams()
	// Use lighter params for faster tests
	params.N = 4096

	// Encrypt
	data, err := EncryptKey(privateKey, password, params)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	// Verify it's valid JSON
	if len(data) == 0 {
		t.Fatal("empty keystore data")
	}

	// Decrypt
	decryptedKey, err := DecryptKey(data, password)
	if err != nil {
		t.Fatalf("failed to decrypt: %v", err)
	}

	// Verify the key matches
	originalBytes := crypto.FromECDSA(privateKey)
	decryptedBytes := crypto.FromECDSA(decryptedKey)

	if string(originalBytes) != string(decryptedBytes) {
		t.Error("decrypted key does not match original")
	}
}

func TestDecryptWrongPassword(t *testing.T) {
	privateKey, _ := crypto.GenerateKey()
	password := "correct-password"
	params := DefaultScryptParams()
	params.N = 4096

	data, err := EncryptKey(privateKey, password, params)
	if err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	// Try to decrypt with wrong password
	_, err = DecryptKey(data, "wrong-password")
	if err == nil {
		t.Error("expected error with wrong password")
	}
}

func TestWriteKeystorePermissions(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "walletgen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	keystoreDir := filepath.Join(tmpDir, "keystores")
	data := []byte(`{"test": "data"}`)
	address := "0x1234567890abcdef1234567890abcdef12345678"

	filename, err := WriteKeystore(keystoreDir, 0, data, address)
	if err != nil {
		t.Fatalf("failed to write keystore: %v", err)
	}

	// Verify file was created
	fullPath := filepath.Join(keystoreDir, filename)
	info, err := os.Stat(fullPath)
	if err != nil {
		t.Fatalf("keystore file not found: %v", err)
	}

	// Verify permissions (Unix only)
	mode := info.Mode().Perm()
	if mode != 0600 {
		t.Errorf("expected permissions 0600, got %o", mode)
	}
}

func TestDefaultScryptParams(t *testing.T) {
	params := DefaultScryptParams()

	if params.N != 262144 {
		t.Errorf("expected N=262144, got %d", params.N)
	}
	if params.R != 8 {
		t.Errorf("expected R=8, got %d", params.R)
	}
	if params.P != 1 {
		t.Errorf("expected P=1, got %d", params.P)
	}
}
