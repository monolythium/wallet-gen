# wallet-gen

Offline mass wallet generator for Monolythium.

## Features

- Generate thousands of wallets offline
- Output format: `bech32:0x:privatekey`
- Optional encrypted keystore mode (Ethereum Keystore V3)
- Progress indicator
- Deterministic mode for testing

## Installation

```bash
go install github.com/monolythium/wallet-gen/cmd/walletgen@latest
```

Or build from source:

```bash
git clone https://github.com/monolythium/wallet-gen.git
cd wallet-gen
go build -o walletgen ./cmd/walletgen
```

## Usage

### Raw Keys (Default)

```bash
walletgen --count 1000 --out wallets.txt
```

### Encrypted Keystores

```bash
# Create password file
echo "your-secure-password" > password.txt

# Generate with encryption
walletgen --count 1000 --out wallets.txt --encrypt --keystore-dir keystores --password-file password.txt
```

## Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--count` | `-n` | Number of wallets to generate | (required) |
| `--out` | `-o` | Output file path | (required) |
| `--prefix` | | Bech32 address prefix | `mono` |
| `--quiet` | | Disable progress output | `false` |
| `--seed` | | Deterministic seed (testing only) | `0` (random) |
| `--encrypt` | | Enable encrypted keystore mode | `false` |
| `--keystore-dir` | | Directory for keystore files | (required if --encrypt) |
| `--password-file` | | File containing encryption password | (required if --encrypt) |
| `--scrypt-n` | | Scrypt N parameter | `262144` |
| `--scrypt-r` | | Scrypt R parameter | `8` |
| `--scrypt-p` | | Scrypt P parameter | `1` |

## Output Format

### Raw Mode

```
mono1abc...:0xABC...:<64-hex-private-key>
```

### Encrypt Mode

```
mono1abc...:0xABC...:<keystore:keystores/wallet-0--abc....json>
```

Keystore files are written to the specified directory in Ethereum Keystore V3 format, compatible with MetaMask and geth.

## Security

- Runs completely offline with no network calls
- Private keys are never printed to stdout
- Output files are created with restrictive permissions (0600)
- Uses audited crypto libraries (go-ethereum)

## License

Business Source License 1.1 - See LICENSE file.
