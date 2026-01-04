# wallet-gen

Offline mass wallet generator for Monolythium.

## Features

- Generate thousands of wallets offline
- Guided interactive script for beginners
- Output format: `bech32:0x:privatekey`
- Optional encrypted keystore mode (Ethereum Keystore V3)
- Progress indicator
- Deterministic mode for testing

## Downloads & Running

### Download Pre-built Binaries

Download the latest release for your operating system from [GitHub Releases](https://github.com/monolythium/wallet-gen/releases):

| Platform | File | Architecture |
|----------|------|--------------|
| macOS (Intel) | `wallet-gen_X.Y.Z_darwin_amd64.tar.gz` | x86_64 |
| macOS (Apple Silicon) | `wallet-gen_X.Y.Z_darwin_arm64.tar.gz` | ARM64 (M1/M2/M3) |
| Linux | `wallet-gen_X.Y.Z_linux_amd64.tar.gz` | x86_64 |
| Linux (ARM) | `wallet-gen_X.Y.Z_linux_arm64.tar.gz` | ARM64 |
| Windows | `wallet-gen_X.Y.Z_windows_amd64.zip` | x86_64 |

### Verify Checksums

Always verify the checksum of downloaded files:

```bash
# Download checksums.txt from the release
sha256sum -c checksums.txt
# or on macOS:
shasum -a 256 -c checksums.txt
```

### macOS Quick Start (Double-Click)

1. Download and extract the archive
2. **Double-click `wallet-gen.command`** in Finder
3. Terminal will open and guide you through wallet generation

If you get a security warning:
- Right-click → Open → Open (first time only)
- Or: System Settings → Privacy & Security → Allow

### Linux Quick Start

```bash
# Extract
tar -xzf wallet-gen_*_linux_amd64.tar.gz

# Make executable
chmod +x walletgen

# Run guided mode
./scripts/wallet-gen.sh

# Or run CLI directly
./walletgen --count 10 --out wallets.txt
```

### Windows Quick Start

> **Note**: The guided script (`wallet-gen.sh`) is designed for macOS/Linux. On Windows, use `walletgen.exe` directly as shown below.

**Option A: PowerShell/CMD (Recommended)**
```powershell
# 1. Download wallet-gen_X.Y.Z_windows_amd64.exe from GitHub Releases
#    (or extract from the .zip for the full package)

# 2. Create output directory
mkdir C:\mono

# 3. Generate wallets (plain text mode)
.\walletgen.exe --count 10 --out C:\mono\wallets.txt

# 4. Or with encryption (recommended)
echo your-secure-password > C:\mono\password.txt
.\walletgen.exe --count 10 --out C:\mono\wallets.txt --encrypt --keystore-dir C:\mono\keystores --password-file C:\mono\password.txt
del C:\mono\password.txt
```

**Option B: WSL (if you prefer the guided script)**
```powershell
# Install WSL if not already installed
wsl --install

# Then run in WSL terminal:
./scripts/wallet-gen.sh
```

> **Important**: Always use absolute paths on Windows (e.g., `C:\mono\wallets.txt`) to know exactly where files are saved.

### Fixing "Permission Denied" Errors

If you get "permission denied" when running downloaded binaries:

```bash
# Make the binary executable
chmod +x walletgen

# For the launcher script
chmod +x wallet-gen.command
chmod +x scripts/wallet-gen.sh
```

### Adding to PATH

To run `walletgen` from anywhere:

**macOS/Linux (zsh):**
```bash
# Option 1: Copy to a directory in PATH
sudo cp walletgen /usr/local/bin/

# Option 2: Add current directory to PATH
echo 'export PATH="$PATH:'"$(pwd)"'"' >> ~/.zshrc
source ~/.zshrc
```

**macOS/Linux (bash):**
```bash
echo 'export PATH="$PATH:'"$(pwd)"'"' >> ~/.bashrc
source ~/.bashrc
```

## Quick Start (From Source)

If you prefer to build from source:

```bash
# Clone the repository
git clone https://github.com/monolythium/wallet-gen.git
cd wallet-gen

# Run the guided script (auto-builds if needed)
./scripts/wallet-gen.sh
```

The script will:
1. Build the `walletgen` binary if needed
2. Guide you through wallet count, output location, and mode selection
3. Warn you about security implications
4. Generate wallets with safe defaults

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

### Guided Mode (macOS/Linux - Recommended for Beginners)

```bash
./scripts/wallet-gen.sh
```

> **Windows users**: Use `walletgen.exe` directly (see [Windows Quick Start](#windows-quick-start) above).

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

### Important: Working Directory Matters

The output location depends on where you run the command from:

```bash
# If you run from /Users/alice/Downloads:
./walletgen --out wallets.txt
# Creates: /Users/alice/Downloads/wallets.txt

# Use absolute paths to be explicit:
./walletgen --out /Users/alice/secure-wallets/wallets.txt
# Creates: /Users/alice/secure-wallets/wallets.txt
```

**Recommendation**: Always create a dedicated folder for wallet output:

```bash
# macOS/Linux
mkdir -p ~/wallet-gen-output
cd ~/wallet-gen-output
walletgen --count 100 --out wallets.txt

# Windows
mkdir C:\mono
cd C:\mono
walletgen.exe --count 100 --out wallets.txt
```

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
