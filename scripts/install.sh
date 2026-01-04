#!/usr/bin/env bash
set -euo pipefail

# wallet-gen auto-installer
# Auto-detects architecture and downloads the correct binary

REPO="monolythium/wallet-gen"
INSTALL_DIR="${INSTALL_DIR:-$HOME/bin}"
BINARY_NAME="walletgen"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*)
            echo "darwin"
            ;;
        Linux*)
            echo "linux"
            ;;
        MINGW* | MSYS* | CYGWIN*)
            echo "windows"
            ;;
        *)
            error "Unsupported operating system: $(uname -s)"
            ;;
    esac
}

# Detect architecture
detect_arch() {
    local arch
    arch="$(uname -m)"

    case "$arch" in
        x86_64 | amd64)
            echo "amd64"
            ;;
        aarch64 | arm64)
            echo "arm64"
            ;;
        armv7l)
            echo "armv7"
            ;;
        i386 | i686)
            error "32-bit architecture not supported. Please upgrade to a 64-bit OS."
            ;;
        *)
            error "Unsupported architecture: $arch"
            ;;
    esac
}

# Get latest release version
get_latest_version() {
    local version
    version=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$version" ]; then
        error "Failed to fetch latest version from GitHub"
    fi

    echo "$version"
}

# Main installation
main() {
    info "wallet-gen installer"
    echo ""

    # Detect system
    local os arch version
    os=$(detect_os)
    arch=$(detect_arch)

    info "Detected OS: $os"
    info "Detected Architecture: $arch ($(uname -m))"
    echo ""

    # Get latest version
    info "Fetching latest release version..."
    version=$(get_latest_version)
    info "Latest version: $version"
    echo ""

    # Construct download URL
    local filename extension
    if [ "$os" = "windows" ]; then
        extension="zip"
        filename="wallet-gen_${version#v}_${os}_${arch}.${extension}"
    else
        extension="tar.gz"
        filename="wallet-gen_${version#v}_${os}_${arch}.${extension}"
    fi

    local download_url="https://github.com/$REPO/releases/download/$version/$filename"

    info "Download URL: $download_url"
    echo ""

    # Create temp directory
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap 'rm -rf "$tmp_dir"' EXIT

    # Download
    info "Downloading $filename..."
    if ! curl -fsSL "$download_url" -o "$tmp_dir/$filename"; then
        error "Failed to download $filename. Please check the release exists for your platform."
    fi

    # Extract
    info "Extracting..."
    if [ "$os" = "windows" ]; then
        unzip -q "$tmp_dir/$filename" -d "$tmp_dir"
    else
        tar -xzf "$tmp_dir/$filename" -C "$tmp_dir"
    fi

    # Create install directory
    mkdir -p "$INSTALL_DIR"

    # Install binary
    local binary_path="$INSTALL_DIR/$BINARY_NAME"
    if [ "$os" = "windows" ]; then
        binary_path="${binary_path}.exe"
    fi

    info "Installing to $binary_path..."
    cp "$tmp_dir/$BINARY_NAME" "$binary_path" 2>/dev/null || cp "$tmp_dir/${BINARY_NAME}.exe" "$binary_path" 2>/dev/null || error "Binary not found in archive"
    chmod +x "$binary_path"

    echo ""
    info "Installation successful!"
    echo ""

    # Check if install dir is in PATH
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        warn "$INSTALL_DIR is not in your PATH"
        echo ""
        echo "Add it to your PATH by running:"
        echo ""
        echo "  echo 'export PATH=\"\$HOME/bin:\$PATH\"' >> ~/.bashrc  # or ~/.zshrc"
        echo "  source ~/.bashrc  # or ~/.zshrc"
        echo ""
    else
        info "You can now run: $BINARY_NAME"
    fi

    # Show version
    if command -v "$BINARY_NAME" &> /dev/null; then
        echo ""
        "$BINARY_NAME" --version 2>/dev/null || true
    fi
}

main "$@"
