// Package wallet provides wallet generation functionality.
package wallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/monolythium/wallet-gen/pkg/bech32"
)

// Wallet represents a generated wallet with all address formats.
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	Bech32     string
	EVMAddress string
	PrivKeyHex string
}

// Generator generates wallets.
type Generator struct {
	prefix string
	rand   io.Reader
}

// NewGenerator creates a new wallet generator.
// If randSource is nil, crypto.GenerateKey uses crypto/rand internally.
func NewGenerator(prefix string, randSource io.Reader) *Generator {
	return &Generator{
		prefix: prefix,
		rand:   randSource,
	}
}

// Generate generates a single wallet.
func (g *Generator) Generate() (*Wallet, error) {
	var privateKey *ecdsa.PrivateKey
	var err error

	if g.rand != nil {
		// Use provided random source (for deterministic testing)
		privateKey, err = ecdsa.GenerateKey(crypto.S256(), g.rand)
	} else {
		// Use crypto/rand via go-ethereum
		privateKey, err = crypto.GenerateKey()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	return g.FromPrivateKey(privateKey)
}

// FromPrivateKey creates a wallet from an existing private key.
func (g *Generator) FromPrivateKey(privateKey *ecdsa.PrivateKey) (*Wallet, error) {
	// Get EVM address (20 bytes)
	evmAddr := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Encode to bech32
	bech32Addr, err := bech32.Encode(g.prefix, evmAddr.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to encode bech32: %w", err)
	}

	// Get private key hex (lowercase, no 0x prefix)
	privKeyBytes := crypto.FromECDSA(privateKey)
	privKeyHex := hex.EncodeToString(privKeyBytes)

	return &Wallet{
		PrivateKey: privateKey,
		Bech32:     bech32Addr,
		EVMAddress: evmAddr.Hex(),
		PrivKeyHex: privKeyHex,
	}, nil
}

// FormatLine formats the wallet as a single line for output.
func (w *Wallet) FormatLine() string {
	return fmt.Sprintf("%s:%s:%s", w.Bech32, w.EVMAddress, w.PrivKeyHex)
}

// FormatLineEncrypted formats the wallet line for encrypted mode.
func (w *Wallet) FormatLineEncrypted(keystorePath string) string {
	return fmt.Sprintf("%s:%s:<keystore:%s>", w.Bech32, w.EVMAddress, keystorePath)
}

// deterministicReader provides deterministic random bytes using SHA256.
// This ensures consistent output across all Go versions for testing.
type deterministicReader struct {
	seed    [32]byte
	counter uint64
	buf     []byte
}

func (d *deterministicReader) Read(p []byte) (n int, err error) {
	for len(p) > 0 {
		if len(d.buf) == 0 {
			// Generate next block: SHA256(seed || counter)
			var input [40]byte
			copy(input[:32], d.seed[:])
			binary.BigEndian.PutUint64(input[32:], d.counter)
			hash := sha256.Sum256(input[:])
			d.buf = hash[:]
			d.counter++
		}
		copied := copy(p, d.buf)
		p = p[copied:]
		d.buf = d.buf[copied:]
		n += copied
	}
	return n, nil
}

// DeterministicRand creates a deterministic random source for testing.
// Uses SHA256-based expansion which is consistent across all Go versions.
func DeterministicRand(seed int64) io.Reader {
	var seedBytes [32]byte
	binary.BigEndian.PutUint64(seedBytes[:8], uint64(seed))
	return &deterministicReader{seed: seedBytes}
}
