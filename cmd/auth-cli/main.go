package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/stuartshay/gcp-automation-api/internal/config"
	"github.com/stuartshay/gcp-automation-api/internal/models"
	"github.com/stuartshay/gcp-automation-api/internal/services"
)

// StoredCredentials represents the stored authentication data
type StoredCredentials struct {
	AccessToken  string                `json:"access_token"`
	TokenType    string                `json:"token_type"`
	ExpiresAt    time.Time             `json:"expires_at"`
	UserInfo     models.GoogleUserInfo `json:"user_info"`
	RefreshToken string                `json:"refresh_token,omitempty"`
}

var (
	cfg         *config.Config
	authService *services.AuthService
)

func main() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	authService = services.NewAuthService(cfg)

	rootCmd := &cobra.Command{
		Use:   "auth-cli",
		Short: "GCP Automation API Authentication CLI",
		Long:  "A CLI tool for managing authentication with the GCP Automation API",
	}

	rootCmd.AddCommand(
		loginCmd(),
		tokenCmd(),
		refreshCmd(),
		profileCmd(),
		testTokenCmd(),
		logoutCmd(),
		statusCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func loginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Login with Google OAuth",
		Long:  "Perform Google OAuth authentication and store credentials locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			return performGoogleLogin()
		},
	}
}

func tokenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "token",
		Short: "Display current access token",
		Long:  "Show the current JWT access token if available",
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, err := loadCredentials()
			if err != nil {
				return fmt.Errorf("no valid credentials found. Please run 'auth-cli login' first")
			}

			if time.Now().After(creds.ExpiresAt) {
				return fmt.Errorf("token has expired. Please run 'auth-cli refresh' or 'auth-cli login'")
			}

			fmt.Println(creds.AccessToken)
			return nil
		},
	}
}

func refreshCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "refresh",
		Short: "Refresh the current token",
		Long:  "Refresh the current JWT token to extend its expiration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return refreshToken()
		},
	}
}

func profileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "profile",
		Short: "Show user profile",
		Long:  "Display the current authenticated user's profile information",
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, err := loadCredentials()
			if err != nil {
				return fmt.Errorf("no valid credentials found. Please run 'auth-cli login' first")
			}

			fmt.Printf("User Profile:\n")
			fmt.Printf("  Name: %s\n", creds.UserInfo.Name)
			fmt.Printf("  Email: %s\n", creds.UserInfo.Email)
			fmt.Printf("  ID: %s\n", creds.UserInfo.Sub)
			if creds.UserInfo.Picture != "" {
				fmt.Printf("  Picture: %s\n", creds.UserInfo.Picture)
			}
			fmt.Printf("  Token Expires: %s\n", creds.ExpiresAt.Format(time.RFC3339))

			return nil
		},
	}
}

func testTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test-token",
		Short: "Generate test token (development only)",
		Long:  "Generate a test JWT token for development and testing purposes",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.IsProduction() {
				return fmt.Errorf("test token generation is not allowed in production")
			}

			userID, _ := cmd.Flags().GetString("user-id")
			email, _ := cmd.Flags().GetString("email")
			name, _ := cmd.Flags().GetString("name")

			if userID == "" || email == "" || name == "" {
				return fmt.Errorf("user-id, email, and name are required")
			}

			token, err := authService.GenerateTestJWT(userID, email, name)
			if err != nil {
				return fmt.Errorf("failed to generate test token: %w", err)
			}

			// Store the test credentials
			creds := &StoredCredentials{
				AccessToken: token,
				TokenType:   "Bearer",
				ExpiresAt:   time.Now().Add(time.Duration(cfg.JWTExpirationHours) * time.Hour),
				UserInfo: models.GoogleUserInfo{
					Sub:   userID,
					Email: email,
					Name:  name,
				},
			}

			if err := saveCredentials(creds); err != nil {
				return fmt.Errorf("failed to save credentials: %w", err)
			}

			fmt.Printf("Test token generated and saved successfully\n")
			fmt.Printf("Token: %s\n", token)
			return nil
		},
	}

	cmd.Flags().String("user-id", "", "User ID for test token")
	cmd.Flags().String("email", "", "Email for test token")
	cmd.Flags().String("name", "", "Name for test token")
	_ = cmd.MarkFlagRequired("user-id")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func logoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Clear stored credentials",
		Long:  "Remove all stored authentication credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			credPath := getCredentialsPath()
			if err := os.Remove(credPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove credentials: %w", err)
			}
			fmt.Println("Logged out successfully")
			return nil
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		Long:  "Display current authentication status and token information",
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, err := loadCredentials()
			if err != nil {
				fmt.Println("Status: Not authenticated")
				fmt.Println("Run 'auth-cli login' to authenticate")
				return nil
			}

			fmt.Println("Status: Authenticated")
			fmt.Printf("User: %s (%s)\n", creds.UserInfo.Name, creds.UserInfo.Email)
			fmt.Printf("Token Type: %s\n", creds.TokenType)
			fmt.Printf("Expires: %s\n", creds.ExpiresAt.Format(time.RFC3339))

			if time.Now().After(creds.ExpiresAt) {
				fmt.Println("⚠️  Token has expired. Run 'auth-cli refresh' or 'auth-cli login'")
			} else {
				remaining := time.Until(creds.ExpiresAt)
				fmt.Printf("Time remaining: %s\n", remaining.Round(time.Minute))
			}

			return nil
		},
	}
}

