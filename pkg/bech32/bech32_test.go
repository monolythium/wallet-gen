package bech32

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestEncodeDecodeRoundtrip(t *testing.T) {
	// Test with known address bytes
	addrHex := "8eaf92eb6a7c48b95998f7c1df79402f3cc8bfa2"
	addrBytes, err := hex.DecodeString(addrHex)
	if err != nil {
		t.Fatalf("failed to decode hex: %v", err)
	}

	// Encode
	encoded, err := Encode("mono", addrBytes)
	if err != nil {
		t.Fatalf("failed to encode: %v", err)
	}

	// Verify prefix
	if encoded[:5] != "mono1" {
		t.Errorf("expected prefix 'mono1', got '%s'", encoded[:5])
	}

	// Decode
	prefix, decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if prefix != "mono" {
		t.Errorf("expected prefix 'mono', got '%s'", prefix)
	}

	if !bytes.Equal(decoded, addrBytes) {
		t.Errorf("roundtrip failed: got %x, want %x", decoded, addrBytes)
	}
}

func TestEncodeWithPrefix(t *testing.T) {
	addrHex := "0102030405060708091011121314151617181920"
	addrBytes, _ := hex.DecodeString(addrHex)

	encoded, err := EncodeWithPrefix(addrBytes)
	if err != nil {
		t.Fatalf("failed to encode: %v", err)
	}

	if encoded[:5] != "mono1" {
		t.Errorf("expected prefix 'mono1', got '%s'", encoded[:5])
	}
}

func TestEncodeInvalidLength(t *testing.T) {
	// 19 bytes - should fail
	addr := make([]byte, 19)
	_, err := Encode("mono", addr)
	if err == nil {
		t.Error("expected error for 19-byte address")
	}

	// 21 bytes - should fail
	addr = make([]byte, 21)
	_, err = Encode("mono", addr)
	if err == nil {
		t.Error("expected error for 21-byte address")
	}
}

func TestDecodeInvalid(t *testing.T) {
	// Invalid bech32
	_, _, err := Decode("invalid")
	if err == nil {
		t.Error("expected error for invalid bech32")
	}

	// Valid bech32 but wrong checksum
	_, _, err = Decode("mono1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq")
	if err == nil {
		t.Error("expected error for invalid checksum")
	}
}

func TestKnownVector(t *testing.T) {
	// Known test vector: dev_fund address
	// EVM: 0x8Eaf92Eb6a7c48b95998F7C1df79402F3cC8BFa2
	// Bech32: mono136he96m203ytjkvc7lqa772q9u7v30azf59lm7
	addrHex := "8eaf92eb6a7c48b95998f7c1df79402f3cc8bfa2"
	addrBytes, _ := hex.DecodeString(addrHex)

	encoded, err := Encode("mono", addrBytes)
	if err != nil {
		t.Fatalf("failed to encode: %v", err)
	}

	expected := "mono136he96m203ytjkvc7lqa772q9u7v30azf59lm7"
	if encoded != expected {
		t.Errorf("encoding mismatch: got %s, want %s", encoded, expected)
	}
}
