package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

// Token represents an authentication token
type Token struct {
	ID          string    `json:"id"`
	TokenHash   string    `json:"token_hash"`
	User        string    `json:"user"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Revoked     bool      `json:"revoked"`
}

// TokenStore holds all tokens
type TokenStore struct {
	Tokens []Token `json:"tokens"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create":
		createCommand()
	case "list":
		listCommand()
	case "revoke":
		revokeCommand()
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func createCommand() {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	user := fs.String("user", "", "username (required)")
	perms := fs.String("permissions", "upload,download,list", "comma-separated permissions")
	days := fs.Int("days", 365, "days until expiration")
	tokenFile := fs.String("file", "tokens.json", "tokens file path")

	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *user == "" {
		fmt.Println("Error: --user is required")
		fs.PrintDefaults()
		os.Exit(1)
	}

	// Generate token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		os.Exit(1)
	}
	tokenStr := hex.EncodeToString(tokenBytes)

	// Hash token for storage
	hash := sha256.Sum256([]byte(tokenStr))
	tokenHash := hex.EncodeToString(hash[:])

	// Parse permissions
	var permissions []string
	if *perms != "" {
		permissions = parsePermissions(*perms)
	}

	// Create token
	token := Token{
		ID:          fmt.Sprintf("tok_%s", tokenStr[:12]),
		TokenHash:   tokenHash,
		User:        *user,
		Permissions: permissions,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().AddDate(0, 0, *days),
		Revoked:     false,
	}

	// Load existing tokens
	store := loadTokenStore(*tokenFile)
	store.Tokens = append(store.Tokens, token)

	// Save tokens
	if err := saveTokenStore(*tokenFile, store); err != nil {
		fmt.Printf("Error saving tokens: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Token created successfully!\n\n")
	fmt.Printf("Token ID:     %s\n", token.ID)
	fmt.Printf("Token:        %s\n", tokenStr)
	fmt.Printf("User:         %s\n", token.User)
	fmt.Printf("Permissions:  %v\n", token.Permissions)
	fmt.Printf("Expires:      %s\n", token.ExpiresAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("\n⚠️  Save this token! It won't be shown again.\n")
}

func listCommand() {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	tokenFile := fs.String("file", "tokens.json", "tokens file path")
	showRevoked := fs.Bool("revoked", false, "show revoked tokens")

	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	store := loadTokenStore(*tokenFile)

	if len(store.Tokens) == 0 {
		fmt.Println("No tokens found.")
		return
	}

	fmt.Printf("%-15s %-15s %-30s %-10s %-20s\n", "ID", "User", "Permissions", "Status", "Expires")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────────────")

	for _, token := range store.Tokens {
		if token.Revoked && !*showRevoked {
			continue
		}

		status := "active"
		if token.Revoked {
			status = "revoked"
		} else if time.Now().After(token.ExpiresAt) {
			status = "expired"
		}

		permsStr := fmt.Sprintf("%v", token.Permissions)
		if len(permsStr) > 30 {
			permsStr = permsStr[:27] + "..."
		}

		fmt.Printf("%-15s %-15s %-30s %-10s %-20s\n",
			token.ID,
			token.User,
			permsStr,
			status,
			token.ExpiresAt.Format("2006-01-02 15:04"),
		)
	}
}

func revokeCommand() {
	fs := flag.NewFlagSet("revoke", flag.ExitOnError)
	tokenFile := fs.String("file", "tokens.json", "tokens file path")

	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if len(fs.Args()) < 1 {
		fmt.Println("Error: token ID required")
		fmt.Println("Usage: goflux-admin revoke <token-id>")
		os.Exit(1)
	}

	tokenID := fs.Args()[0]

	store := loadTokenStore(*tokenFile)
	found := false

	for i := range store.Tokens {
		if store.Tokens[i].ID == tokenID {
			if store.Tokens[i].Revoked {
				fmt.Printf("Token %s is already revoked.\n", tokenID)
				return
			}
			store.Tokens[i].Revoked = true
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Token not found: %s\n", tokenID)
		os.Exit(1)
	}

	if err := saveTokenStore(*tokenFile, store); err != nil {
		fmt.Printf("Error saving tokens: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Token %s has been revoked.\n", tokenID)
}

func loadTokenStore(filename string) *TokenStore {
	store := &TokenStore{Tokens: []Token{}}

	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return store
		}
		fmt.Printf("Error reading token file: %v\n", err)
		os.Exit(1)
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, store); err != nil {
			fmt.Printf("Error parsing token file: %v\n", err)
			os.Exit(1)
		}
	}

	return store
}

func saveTokenStore(filename string, store *TokenStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0600)
}

func parsePermissions(perms string) []string {
	var result []string
	current := ""
	for _, char := range perms {
		if char == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func printUsage() {
	fmt.Println("goflux-admin - Token management for goflux server")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  goflux-admin <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  create    Create a new authentication token")
	fmt.Println("  list      List all tokens")
	fmt.Println("  revoke    Revoke a token")
	fmt.Println("  help      Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  goflux-admin create --user alice --permissions upload,download")
	fmt.Println("  goflux-admin list")
	fmt.Println("  goflux-admin list --revoked")
	fmt.Println("  goflux-admin revoke tok_abc123def456")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  create:")
	fmt.Println("    --user <username>           User name (required)")
	fmt.Println("    --permissions <perms>       Comma-separated permissions (default: upload,download,list)")
	fmt.Println("    --days <days>               Days until expiration (default: 365)")
	fmt.Println("    --file <path>               Tokens file path (default: tokens.json)")
	fmt.Println()
	fmt.Println("  list:")
	fmt.Println("    --file <path>               Tokens file path (default: tokens.json)")
	fmt.Println("    --revoked                   Show revoked tokens")
	fmt.Println()
	fmt.Println("  revoke:")
	fmt.Println("    --file <path>               Tokens file path (default: tokens.json)")
}
