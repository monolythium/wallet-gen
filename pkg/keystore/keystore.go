// Package keystore provides Ethereum Keystore V3 encryption.
package keystore

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
)

// ScryptParams holds scrypt KDF parameters.
type ScryptParams struct {
	N int // CPU/memory cost parameter (must be power of 2)
	R int // Block size parameter
	P int // Parallelization parameter
}

// DefaultScryptParams returns the default scrypt parameters.
func DefaultScryptParams() ScryptParams {
	return ScryptParams{
		N: 262144, // 2^18
		R: 8,
		P: 1,
	}
}

// EncryptedKey represents an encrypted keystore file.
type EncryptedKey struct {
	Address string          `json:"address"`
	Crypto  json.RawMessage `json:"crypto"`
	ID      string          `json:"id"`
	Version int             `json:"version"`
}

// EncryptKey encrypts a private key to Ethereum Keystore V3 format.
func EncryptKey(privateKey *ecdsa.PrivateKey, password string, params ScryptParams) ([]byte, error) {
	// Create a temporary keystore to use go-ethereum's encryption
	tmpDir, err := os.MkdirTemp("", "walletgen-ks-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Use go-ethereum's keystore with specified scrypt params
	ks := keystore.NewKeyStore(tmpDir, params.N, params.P)

	// Import the private key
	account, err := ks.ImportECDSA(privateKey, password)
	if err != nil {
		return nil, fmt.Errorf("failed to import key: %w", err)
	}

	// Read the generated keystore file
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore dir: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no keystore file generated")
	}

	keystoreFile := filepath.Join(tmpDir, files[0].Name())
	data, err := os.ReadFile(keystoreFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore file: %w", err)
	}

	// Verify the keystore is valid by parsing it
	var ek EncryptedKey
	if err := json.Unmarshal(data, &ek); err != nil {
		return nil, fmt.Errorf("failed to parse keystore: %w", err)
	}

	// Verify the address matches (keystore stores lowercase without 0x)
	expectedAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	expectedAddrLower := strings.ToLower(expectedAddr.Hex()[2:])
	if strings.ToLower(ek.Address) != expectedAddrLower {
		return nil, fmt.Errorf("address mismatch in keystore: got %s, want %s", ek.Address, expectedAddrLower)
	}

	_ = account // silence unused warning
	return data, nil
}

// WriteKeystore writes encrypted keystore JSON to a file.
func WriteKeystore(keystoreDir string, index int, data []byte, address string) (string, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(keystoreDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create keystore dir: %w", err)
	}

	// Generate filename: UTC--<timestamp>--<address>
	// For simplicity, use wallet index
	filename := fmt.Sprintf("wallet-%d--%s.json", index, address[2:]) // remove 0x
	filepath := filepath.Join(keystoreDir, filename)

	// Write with restrictive permissions
	if err := os.WriteFile(filepath, data, 0600); err != nil {
		return "", fmt.Errorf("failed to write keystore: %w", err)
	}

	return filename, nil
}

// DecryptKey decrypts a keystore file and returns the private key.
// Used for testing.
func DecryptKey(data []byte, password string) (*ecdsa.PrivateKey, error) {
	key, err := keystore.DecryptKey(data, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key: %w", err)
	}
	return key.PrivateKey, nil
}

// GenerateUUID generates a new UUID for keystore files.
func GenerateUUID() string {
	return uuid.New().String()
}
