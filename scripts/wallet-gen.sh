#!/usr/bin/env bash
#
# wallet-gen.sh - Guided wrapper for Monolythium wallet generator
#
# This script provides a noob-friendly interactive interface for generating
# wallets offline with safe defaults. It calls the walletgen CLI internally.
#
# SECURITY NOTES:
# - Runs completely offline (no network calls)
# - Private keys are NEVER printed to stdout
# - Output files are created with restrictive permissions (0600)
#
# Usage:
#   ./scripts/wallet-gen.sh           # Interactive mode
#   ./scripts/wallet-gen.sh --help    # Show help
#
# For advanced usage, use the walletgen CLI directly.

set -euo pipefail

# =============================================================================
# Configuration
# =============================================================================

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
WALLETGEN_BIN=""

# Colors (disabled if not a terminal)
if [[ -t 1 ]]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    BLUE='\033[0;34m'
    BOLD='\033[1m'
    NC='\033[0m' # No Color
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    BOLD=''
    NC=''
fi

# =============================================================================
# Helper Functions
# =============================================================================

print_header() {
    echo ""
    echo -e "${BLUE}${BOLD}========================================${NC}"
    echo -e "${BLUE}${BOLD}  Monolythium Wallet Generator${NC}"
    echo -e "${BLUE}${BOLD}  Guided Mode (Offline Only)${NC}"
    echo -e "${BLUE}${BOLD}========================================${NC}"
    echo ""
}

print_warning() {
    echo -e "${YELLOW}WARNING: $1${NC}"
}

print_error() {
    echo -e "${RED}ERROR: $1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}$1${NC}"
}

print_info() {
    echo -e "${BLUE}$1${NC}"
}

# Prompt for input with default value
prompt_with_default() {
    local prompt="$1"
    local default="$2"
    local result

    if [[ -n "$default" ]]; then
        read -r -p "$prompt [$default]: " result
        echo "${result:-$default}"
    else
        read -r -p "$prompt: " result
        echo "$result"
    fi
}

# Prompt for password (hidden input)
prompt_password() {
    local prompt="$1"
    local password

    read -r -s -p "$prompt: " password
    echo ""
    echo "$password"
}

# Confirm yes/no
confirm() {
    local prompt="$1"
    local response

    read -r -p "$prompt [y/N]: " response
    case "$response" in
        [yY][eE][sS]|[yY]) return 0 ;;
        *) return 1 ;;
    esac
}

# Find the walletgen binary
find_walletgen() {
    # Check if already built in repo
    if [[ -x "$REPO_ROOT/walletgen" ]]; then
        WALLETGEN_BIN="$REPO_ROOT/walletgen"
        return 0
    fi

    # Check if in PATH
    if command -v walletgen &> /dev/null; then
        WALLETGEN_BIN="walletgen"
        return 0
    fi

    # Try to build it
    echo "walletgen binary not found. Attempting to build..."
    if command -v go &> /dev/null; then
        (cd "$REPO_ROOT" && go build -o walletgen ./cmd/walletgen)
        if [[ -x "$REPO_ROOT/walletgen" ]]; then
            WALLETGEN_BIN="$REPO_ROOT/walletgen"
            print_success "Built walletgen successfully"
            return 0
        fi
    fi

    print_error "Could not find or build walletgen binary"
    echo "Please run: go build -o walletgen ./cmd/walletgen"
    exit 1
}

# Validate positive integer
validate_positive_int() {
    local value="$1"
    if [[ ! "$value" =~ ^[0-9]+$ ]] || [[ "$value" -le 0 ]]; then
        return 1
    fi
    return 0
}

# Create directory with secure permissions
secure_mkdir() {
    local dir="$1"
    mkdir -p "$dir"
    chmod 700 "$dir"
}

# Set secure file permissions
secure_file() {
    local file="$1"
    if [[ -f "$file" ]]; then
        chmod 600 "$file"
    fi
}

# Generate timestamp for default directory
get_timestamp() {
    date +"%Y%m%d-%H%M%S"
}

# Show help
show_help() {
    cat << EOF
Monolythium Wallet Generator - Guided Mode

USAGE:
    $0 [OPTIONS]

OPTIONS:
    --help, -h      Show this help message
    --version       Show version info

DESCRIPTION:
    This is a guided, interactive wrapper for the walletgen CLI tool.
    It provides a noob-friendly interface with safe defaults for generating
    Monolythium wallets offline.

    For advanced usage or scripting, use the walletgen CLI directly:
        walletgen --count 100 --out wallets.txt

SECURITY:
    - Runs completely offline (no network calls)
    - Private keys are NEVER printed to stdout
    - Output files are created with restrictive permissions (0600)
    - Encrypted mode uses Ethereum Keystore V3 format

MODES:
    1) Plain text - Outputs private keys directly (NOT RECOMMENDED)
    2) Encrypted keystores - Uses secure encryption (RECOMMENDED)
    3) Both - Creates both outputs

