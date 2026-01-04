#!/usr/bin/env bash
#
# wallet-gen.command - macOS double-click launcher for Monolythium Wallet Generator
#
# Double-click this file in Finder to open Terminal and run the guided wallet
# generation script. This provides an easy "double click to start" experience
# for Mac users without requiring terminal knowledge.
#
# Note: This is still a terminal UI (by design). No GUI app signing required.
#

# Exit on error
set -euo pipefail

# =============================================================================
# Determine script location and repo root
# =============================================================================

# Get the directory where this .command file is located
# This works whether double-clicked or run from terminal
SCRIPT_PATH="$0"

# Resolve symlinks to get the real path
if [[ -L "$SCRIPT_PATH" ]]; then
    SCRIPT_PATH="$(readlink "$SCRIPT_PATH")"
fi

# Get the directory containing this script
SCRIPT_DIR="$(cd "$(dirname "$SCRIPT_PATH")" && pwd)"

# The repo root is where this .command file lives
REPO_ROOT="$SCRIPT_DIR"

# =============================================================================
# Terminal setup for double-click experience
# =============================================================================

# Clear terminal for clean look
clear

# Show welcome banner
echo ""
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║                                                                  ║"
echo "║           MONOLYTHIUM WALLET GENERATOR                           ║"
echo "║                                                                  ║"
echo "║   Offline Mass Wallet Generation Tool                           ║"
echo "║                                                                  ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo ""
echo "This tool generates Monolythium wallets completely OFFLINE."
echo "Private keys are NEVER transmitted over the network."
echo ""
echo "Starting guided mode..."
echo ""

# =============================================================================
# Launch the guided script
# =============================================================================

# Check if the guided script exists
GUIDED_SCRIPT="$REPO_ROOT/scripts/wallet-gen.sh"

if [[ ! -f "$GUIDED_SCRIPT" ]]; then
    echo "ERROR: Could not find scripts/wallet-gen.sh"
    echo ""
    echo "Expected location: $GUIDED_SCRIPT"
    echo ""
    echo "Please ensure you have the complete wallet-gen distribution."
    echo ""
    echo "Press Enter to close..."
    read -r
    exit 1
fi

# Ensure script is executable
if [[ ! -x "$GUIDED_SCRIPT" ]]; then
    echo "Making wallet-gen.sh executable..."
    chmod +x "$GUIDED_SCRIPT"
fi

# Change to repo root so relative paths work correctly
cd "$REPO_ROOT"

# Run the guided script
exec "$GUIDED_SCRIPT"
