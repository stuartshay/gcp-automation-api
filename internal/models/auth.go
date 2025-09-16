package models

import "github.com/golang-jwt/jwt/v5"

// GoogleUserInfo represents user information from Google OAuth
type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// LoginRequest represents a login request with Google ID token
type LoginRequest struct {
	GoogleIDToken string `json:"google_id_token" validate:"required" binding:"required"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	AccessToken string         `json:"access_token"`
	TokenType   string         `json:"token_type"`
	ExpiresIn   int            `json:"expires_in"`
	UserInfo    GoogleUserInfo `json:"user_info"`
}

// OAuthTokenResponse represents the OAuth2 token exchange response from Google
type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Picture   string `json:"picture,omitempty"`
	GoogleSub string `json:"google_sub,omitempty"`
	jwt.RegisteredClaims
}
