# wallet-gen

Offline mass wallet generator for Monolythium.

## Features

- **Double-click to run** - Built-in interactive wizard, no terminal experience needed
- Generate thousands of wallets offline
- Output format: `bech32:0x:privatekey`
- Optional encrypted keystore mode (Ethereum Keystore V3)
- Progress indicator
- Deterministic mode for testing
- Cross-platform: Windows, macOS, Linux

## Download & Run (Recommended)

### Step 1: Download Pre-built Binary

Download the latest release for your operating system from [GitHub Releases](https://github.com/monolythium/wallet-gen/releases):

| Platform | File | Architecture |
|----------|------|--------------|
| macOS (Intel) | `wallet-gen_X.Y.Z_darwin_amd64.tar.gz` | x86_64 |
| macOS (Apple Silicon) | `wallet-gen_X.Y.Z_darwin_arm64.tar.gz` | ARM64 (M1/M2/M3) |
| Linux | `wallet-gen_X.Y.Z_linux_amd64.tar.gz` | x86_64 |
| Linux (ARM) | `wallet-gen_X.Y.Z_linux_arm64.tar.gz` | ARM64 |
| Windows | `wallet-gen_X.Y.Z_windows_amd64.zip` | x86_64 |

### Step 2: Verify Checksums

Always verify the checksum of downloaded files before running:

```bash
# Download checksums.txt from the same release page
# Then verify (Linux):
sha256sum -c checksums.txt

# Or on macOS:
shasum -a 256 -c checksums.txt
```

On Windows (PowerShell):
```powershell
# Check a single file
(Get-FileHash walletgen.exe -Algorithm SHA256).Hash
# Compare with the hash in checksums.txt
```

### Step 3: Extract and Run

#### macOS Quick Start

1. Extract the archive (double-click the `.tar.gz` file)
2. **Double-click `wallet-gen.command`** in Finder
3. Terminal will open and guide you through wallet generation

**If macOS shows "cannot be opened" or security warning:**
- Right-click → Open → Open (first time only)
- Or: System Settings → Privacy & Security → scroll down → Allow

#### Linux Quick Start

```bash
# Extract
tar -xzf wallet-gen_*_linux_amd64.tar.gz

# Make executable
chmod +x walletgen

# Run interactive wizard
./walletgen

# Or use the guided script
./scripts/wallet-gen.sh
```

#### Windows Quick Start

1. Extract the ZIP file
2. **Double-click `walletgen.exe`**
3. Follow the interactive wizard prompts

The wizard will guide you through wallet count, output location, and encryption options.

**Default output location:**
```
C:\Users\YourName\wallet-gen-output\20240115-143022\
  wallets.txt
  keystores\  (if encrypted mode selected)
```

---

## Troubleshooting

### Fixing "Permission Denied" (macOS/Linux)

If you get "permission denied" when running the binary:

```bash
# Make the binary executable
chmod +x walletgen

# Also make scripts executable
chmod +x wallet-gen.command
chmod +x scripts/wallet-gen.sh
```

### macOS Gatekeeper / Quarantine Issues

If macOS blocks the binary with "cannot be opened because the developer cannot be verified":

**Option 1: Right-click method**
- Right-click the file → Open → Open

**Option 2: Remove quarantine attribute**
```bash
xattr -dr com.apple.quarantine walletgen
xattr -dr com.apple.quarantine wallet-gen.command
```

> **Warning:** Only remove quarantine for files you downloaded from trusted sources (official GitHub releases). Never run `xattr -dr` on files from unknown origins.

---

## Adding to PATH (Optional)

Adding `walletgen` to your PATH lets you run it from any directory.

### macOS/Linux

**Option 1: Create ~/bin and add to PATH (recommended)**
```bash
mkdir -p ~/bin
cp walletgen ~/bin/
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.zshrc  # or ~/.bashrc
source ~/.zshrc
```

**Option 2: Copy to system PATH**
```bash
sudo cp walletgen /usr/local/bin/
```

### Windows

1. Create a folder for CLI tools, e.g., `C:\mono\bin`
2. Copy `walletgen.exe` to that folder
3. Add the folder to your PATH:
   - Search "Environment Variables" in Start menu
   - Edit "Path" under User variables
   - Add `C:\mono\bin`
   - Click OK and restart your terminal

Now you can run `walletgen` from any directory.

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

### Interactive Wizard Mode (All Platforms - Recommended)

Simply run the binary with no arguments:

```bash
# macOS/Linux
./walletgen

# Windows (double-click walletgen.exe or run in terminal)
.\walletgen.exe
```

The interactive wizard will guide you through:
1. How many wallets to generate (1-100,000)
2. Output directory (defaults to `~/wallet-gen-output/<timestamp>/` or `%USERPROFILE%\wallet-gen-output\<timestamp>\`)
3. Output mode:
   - **Raw text** - Private keys in plain text (requires security confirmation)
   - **Encrypted keystores** - Password-protected (recommended)
   - **Both** - Raw backup + encrypted keystores
4. Password setup (for encrypted mode)
5. Confirmation summary before generation

### Legacy Script Mode (macOS/Linux)

The bash script is still available for those who prefer it:

```bash
./scripts/wallet-gen.sh
```

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

### Default Output Folders (Wizard Mode)

When you run the interactive wizard (double-click or run without flags), it uses safe default locations:

| Platform | Default Output Path |
|----------|---------------------|
| macOS/Linux | `~/wallet-gen-output/<timestamp>/` |
| Windows | `%USERPROFILE%\wallet-gen-output\<timestamp>\` |

Example structure:
```
wallet-gen-output/
  20240115-143022/           # Timestamped directory
    wallets.txt              # Wallet addresses and keys
    keystores/               # Encrypted keystore files (if enabled)
      wallet-0--abc123....json
      wallet-1--def456....json
```

### CLI Mode: Working Directory Matters

When using CLI flags, output is relative to your current directory:

```bash
# If you run from /Users/alice/Downloads:
./walletgen --count 10 --out wallets.txt
# Creates: /Users/alice/Downloads/wallets.txt

# Use absolute paths to be explicit:
./walletgen --count 10 --out ~/secure-wallets/wallets.txt
# Creates: /Users/alice/secure-wallets/wallets.txt
```

**Recommendation:** Use the interactive wizard or explicit absolute paths to avoid accidentally writing sensitive files to unexpected locations.

### File Permissions

- **Output files**: Created with `0600` (owner read/write only)
- **Keystore directories**: Created with `0700` (owner access only)

Verify permissions after generation:
```bash
ls -la ~/wallet-gen-output/
```

## Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--count` | `-n` | Number of wallets to generate | (wizard mode if omitted) |
| `--out` | `-o` | Output file path | (wizard mode if omitted) |
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

### Raw Mode is Unsafe for Production

- **Raw text mode stores private keys unencrypted.** Anyone with file access can steal all funds.
- **Always prefer encrypted keystore mode** for any real-world use.
- If you must use raw mode, delete the file securely after importing keys:
  ```bash
  shred -u wallets.txt  # Linux
  rm -P wallets.txt     # macOS
  ```

### Password File Handling

- Password files are stored in **plain text**.
- The wizard can prompt for passwords interactively (no file needed).
- If you use `--password-file`:
  - Set restrictive permissions first: `chmod 600 password.txt`
  - **Delete it immediately after wallet generation**
  - Never commit password files to version control

### Offline Operation

This tool runs **completely offline**:
- No network calls are made
- No telemetry or external connections
- Safe to run on air-gapped machines

### Summary

| Feature | Status |
|---------|--------|
| Offline operation | Yes, no network calls |
| Private keys to stdout | Never |
| File permissions | `0600` (owner only) |
| Crypto library | go-ethereum (audited) |
| Keystore format | Ethereum Keystore V3 (scrypt) |

## Testing

```bash
go test ./...
```

## License

Business Source License 1.1 - See LICENSE file.
