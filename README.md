# wallet-gen

Offline mass wallet generator for Monolythium.

## Features

- Generate thousands of wallets offline
- Guided interactive script for beginners
- Output format: `bech32:0x:privatekey`
- Optional encrypted keystore mode (Ethereum Keystore V3)
- Progress indicator
- Deterministic mode for testing

## Quick Start

The easiest way to get started is using the guided script:

```bash
# Clone the repository
git clone https://github.com/monolythium/wallet-gen.git
cd wallet-gen

# Run the guided script
./scripts/wallet-gen.sh
```

The script will:
1. Build the `walletgen` binary if needed
2. Guide you through wallet count, output location, and mode selection
3. Warn you about security implications
4. Generate wallets with safe defaults

> **Windows Users**: Use WSL (Windows Subsystem for Linux) or run the `walletgen` CLI directly.

## Installation

### Using Go Install

```bash
go install github.com/monolythium/wallet-gen/cmd/walletgen@latest
```

### Build from Source

```bash
git clone https://github.com/monolythium/wallet-gen.git
cd wallet-gen
go build -o walletgen ./cmd/walletgen
```

## Usage

### Guided Mode (Recommended for Beginners)

```bash
./scripts/wallet-gen.sh
```

Follow the interactive prompts to:
- Set wallet count
- Choose output directory
- Select output mode (plain text or encrypted)
- Configure encryption password

### CLI Mode (Advanced)

#### Raw Keys (Default)

```bash
walletgen --count 1000 --out wallets.txt
```

#### Encrypted Keystores (Recommended)

```bash
# Create password file
echo "your-secure-password" > password.txt
chmod 600 password.txt

# Generate with encryption
walletgen --count 1000 --out wallets.txt \
  --encrypt --keystore-dir keystores --password-file password.txt

# IMPORTANT: Delete password file after use
rm password.txt
```

## Output Location

**Always use explicit output paths.** This ensures you know exactly where sensitive files are written.

### Recommended Output Structure

```
output/
  20240115-143022/           # Timestamped directory
    wallets.txt              # Wallet index file (0600 permissions)
    keystores/               # Encrypted keystore files (0700 directory)
      wallet-0--abc123....json
      wallet-1--def456....json
      ...
    run-summary.txt          # Generation summary
```

### Default Behavior

- **Guided script**: Creates timestamped directories under `./output/`
- **CLI tool**: Writes to the exact path specified with `--out`
- **File permissions**: All sensitive files are created with `0600` (owner read/write only)
- **Directory permissions**: Keystore directories are created with `0700` (owner access only)

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

## Security Warnings

### Private Key Safety

- **Plain text mode stores private keys unencrypted.** Anyone with file access can steal all funds.
- **Always prefer encrypted keystore mode** for any real-world use.
- **Delete plain text files securely** after importing keys to a secure wallet.
- Consider using `shred` or `srm` for secure deletion:
  ```bash
  shred -u wallets.txt  # Linux
  rm -P wallets.txt     # macOS
  ```

### Password File Security

- Password files are stored in **plain text**.
- The guided script can create a temporary password file that is deleted after generation.
- If you create your own password file:
  - Set restrictive permissions: `chmod 600 password.txt`
  - Delete it immediately after use
  - Never commit password files to version control

### Offline Operation

This tool is designed to run **completely offline**:
- No network calls are made
- No telemetry or external connections
- Safe to run on air-gapped machines

### File Permissions

The tool automatically sets restrictive permissions:
- Output files: `0600` (read/write by owner only)
- Keystore directories: `0700` (access by owner only)

Verify permissions after generation:
```bash
ls -la output/
```

## Security

- Runs completely offline with no network calls
- Private keys are never printed to stdout
- Output files are created with restrictive permissions (0600)
- Uses audited crypto libraries (go-ethereum)
- Encrypted keystores use Ethereum Keystore V3 format with scrypt KDF

## Testing

```bash
go test ./...
```

## License

Business Source License 1.1 - See LICENSE file.
