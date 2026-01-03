package wallet

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/monolythium/wallet-gen/pkg/bech32"
)

func TestGenerateWallet(t *testing.T) {
	gen := NewGenerator("mono", nil)

	wallet, err := gen.Generate()
	if err != nil {
		t.Fatalf("failed to generate wallet: %v", err)
	}

	// Verify bech32 prefix
	if !strings.HasPrefix(wallet.Bech32, "mono1") {
		t.Errorf("expected bech32 prefix 'mono1', got '%s'", wallet.Bech32[:5])
	}

	// Verify EVM address format
	if !strings.HasPrefix(wallet.EVMAddress, "0x") {
		t.Errorf("expected EVM address prefix '0x', got '%s'", wallet.EVMAddress[:2])
	}
	if len(wallet.EVMAddress) != 42 {
		t.Errorf("expected EVM address length 42, got %d", len(wallet.EVMAddress))
	}

	// Verify private key hex length (32 bytes = 64 hex chars)
	if len(wallet.PrivKeyHex) != 64 {
		t.Errorf("expected privkey hex length 64, got %d", len(wallet.PrivKeyHex))
	}

	// Verify private key is lowercase
	if wallet.PrivKeyHex != strings.ToLower(wallet.PrivKeyHex) {
		t.Error("private key hex should be lowercase")
	}
}

func TestDeterministicGeneration(t *testing.T) {
	// Generate with same seed twice
	gen1 := NewGenerator("mono", DeterministicRand(12345))
	gen2 := NewGenerator("mono", DeterministicRand(12345))

	wallet1, err := gen1.Generate()
	if err != nil {
		t.Fatalf("failed to generate wallet1: %v", err)
	}

	wallet2, err := gen2.Generate()
	if err != nil {
		t.Fatalf("failed to generate wallet2: %v", err)
	}

	// Should produce same wallet
	if wallet1.PrivKeyHex != wallet2.PrivKeyHex {
		t.Error("deterministic generation should produce same private key")
	}
	if wallet1.Bech32 != wallet2.Bech32 {
		t.Error("deterministic generation should produce same bech32 address")
	}
	if wallet1.EVMAddress != wallet2.EVMAddress {
		t.Error("deterministic generation should produce same EVM address")
	}
}

func TestFromPrivateKey(t *testing.T) {
	// Known test vector
	privKeyHex := "4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"
	privKeyBytes, _ := hex.DecodeString(privKeyHex)
	privateKey, _ := crypto.ToECDSA(privKeyBytes)

	gen := NewGenerator("mono", nil)
	wallet, err := gen.FromPrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to create wallet from private key: %v", err)
	}

	// Verify private key matches
	if wallet.PrivKeyHex != privKeyHex {
		t.Errorf("privkey mismatch: got %s, want %s", wallet.PrivKeyHex, privKeyHex)
	}

	// Verify we can derive the address
	expectedAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	if wallet.EVMAddress != expectedAddr.Hex() {
		t.Errorf("EVM address mismatch: got %s, want %s", wallet.EVMAddress, expectedAddr.Hex())
	}
}

func TestFormatLine(t *testing.T) {
	gen := NewGenerator("mono", DeterministicRand(99999))
	wallet, _ := gen.Generate()

	line := wallet.FormatLine()

	// Should have 3 parts separated by colons
	parts := strings.Split(line, ":")
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parts))
	}

	// Part 0: bech32
	if !strings.HasPrefix(parts[0], "mono1") {
		t.Errorf("part 0 should be bech32 address, got %s", parts[0])
	}

	// Part 1: EVM address
	if !strings.HasPrefix(parts[1], "0x") {
		t.Errorf("part 1 should be EVM address, got %s", parts[1])
	}

	// Part 2: private key (64 hex chars)
	if len(parts[2]) != 64 {
		t.Errorf("part 2 should be 64-char privkey, got %d chars", len(parts[2]))
	}
}

func TestFormatLineEncrypted(t *testing.T) {
	gen := NewGenerator("mono", DeterministicRand(99999))
	wallet, _ := gen.Generate()

	line := wallet.FormatLineEncrypted("keystores/wallet-0.json")

	if !strings.Contains(line, "<keystore:") {
		t.Errorf("expected <keystore:...> marker, got %s", line)
	}

	if !strings.HasPrefix(line, "mono1") {
		t.Error("line should start with bech32 address")
	}
}

func TestAddressConsistency(t *testing.T) {
	// Verify bech32 and EVM addresses encode same 20 bytes
	gen := NewGenerator("mono", nil)

	for i := 0; i < 10; i++ {
		wallet, err := gen.Generate()
		if err != nil {
			t.Fatalf("failed to generate wallet: %v", err)
		}

		// Get address bytes from EVM address
		evmAddrBytes, err := hex.DecodeString(wallet.EVMAddress[2:])
		if err != nil {
			t.Fatalf("failed to decode EVM address: %v", err)
		}

		// Decode bech32 to get address bytes
		_, bech32Bytes, err := bech32.Decode(wallet.Bech32)
		if err != nil {
			t.Fatalf("failed to decode bech32: %v", err)
		}

		// They should match
		if hex.EncodeToString(evmAddrBytes) != hex.EncodeToString(bech32Bytes) {
			t.Errorf("address mismatch: EVM=%x, Bech32=%x", evmAddrBytes, bech32Bytes)
		}
	}
}
