// Package main provides the walletgen CLI tool.
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/monolythium/wallet-gen/pkg/keystore"
	"github.com/monolythium/wallet-gen/pkg/wallet"
	"github.com/spf13/cobra"
)

var (
	// Flags
	count        int
	outputFile   string
	prefix       string
	quiet        bool
	seed         int64
	format       string
	encrypt      bool
	keystoreDir  string
	passwordFile string
	scryptN      int
	scryptR      int
	scryptP      int
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "walletgen",
		Short: "Offline mass wallet generator for Monolythium",
		Long: `walletgen generates multiple wallets offline and writes them to a file.

Default output format: bech32:0x:privatekey
Encrypt mode writes Ethereum Keystore V3 JSON files.

This tool runs completely offline with no network calls.`,
		RunE: run,
	}

	// Required flags
	rootCmd.Flags().IntVarP(&count, "count", "n", 0, "Number of wallets to generate (required)")
	rootCmd.Flags().StringVarP(&outputFile, "out", "o", "", "Output file path (required)")
	rootCmd.MarkFlagRequired("count")
	rootCmd.MarkFlagRequired("out")

	// Optional flags
	rootCmd.Flags().StringVar(&prefix, "prefix", "mono", "Bech32 address prefix")
	rootCmd.Flags().BoolVar(&quiet, "quiet", false, "Disable progress output")
	rootCmd.Flags().Int64Var(&seed, "seed", 0, "Deterministic seed (testing only, 0 = random)")
	rootCmd.Flags().StringVar(&format, "format", "bech32:0x:privkey", "Output format (informational)")

	// Encrypt mode flags
	rootCmd.Flags().BoolVar(&encrypt, "encrypt", false, "Enable encrypted keystore mode")
	rootCmd.Flags().StringVar(&keystoreDir, "keystore-dir", "", "Directory for keystore files (required if --encrypt)")
	rootCmd.Flags().StringVar(&passwordFile, "password-file", "", "File containing encryption password (required if --encrypt)")
	rootCmd.Flags().IntVar(&scryptN, "scrypt-n", 262144, "Scrypt N parameter (CPU/memory cost)")
	rootCmd.Flags().IntVar(&scryptR, "scrypt-r", 8, "Scrypt R parameter (block size)")
	rootCmd.Flags().IntVar(&scryptP, "scrypt-p", 1, "Scrypt P parameter (parallelization)")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Validate count
	if count <= 0 {
		return fmt.Errorf("count must be greater than 0")
	}

	// Validate encrypt mode requirements
	if encrypt {
		if keystoreDir == "" {
			return fmt.Errorf("--keystore-dir is required when using --encrypt")
		}
		if passwordFile == "" {
			return fmt.Errorf("--password-file is required when using --encrypt")
		}
	}

	// Read password if encrypt mode
	var password string
	if encrypt {
		data, err := os.ReadFile(passwordFile)
		if err != nil {
			return fmt.Errorf("failed to read password file: %w", err)
		}
		password = strings.TrimSpace(string(data))
		if password == "" {
			return fmt.Errorf("password file is empty")
		}
	}

	// Create output file with restrictive permissions
	file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Use buffered writer for streaming
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Create wallet generator
	var gen *wallet.Generator
	if seed != 0 {
		gen = wallet.NewGenerator(prefix, wallet.DeterministicRand(seed))
	} else {
		gen = wallet.NewGenerator(prefix, nil)
	}

	// Create keystore directory if encrypt mode
	if encrypt {
		if err := os.MkdirAll(keystoreDir, 0700); err != nil {
			return fmt.Errorf("failed to create keystore directory: %w", err)
		}
	}

	// Scrypt params for encrypt mode
	scryptParams := keystore.ScryptParams{
		N: scryptN,
		R: scryptR,
		P: scryptP,
	}

	// Progress tracking
	lastPercent := -1
	updateInterval := count / 100
	if updateInterval < 1 {
		updateInterval = 1
	}

	// Generate wallets
	for i := 0; i < count; i++ {
		w, err := gen.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate wallet %d: %w", i, err)
		}

		var line string
		if encrypt {
			// Encrypt and save keystore
			ksData, err := keystore.EncryptKey(w.PrivateKey, password, scryptParams)
			if err != nil {
				return fmt.Errorf("failed to encrypt wallet %d: %w", i, err)
			}

			filename, err := keystore.WriteKeystore(keystoreDir, i, ksData, w.EVMAddress)
			if err != nil {
				return fmt.Errorf("failed to write keystore %d: %w", i, err)
			}

			line = w.FormatLineEncrypted(filepath.Join(keystoreDir, filename))
		} else {
			line = w.FormatLine()
		}

		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write wallet %d: %w", i, err)
		}

		// Update progress
		if !quiet {
			percent := (i + 1) * 100 / count
			if percent != lastPercent && (i+1)%updateInterval == 0 {
				fmt.Fprintf(os.Stderr, "\rProgress: %d%%", percent)
				lastPercent = percent
			}
		}
	}

	// Final flush
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush output: %w", err)
	}

	if !quiet {
		fmt.Fprintf(os.Stderr, "\rProgress: 100%%\n")
		fmt.Fprintf(os.Stderr, "Generated %d wallets to %s\n", count, outputFile)
		if encrypt {
			fmt.Fprintf(os.Stderr, "Keystores written to %s/\n", keystoreDir)
		}
	}

	return nil
}
