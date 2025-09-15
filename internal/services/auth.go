package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"

	"github.com/stuartshay/gcp-automation-api/internal/config"
	"github.com/stuartshay/gcp-automation-api/internal/models"
)

// AuthService handles authentication operations
type AuthService struct {
	config *config.Config
}

// NewAuthService creates a new authentication service instance
func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		config: cfg,
	}
}

// LoginWithGoogle authenticates a user with Google ID token and returns a JWT
func (as *AuthService) LoginWithGoogle(ctx context.Context, googleIDToken string) (*models.LoginResponse, error) {
	if !as.config.EnableGoogleAuth {
		return nil, fmt.Errorf("Google authentication is disabled")
	}

	// Validate Google ID token
	userInfo, err := as.validateGoogleIDToken(ctx, googleIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate Google ID token: %w", err)
	}

	// Generate JWT token for the user
	jwtToken, err := as.generateJWT(userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	log.Printf("User %s (%s) authenticated successfully", userInfo.Name, userInfo.Email)

	// Prepare response
	response := &models.LoginResponse{
		AccessToken: jwtToken,
		TokenType:   "Bearer",
		ExpiresIn:   as.config.JWTExpirationHours * 3600, // Convert hours to seconds
		UserInfo:    *userInfo,
	}

	return response, nil
}

// ValidateJWT validates a JWT token and returns the claims
func (as *AuthService) ValidateJWT(tokenString string) (*models.JWTClaims, error) {
	// Parse token with custom claims
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(as.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Validate and extract claims
	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		// Check if token is expired
		if time.Now().After(claims.ExpiresAt.Time) {
			return nil, fmt.Errorf("token has expired")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// GenerateTestJWT generates a JWT token for testing purposes (development only)
func (as *AuthService) GenerateTestJWT(userID, email, name string) (string, error) {
	if as.config.IsProduction() {
		return "", fmt.Errorf("test JWT generation is not allowed in production")
	}

	userInfo := &models.GoogleUserInfo{
		Sub:           userID,
		Email:         email,
		Name:          name,
		EmailVerified: true,
		Picture:       "",
		Locale:        "en",
	}

	return as.generateJWT(userInfo)
}

// RefreshJWT generates a new JWT token using existing valid claims
func (as *AuthService) RefreshJWT(claims *models.JWTClaims) (string, error) {
	// Create new user info from existing claims
	userInfo := &models.GoogleUserInfo{
		Sub:           claims.GoogleSub,
		Email:         claims.Email,
		Name:          claims.Name,
		EmailVerified: true,
		Picture:       claims.Picture,
	}

	return as.generateJWT(userInfo)
}

// validateGoogleIDToken validates a Google ID token and extracts user information
func (as *AuthService) validateGoogleIDToken(ctx context.Context, idToken string) (*models.GoogleUserInfo, error) {
	// Validate the Google ID token
	payload, err := idtoken.Validate(ctx, idToken, as.config.GoogleClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate Google ID token: %w", err)
	}

	// Helper function to safely extract string from claims
	getString := func(claims map[string]interface{}, key string) string {
		if val, ok := claims[key]; ok {
			if str, ok := val.(string); ok {
				return str
			}
		}
		return ""
	}

	// Helper function to safely extract bool from claims
	getBool := func(claims map[string]interface{}, key string) bool {
		if val, ok := claims[key]; ok {
			if b, ok := val.(bool); ok {
				return b
			}
		}
		return false
	}

	// Extract user information from payload
	userInfo := &models.GoogleUserInfo{
		Sub:           payload.Subject,
		Email:         getString(payload.Claims, "email"),
		EmailVerified: getBool(payload.Claims, "email_verified"),
		Name:          getString(payload.Claims, "name"),
		GivenName:     getString(payload.Claims, "given_name"),
		FamilyName:    getString(payload.Claims, "family_name"),
		Picture:       getString(payload.Claims, "picture"),
		Locale:        getString(payload.Claims, "locale"),
	}

	// Verify required fields
	if userInfo.Email == "" {
		return nil, fmt.Errorf("email not found in Google ID token")
	}

	if !userInfo.EmailVerified {
		return nil, fmt.Errorf("Google account email not verified")
	}

	return userInfo, nil
}

// generateJWT generates a new JWT token with user information
func (as *AuthService) generateJWT(userInfo *models.GoogleUserInfo) (string, error) {
	// Set token expiration
	expirationTime := time.Now().Add(time.Duration(as.config.JWTExpirationHours) * time.Hour)

	// Create claims
	claims := &models.JWTClaims{
		UserID:    userInfo.Sub,
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		Picture:   userInfo.Picture,
		GoogleSub: userInfo.Sub,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gcp-automation-api",
			Subject:   userInfo.Sub,
			Audience:  []string{"gcp-automation-api"},
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate signed token string
	tokenString, err := token.SignedString([]byte(as.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, nil
}

// GetUserContext extracts user information that can be used in API handlers
func (as *AuthService) GetUserContext(claims *models.JWTClaims) map[string]interface{} {
	return map[string]interface{}{
		"user_id":    claims.UserID,
		"email":      claims.Email,
		"name":       claims.Name,
		"picture":    claims.Picture,
		"google_sub": claims.GoogleSub,
	}
}

// IsTokenExpired checks if a JWT token is expired
func (as *AuthService) IsTokenExpired(claims *models.JWTClaims) bool {
	return time.Now().After(claims.ExpiresAt.Time)
}

// GetConfig returns the configuration for external access
func (as *AuthService) GetConfig() *config.Config {
	return as.config
}