func performGoogleLogin() error {
	if cfg.GoogleClientID == "" {
		return fmt.Errorf("GOOGLE_CLIENT_ID not configured")
	}

	// Generate state parameter for security
	state, err := generateRandomString(32)
	if err != nil {
		return fmt.Errorf("failed to generate state parameter: %w", err)
	}

	// Build OAuth URL
	authURL := buildGoogleAuthURL(state)

	// Start local server to handle callback
	server := &http.Server{
		Addr:              ":" + cfg.OAuthCallbackPort,
		ReadHeaderTimeout: 10 * time.Second,
	}
	var authCode string
	var authError error

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			authError = fmt.Errorf("invalid state parameter")
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		if errParam := r.URL.Query().Get("error"); errParam != "" {
			authError = fmt.Errorf("OAuth error: %s", errParam)
			http.Error(w, "OAuth error: "+errParam, http.StatusBadRequest)
			return
		}

		authCode = r.URL.Query().Get("code")
		if authCode == "" {
			authError = fmt.Errorf("no authorization code received")
			http.Error(w, "No authorization code received", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head><title>Authentication Successful</title></head>
<body>
<h1>Authentication Successful!</h1>
<p>You can close this window and return to the terminal.</p>
<script>window.close();</script>
</body>
</html>`)

		// Shutdown server
		go func() {
			time.Sleep(1 * time.Second)
			_ = server.Shutdown(context.Background())
		}()
	})

	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			authError = fmt.Errorf("failed to start callback server: %w", err)
		}
	}()

	// Open browser
	fmt.Printf("Opening browser for Google authentication...\n")
	fmt.Printf("If the browser doesn't open automatically, visit: %s\n", authURL)

	if err := openBrowser(authURL); err != nil {
		fmt.Printf("Failed to open browser automatically: %v\n", err)
		fmt.Printf("Please manually visit: %s\n", authURL)
	}

	// Wait for callback or timeout
	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			_ = server.Shutdown(context.Background())
			return fmt.Errorf("authentication timeout")
		case <-ticker.C:
			if authError != nil {
				return authError
			}
			if authCode != "" {
				// Exchange code for token
				return exchangeCodeForToken(authCode)
			}
		}
	}
}

func buildGoogleAuthURL(state string) string {
	baseURL := "https://accounts.google.com/o/oauth2/v2/auth"
	params := url.Values{
		"client_id":     {cfg.GoogleClientID},
		"redirect_uri":  {cfg.OAuthRedirectURI},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"state":         {state},
		"access_type":   {"offline"},
		"prompt":        {"consent"},
	}
	return baseURL + "?" + params.Encode()
}

func exchangeCodeForToken(code string) error {
	// Exchange authorization code for tokens
	data := url.Values{
		"client_id":     {cfg.GoogleClientID},
		"client_secret": {cfg.GoogleClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {cfg.OAuthRedirectURI},
	}

	resp, err := http.PostForm(cfg.OAuthTokenURL, data)
	if err != nil {
		return fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	// Use the ID token to authenticate with our service
	loginResp, err := authService.LoginWithGoogle(context.Background(), tokenResp.IDToken)
	if err != nil {
		return fmt.Errorf("failed to authenticate with service: %w", err)
	}

	// Store credentials
	creds := &StoredCredentials{
		AccessToken: loginResp.AccessToken,
		TokenType:   loginResp.TokenType,
		ExpiresAt:   time.Now().Add(time.Duration(loginResp.ExpiresIn) * time.Second),
		UserInfo:    loginResp.UserInfo,
	}

	if err := saveCredentials(creds); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	fmt.Printf("Authentication successful!\n")
	fmt.Printf("Welcome, %s (%s)\n", creds.UserInfo.Name, creds.UserInfo.Email)
	return nil
}

func refreshToken() error {
	creds, err := loadCredentials()
	if err != nil {
		return fmt.Errorf("no credentials found. Please run 'auth-cli login' first")
	}

	// Validate current token to get claims
	claims, err := authService.ValidateJWT(creds.AccessToken)
	if err != nil {
		return fmt.Errorf("current token is invalid. Please run 'auth-cli login' again")
	}

	// Generate new token
	newToken, err := authService.RefreshJWT(claims)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	// Update stored credentials
	creds.AccessToken = newToken
	creds.ExpiresAt = time.Now().Add(time.Duration(cfg.JWTExpirationHours) * time.Hour)

	if err := saveCredentials(creds); err != nil {
		return fmt.Errorf("failed to save refreshed credentials: %w", err)
	}

	fmt.Println("Token refreshed successfully")
	return nil
}

func loadCredentials() (*StoredCredentials, error) {
	credPath := getCredentialsPath()
	if err := validateFilePath(credPath); err != nil {
		return nil, fmt.Errorf("invalid credentials path: %w", err)
	}

	// #nosec G304 - credPath is validated by validateFilePath to prevent path traversal
	data, err := os.ReadFile(credPath)
	if err != nil {
		return nil, err
	}

	var creds StoredCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

func saveCredentials(creds *StoredCredentials) error {
	credPath := getCredentialsPath()

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(credPath), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(credPath, data, 0600)
}

func getCredentialsPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return filepath.Join(".", cfg.CredentialsDir, cfg.CredentialsFile)
	}
	return filepath.Join(homeDir, cfg.CredentialsDir, cfg.CredentialsFile)
}

// validateFilePath ensures the file path is safe and within expected boundaries
func validateFilePath(path string) error {
	// Clean the path to resolve any relative path components
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	// Ensure the path contains expected credentials directory
	if !strings.Contains(cleanPath, cfg.CredentialsDir) {
		return fmt.Errorf("invalid credentials path")
	}

	return nil
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func openBrowser(urlStr string) error {
	// Validate URL to prevent command injection
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Only allow http and https schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}

	// Validate the command exists and add the URL as the last argument
	if _, err := exec.LookPath(cmd); err != nil {
		return fmt.Errorf("browser command not found: %w", err)
	}

	args = append(args, urlStr)
	// #nosec G204 - cmd is validated via exec.LookPath and urlStr is validated URL
	return exec.Command(cmd, args...).Start()
}
