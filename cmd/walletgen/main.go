// Package main provides the walletgen CLI tool.
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

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

This tool runs completely offline with no network calls.

Run with no arguments for interactive wizard mode.`,
		RunE: run,
	}

	// Flags (not required - wizard mode handles missing values)
	rootCmd.Flags().IntVarP(&count, "count", "n", 0, "Number of wallets to generate")
	rootCmd.Flags().StringVarP(&outputFile, "out", "o", "", "Output file path")

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
	// Check if we should run wizard mode (no essential flags provided)
	if count == 0 && outputFile == "" {
		return runWizard()
	}

	// CLI mode - validate required params
	if count <= 0 {
		return fmt.Errorf("--count (-n) is required and must be greater than 0")
	}
	if outputFile == "" {
		return fmt.Errorf("--out (-o) is required")
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

	return generateWallets(count, outputFile, encrypt, keystoreDir, password)
}

// getDefaultOutputDir returns the OS-appropriate default output directory
func getDefaultOutputDir() string {
	var baseDir string
	if runtime.GOOS == "windows" {
		// Windows: %USERPROFILE%\wallet-gen-output\
		baseDir = os.Getenv("USERPROFILE")
		if baseDir == "" {
			baseDir = "C:\\"
		}
	} else {
		// macOS/Linux: ~/wallet-gen-output/
		baseDir = os.Getenv("HOME")
		if baseDir == "" {
			baseDir = "."
		}
	}
	return filepath.Join(baseDir, "wallet-gen-output")
}

// runWizard runs the interactive wizard mode
func runWizard() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("===========================================")
	fmt.Println("    Monolythium Wallet Generator Wizard")
	fmt.Println("===========================================")
	fmt.Println()
	fmt.Println("This tool generates wallets OFFLINE.")
	fmt.Println("No network connection is used.")
	fmt.Println()

	// Step 1: Number of wallets
	var wizardCount int
	for {
		fmt.Print("How many wallets do you want to generate? [1-100000]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("  Please enter a number.")
			continue
		}
		n, err := strconv.Atoi(input)
		if err != nil || n < 1 || n > 100000 {
			fmt.Println("  Please enter a number between 1 and 100000.")
			continue
		}
		wizardCount = n
		break
	}
	fmt.Println()

	// Step 2: Output directory
	defaultDir := getDefaultOutputDir()
	timestamp := time.Now().Format("20060102-150405")
	defaultPath := filepath.Join(defaultDir, timestamp)

	fmt.Printf("Output directory [%s]: ", defaultPath)
	dirInput, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	dirInput = strings.TrimSpace(dirInput)
	outputDir := defaultPath
	if dirInput != "" {
		outputDir = dirInput
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0700); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	fmt.Println()

	// Step 3: Output mode
	fmt.Println("Output mode:")
	fmt.Println("  1) Raw text (private keys in plain text) - LESS SECURE")
	fmt.Println("  2) Encrypted keystores (password protected) - RECOMMENDED")
	fmt.Println("  3) Both (raw text + encrypted keystores)")
	fmt.Println()

	var mode int
	for {
		fmt.Print("Choose mode [1/2/3]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		input = strings.TrimSpace(input)
		m, err := strconv.Atoi(input)
		if err != nil || m < 1 || m > 3 {
			fmt.Println("  Please enter 1, 2, or 3.")
			continue
		}
		mode = m
		break
	}
	fmt.Println()

	// Security warning for raw mode
	if mode == 1 || mode == 3 {
		fmt.Println("!!! SECURITY WARNING !!!")
		fmt.Println()
		fmt.Println("Raw text mode stores private keys WITHOUT encryption.")
		fmt.Println("Anyone with access to the file can steal ALL funds.")
		fmt.Println()
		fmt.Println("Only use this if you understand the risks and will:")
		fmt.Println("  - Keep the file on an ENCRYPTED, OFFLINE drive")
		fmt.Println("  - Delete it securely after importing keys")
		fmt.Println()

		for {
			fmt.Print("Type YES-I-UNDERSTAND to continue: ")
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			input = strings.TrimSpace(input)
			if input == "YES-I-UNDERSTAND" {
				break
			}
			fmt.Println("  You must type exactly: YES-I-UNDERSTAND")
		}
		fmt.Println()
	}

	// Password for encrypted mode
	var password string
	if mode == 2 || mode == 3 {
		fmt.Println("Enter a password to encrypt your keystores.")
		fmt.Println("IMPORTANT: You MUST remember this password to access your wallets!")
		fmt.Println()

		for {
			fmt.Print("Password (min 8 characters): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			password = strings.TrimSpace(input)
			if len(password) < 8 {
				fmt.Println("  Password must be at least 8 characters.")
				continue
			}

			fmt.Print("Confirm password: ")
			confirm, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			confirm = strings.TrimSpace(confirm)
			if confirm != password {
				fmt.Println("  Passwords do not match. Try again.")
				continue
			}
			break
		}
		fmt.Println()
	}

	// Confirmation summary
	wizardOutputFile := filepath.Join(outputDir, "wallets.txt")
	wizardKeystoreDir := filepath.Join(outputDir, "keystores")

	fmt.Println("===========================================")
	fmt.Println("                 SUMMARY")
	fmt.Println("===========================================")
	fmt.Printf("  Wallets to generate: %d\n", wizardCount)
	fmt.Printf("  Output directory:    %s\n", outputDir)
	switch mode {
	case 1:
		fmt.Printf("  Output file:         %s\n", wizardOutputFile)
		fmt.Println("  Mode:                Raw text (UNENCRYPTED)")
	case 2:
		fmt.Printf("  Output file:         %s\n", wizardOutputFile)
		fmt.Printf("  Keystore directory:  %s\n", wizardKeystoreDir)
		fmt.Println("  Mode:                Encrypted keystores")
	case 3:
		fmt.Printf("  Output file:         %s\n", wizardOutputFile)
		fmt.Printf("  Keystore directory:  %s\n", wizardKeystoreDir)
		fmt.Println("  Mode:                Both (raw + encrypted)")
	}
	fmt.Println("===========================================")
	fmt.Println()

	for {
		fmt.Print("Proceed with generation? [y/n]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "n" || input == "no" {
			fmt.Println("Aborted.")
			return nil
		}
		if input == "y" || input == "yes" {
			break
		}
		fmt.Println("  Please enter y or n.")
	}
	fmt.Println()

	// Generate based on mode
	wizardEncrypt := mode == 2 || mode == 3
	wizardKeystoreDirArg := ""
	if wizardEncrypt {
		wizardKeystoreDirArg = wizardKeystoreDir
	}

	if err := generateWallets(wizardCount, wizardOutputFile, wizardEncrypt, wizardKeystoreDirArg, password); err != nil {
		return err
	}

	// If mode 3, also generate raw version
	if mode == 3 {
		// The generateWallets already outputs with keystore references
		// We need to also generate a raw-only file
		rawFile := filepath.Join(outputDir, "wallets-raw.txt")
		fmt.Println()
		fmt.Println("Generating raw text backup...")
		if err := generateWalletsRaw(wizardCount, rawFile); err != nil {
			return fmt.Errorf("failed to generate raw backup: %w", err)
		}
		fmt.Printf("Raw backup written to %s\n", rawFile)
	}

	fmt.Println()
	fmt.Println("===========================================")
	fmt.Println("              GENERATION COMPLETE")
	fmt.Println("===========================================")
	fmt.Printf("Output directory: %s\n", outputDir)
	fmt.Println()
	fmt.Println("IMPORTANT SECURITY REMINDERS:")
	fmt.Println("  - Back up this folder to a secure location")
	fmt.Println("  - Consider an encrypted USB drive or hardware wallet")
	fmt.Println("  - Never share your private keys or password")
	fmt.Println("===========================================")
	fmt.Println()

	// Keep window open on Windows
	if runtime.GOOS == "windows" {
		fmt.Println("Press Enter to exit...")
		reader.ReadString('\n')
	}

	return nil
}

// generateWalletsRaw generates wallets in raw format (for mode 3 backup)
func generateWalletsRaw(count int, outputFile string) error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(outputFile), 0700); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	gen := wallet.NewGenerator(prefix, nil)

	for i := 0; i < count; i++ {
		w, err := gen.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate wallet %d: %w", i, err)
		}
		if _, err := writer.WriteString(w.FormatLine() + "\n"); err != nil {
			return fmt.Errorf("failed to write wallet %d: %w", i, err)
		}
	}

	return writer.Flush()
}

// generateWallets is the core wallet generation logic
func generateWallets(count int, outputFile string, doEncrypt bool, ksDir string, password string) error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(outputFile), 0700); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
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
	if doEncrypt {
		if err := os.MkdirAll(ksDir, 0700); err != nil {
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
		if doEncrypt {
			// Encrypt and save keystore
			ksData, err := keystore.EncryptKey(w.PrivateKey, password, scryptParams)
			if err != nil {
				return fmt.Errorf("failed to encrypt wallet %d: %w", i, err)
			}

			filename, err := keystore.WriteKeystore(ksDir, i, ksData, w.EVMAddress)
			if err != nil {
				return fmt.Errorf("failed to write keystore %d: %w", i, err)
			}

			line = w.FormatLineEncrypted(filepath.Join(ksDir, filename))
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
		if doEncrypt {
			fmt.Fprintf(os.Stderr, "Keystores written to %s/\n", ksDir)
		}
	}

	return nil
}
