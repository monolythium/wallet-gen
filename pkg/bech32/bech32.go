// Package bech32 provides bech32 encoding for Monolythium addresses.
package bech32

import (
	"fmt"

	"github.com/btcsuite/btcutil/bech32"
)

// Encode encodes a 20-byte address to bech32 with the given prefix.
func Encode(prefix string, addr []byte) (string, error) {
	if len(addr) != 20 {
		return "", fmt.Errorf("address must be 20 bytes, got %d", len(addr))
	}

	// Convert 8-bit bytes to 5-bit groups for bech32
	conv, err := bech32.ConvertBits(addr, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("failed to convert bits: %w", err)
	}

	encoded, err := bech32.Encode(prefix, conv)
	if err != nil {
		return "", fmt.Errorf("failed to encode bech32: %w", err)
	}

	return encoded, nil
}

// Decode decodes a bech32 address and returns the prefix and 20-byte address.
func Decode(encoded string) (string, []byte, error) {
	prefix, data, err := bech32.Decode(encoded)
	if err != nil {
		return "", nil, fmt.Errorf("failed to decode bech32: %w", err)
	}

	// Convert 5-bit groups back to 8-bit bytes
	conv, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return "", nil, fmt.Errorf("failed to convert bits: %w", err)
	}

	if len(conv) != 20 {
		return "", nil, fmt.Errorf("decoded address must be 20 bytes, got %d", len(conv))
	}

	return prefix, conv, nil
}

// EncodeWithPrefix is a convenience function for encoding with "mono" prefix.
func EncodeWithPrefix(addr []byte) (string, error) {
	return Encode("mono", addr)
}