EOF
}

# Spinner for long operations
spinner() {
    local pid=$1
    local delay=0.1
    local spinstr='|/-\'
    while ps -p "$pid" > /dev/null 2>&1; do
        local temp=${spinstr#?}
        printf " [%c]  " "$spinstr"
        local spinstr=$temp${spinstr%"$temp"}
        sleep $delay
        printf "\b\b\b\b\b\b"
    done
    printf "      \b\b\b\b\b\b"
}

# =============================================================================
# Main Interactive Flow
# =============================================================================

run_interactive() {
    print_header

    echo "This script will guide you through generating Monolythium wallets."
    echo "All operations run completely OFFLINE - no network calls are made."
    echo ""

    # Find walletgen binary
    find_walletgen
    echo ""

    # -------------------------------------------------------------------------
    # Step 1: Number of wallets
    # -------------------------------------------------------------------------
    print_info "Step 1/6: Wallet Count"
    echo "-----------------------"

    local count
    while true; do
        count=$(prompt_with_default "How many wallets to generate?" "10")
        if validate_positive_int "$count"; then
            break
        fi
        print_error "Please enter a positive integer"
    done
    echo ""

    # -------------------------------------------------------------------------
    # Step 2: Output directory
    # -------------------------------------------------------------------------
    print_info "Step 2/6: Output Directory"
    echo "---------------------------"

    local timestamp
    timestamp=$(get_timestamp)
    local default_dir="./output/${timestamp}"

    echo "Output will be saved to a directory of your choice."
    echo "The directory will be created if it doesn't exist."

    local output_dir
    output_dir=$(prompt_with_default "Output directory" "$default_dir")
    echo ""

    # -------------------------------------------------------------------------
    # Step 3: Output filename
    # -------------------------------------------------------------------------
    print_info "Step 3/6: Output Filename"
    echo "--------------------------"

    local output_file
    output_file=$(prompt_with_default "Output filename (inside output directory)" "wallets.txt")
    echo ""

    # -------------------------------------------------------------------------
    # Step 4: Output mode
    # -------------------------------------------------------------------------
    print_info "Step 4/6: Output Mode"
    echo "----------------------"

    echo ""
    echo "Choose output mode:"
    echo ""
    echo -e "  ${RED}1) Plain text (UNSAFE)${NC}"
    echo "     Format: bech32:0x:privatekey"
    echo "     Private keys are stored in clear text."
    echo ""
    echo -e "  ${GREEN}2) Encrypted keystores (RECOMMENDED)${NC}"
    echo "     Format: bech32:0x:<keystore:path>"
    echo "     Private keys are encrypted with your password."
    echo "     Compatible with MetaMask and geth."
    echo ""
    echo "  3) Both (plain text + keystores)"
    echo ""

    local mode
    while true; do
        mode=$(prompt_with_default "Enter choice (1/2/3)" "2")
        case "$mode" in
            1|2|3) break ;;
            *) print_error "Please enter 1, 2, or 3" ;;
        esac
    done
    echo ""

    # Handle plain text warning
    local use_plaintext=false
    local use_encrypt=false

    case "$mode" in
        1)
            use_plaintext=true
            ;;
        2)
            use_encrypt=true
            ;;
        3)
            use_plaintext=true
            use_encrypt=true
            ;;
    esac

    # Plain text safety confirmation
    if [[ "$use_plaintext" == true ]]; then
        echo ""
        echo -e "${RED}${BOLD}======================================${NC}"
        echo -e "${RED}${BOLD}         SECURITY WARNING${NC}"
        echo -e "${RED}${BOLD}======================================${NC}"
        echo ""
        echo -e "${YELLOW}You have selected plain text output mode.${NC}"
        echo ""
        echo "This will write private keys in CLEAR TEXT to disk."
        echo "Anyone with access to this file can steal all funds."
        echo ""
        echo "Risks:"
        echo "  - Private keys are not encrypted"
        echo "  - File can be read by anyone with file access"
        echo "  - Malware can easily scan for and exfiltrate keys"
        echo "  - No protection if device is lost or stolen"
        echo ""
        echo -e "${BOLD}To proceed, type exactly: YES-I-UNDERSTAND${NC}"
        echo ""

        local confirmation
        read -r -p "Confirmation: " confirmation

        if [[ "$confirmation" != "YES-I-UNDERSTAND" ]]; then
            print_error "Confirmation not received. Aborting."
            exit 1
        fi
        echo ""
        print_warning "Proceeding with plain text mode at your own risk."
        echo ""
    fi

    # Password handling for encrypt mode
    local password=""
    local password_file=""
    local temp_password_file=""
    local delete_password_file=false

    if [[ "$use_encrypt" == true ]]; then
        print_info "Password Setup (for encryption)"
        echo "--------------------------------"
        echo ""
        echo "Encrypted keystores require a password."
        echo ""
        echo "Choose password method:"
        echo "  A) Enter password interactively (recommended)"
        echo "  B) Use existing password file"
        echo ""

        local pw_method
        while true; do
            pw_method=$(prompt_with_default "Enter choice (A/B)" "A")
            case "$pw_method" in
                [aA]) pw_method="A"; break ;;
                [bB]) pw_method="B"; break ;;
                *) print_error "Please enter A or B" ;;
            esac
        done
        echo ""

        if [[ "$pw_method" == "A" ]]; then
            # Interactive password entry
            local pw1 pw2
            while true; do
                pw1=$(prompt_password "Enter encryption password")
                pw2=$(prompt_password "Confirm encryption password")

                if [[ "$pw1" != "$pw2" ]]; then
                    print_error "Passwords do not match. Please try again."
                    echo ""
                    continue
                fi

                if [[ -z "$pw1" ]]; then
                    print_error "Password cannot be empty."
                    echo ""
                    continue
                fi

                if [[ ${#pw1} -lt 8 ]]; then
                    print_warning "Password is less than 8 characters. This is not recommended."
                    if ! confirm "Continue anyway?"; then
                        echo ""
                        continue
                    fi
                fi

                password="$pw1"
                break
            done

            # Will create temp password file later
            temp_password_file="${output_dir}/.password.tmp"
            delete_password_file=true
            echo ""

        else
            # Use existing password file
            while true; do
                password_file=$(prompt_with_default "Path to password file" "")

                if [[ -z "$password_file" ]]; then
                    print_error "Password file path is required."
                    continue
                fi

                if [[ ! -f "$password_file" ]]; then
                    print_error "File not found: $password_file"
                    continue
                fi

                break
            done
            echo ""
        fi

        echo ""
        print_warning "Password file is stored in plaintext. Keep it secure!"
        echo "Consider deleting the password file after use."
        echo ""
    fi

    # -------------------------------------------------------------------------
    # Step 5: Prefix
    # -------------------------------------------------------------------------
    print_info "Step 5/6: Address Prefix"
    echo "-------------------------"

    echo "Bech32 address prefix (default: mono for Monolythium)"
    local prefix
    prefix=$(prompt_with_default "Prefix" "mono")
    echo ""

    # -------------------------------------------------------------------------
    # Step 6: Confirmation
    # -------------------------------------------------------------------------
    print_info "Step 6/6: Confirm Settings"
    echo "---------------------------"
    echo ""
    echo "Please review your settings:"
    echo ""
    echo "  Wallet count:     $count"
    echo "  Output directory: $output_dir"
    echo "  Output file:      $output_file"
    echo "  Address prefix:   $prefix"
    echo ""

    if [[ "$use_plaintext" == true ]] && [[ "$use_encrypt" == true ]]; then
        echo "  Mode:             Both (plain text + encrypted)"
        echo "  Plain output:     $output_dir/$output_file"
        echo "  Encrypted output: $output_dir/${output_file%.txt}-encrypted.txt"
        echo "  Keystore dir:     $output_dir/keystores/"
    elif [[ "$use_encrypt" == true ]]; then
        echo "  Mode:             Encrypted keystores (recommended)"
        echo "  Keystore dir:     $output_dir/keystores/"
    else
        echo -e "  Mode:             ${RED}Plain text (UNSAFE)${NC}"
    fi
    echo ""

    if ! confirm "Proceed with generation?"; then
        echo "Aborted."
        exit 0
    fi
    echo ""

    # -------------------------------------------------------------------------
    # Execute generation
    # -------------------------------------------------------------------------
    print_info "Generating wallets..."
    echo ""

    # Create output directory with secure permissions
    secure_mkdir "$output_dir"

    # Write password file if needed
    if [[ -n "$temp_password_file" ]]; then
        echo "$password" > "$temp_password_file"
        secure_file "$temp_password_file"
        password_file="$temp_password_file"
    fi

    local output_path="$output_dir/$output_file"
    local encrypted_output_path="$output_dir/${output_file%.txt}-encrypted.txt"
    local keystore_dir="$output_dir/keystores"
    local summary_file="$output_dir/run-summary.txt"

    local start_time
    start_time=$(date +%s)

    # Run generation
    if [[ "$use_plaintext" == true ]] && [[ "$use_encrypt" == true ]]; then
        # Both modes: run twice
        echo "Generating plain text output..."
        "$WALLETGEN_BIN" --count "$count" --out "$output_path" --prefix "$prefix"
        secure_file "$output_path"

        echo ""
        echo "Generating encrypted keystores..."
        secure_mkdir "$keystore_dir"
        "$WALLETGEN_BIN" --count "$count" --out "$encrypted_output_path" --prefix "$prefix" \
            --encrypt --keystore-dir "$keystore_dir" --password-file "$password_file"
        secure_file "$encrypted_output_path"

    elif [[ "$use_encrypt" == true ]]; then
        # Encrypted only
        secure_mkdir "$keystore_dir"
        "$WALLETGEN_BIN" --count "$count" --out "$output_path" --prefix "$prefix" \
            --encrypt --keystore-dir "$keystore_dir" --password-file "$password_file"
        secure_file "$output_path"

    else
        # Plain text only
        "$WALLETGEN_BIN" --count "$count" --out "$output_path" --prefix "$prefix"
        secure_file "$output_path"
    fi

    local end_time
    end_time=$(date +%s)
    local duration=$((end_time - start_time))

    # Clean up temp password file
    if [[ "$delete_password_file" == true ]] && [[ -f "$temp_password_file" ]]; then
        if confirm "Delete temporary password file?"; then
            rm -f "$temp_password_file"
            print_success "Password file deleted."
        else
            print_warning "Password file kept at: $temp_password_file"
            echo "Remember to delete it securely when no longer needed!"
        fi
    fi

    # Generate summary
    {
        echo "Wallet Generation Summary"
        echo "========================="
        echo ""
        echo "Timestamp:       $(date -u +"%Y-%m-%d %H:%M:%S UTC")"
        echo "Duration:        ${duration}s"
        echo "Wallet count:    $count"
        echo "Address prefix:  $prefix"
        echo ""
        echo "Output Mode:"
        if [[ "$use_plaintext" == true ]] && [[ "$use_encrypt" == true ]]; then
            echo "  - Plain text + Encrypted keystores"
        elif [[ "$use_encrypt" == true ]]; then
            echo "  - Encrypted keystores only (recommended)"
        else
            echo "  - Plain text only (UNSAFE)"
        fi
        echo ""
        echo "Output Files:"
        if [[ "$use_plaintext" == true ]]; then
            echo "  - Plain text:    $output_path"
        fi
        if [[ "$use_encrypt" == true ]]; then
            if [[ "$use_plaintext" == true ]]; then
                echo "  - Encrypted:     $encrypted_output_path"
            else
                echo "  - Wallet index:  $output_path"
            fi
            echo "  - Keystores:     $keystore_dir/"
        fi
        echo ""
        echo "Security Notes:"
        echo "  - Output files have 0600 permissions (owner read/write only)"
        echo "  - Keystore directory has 0700 permissions (owner access only)"
        if [[ "$use_plaintext" == true ]]; then
            echo "  - WARNING: Plain text output contains unencrypted private keys!"
            echo "    Consider deleting after importing to secure storage."
        fi
        echo ""
        echo "Generated by: wallet-gen.sh"
    } > "$summary_file"
    secure_file "$summary_file"

    # Print completion message
    echo ""
    echo -e "${GREEN}${BOLD}======================================${NC}"
    echo -e "${GREEN}${BOLD}      Generation Complete!${NC}"
    echo -e "${GREEN}${BOLD}======================================${NC}"
    echo ""
    echo "Generated $count wallets in ${duration}s"
    echo ""
    echo "Output location: $output_dir/"
    echo ""
    echo "Files created:"

    if [[ "$use_plaintext" == true ]]; then
        local plain_size
        plain_size=$(wc -c < "$output_path" 2>/dev/null || echo "?")
        echo "  - $output_file ($plain_size bytes)"
    fi

    if [[ "$use_encrypt" == true ]]; then
        if [[ "$use_plaintext" == true ]]; then
            local enc_size
            enc_size=$(wc -c < "$encrypted_output_path" 2>/dev/null || echo "?")
            echo "  - ${output_file%.txt}-encrypted.txt ($enc_size bytes)"
        else
            local out_size
            out_size=$(wc -c < "$output_path" 2>/dev/null || echo "?")
            echo "  - $output_file ($out_size bytes)"
        fi
        local ks_count
        ks_count=$(find "$keystore_dir" -name "*.json" 2>/dev/null | wc -l | tr -d ' ')
        echo "  - keystores/ ($ks_count keystore files)"
    fi

    echo "  - run-summary.txt"
    echo ""

    if [[ "$use_plaintext" == true ]]; then
        print_warning "Plain text output contains unencrypted private keys!"
        echo "Consider secure deletion after importing to secure storage."
        echo ""
    fi

    print_success "Done! Check $summary_file for full details."
}

# =============================================================================
# Main Entry Point
# =============================================================================

main() {
    case "${1:-}" in
        --help|-h)
            show_help
            exit 0
            ;;
        --version)
            echo "wallet-gen.sh v1.0.0"
            exit 0
            ;;
        "")
            run_interactive
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help for usage information."
            exit 1
            ;;
    esac
}

main "$@"
