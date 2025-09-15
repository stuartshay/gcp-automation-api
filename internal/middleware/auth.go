package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"

	"github.com/stuartshay/gcp-automation-api/internal/config"
	"github.com/stuartshay/gcp-automation-api/internal/models"
)

// AuthMiddleware provides JWT authentication middleware
type AuthMiddleware struct {
	config *config.Config
}

// NewAuthMiddleware creates a new authentication middleware instance
func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		config: cfg,
	}
}

// JWTMiddleware returns Echo JWT middleware configured for our application
func (am *AuthMiddleware) JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(am.config.JWTSecret),
		TokenLookup: "header:Authorization:Bearer ,query:token",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &models.JWTClaims{}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid or missing JWT token",
				Code:    http.StatusUnauthorized,
			})
		},
		SuccessHandler: func(c echo.Context) {
			// Extract user information from JWT claims and add to context
			token := c.Get("user").(*jwt.Token)
			if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Set("user_name", claims.Name)
			}
		},
	})
}

// GenerateJWT generates a new JWT token with user information
func (am *AuthMiddleware) GenerateJWT(userID, email, name, picture, googleSub string) (string, error) {
	// Set token expiration
	expirationTime := time.Now().Add(time.Duration(am.config.JWTExpirationHours) * time.Hour)

	// Create claims
	claims := &models.JWTClaims{
		UserID:    userID,
		Email:     email,
		Name:      name,
		Picture:   picture,
		GoogleSub: googleSub,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gcp-automation-api",
			Subject:   userID,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate signed token string
	tokenString, err := token.SignedString([]byte(am.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token and returns the claims
func (am *AuthMiddleware) ValidateJWT(tokenString string) (*models.JWTClaims, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(am.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Validate claims
	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// ValidateGoogleIDToken validates a Google ID token and extracts user information
func (am *AuthMiddleware) ValidateGoogleIDToken(ctx context.Context, idToken string) (*models.GoogleUserInfo, error) {
	if !am.config.EnableGoogleAuth {
		return nil, fmt.Errorf("Google authentication is disabled")
	}

	// Validate the Google ID token
	payload, err := idtoken.Validate(ctx, idToken, am.config.GoogleClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate Google ID token: %w", err)
	}

	// Extract user information from payload
	userInfo := &models.GoogleUserInfo{
		Sub:           payload.Subject,
		Email:         payload.Claims["email"].(string),
		EmailVerified: payload.Claims["email_verified"].(bool),
		Name:          payload.Claims["name"].(string),
		GivenName:     payload.Claims["given_name"].(string),
		FamilyName:    payload.Claims["family_name"].(string),
		Picture:       payload.Claims["picture"].(string),
		Locale:        payload.Claims["locale"].(string),
	}

	// Verify email is verified
	if !userInfo.EmailVerified {
		return nil, fmt.Errorf("Google account email not verified")
	}

	return userInfo, nil
}

// RequireAuth is a convenience function that returns the JWT middleware
func (am *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	return am.JWTMiddleware()
}

// GetUserFromContext extracts user information from Echo context
func GetUserFromContext(c echo.Context) (userID, email, name string) {
	if uid, ok := c.Get("user_id").(string); ok {
		userID = uid
	}
	if em, ok := c.Get("user_email").(string); ok {
		email = em
	}
	if nm, ok := c.Get("user_name").(string); ok {
		name = nm
	}
	return userID, email, name
}

// SkipAuth returns a middleware that skips authentication for specific paths
func SkipAuth(paths ...string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			for _, skipPath := range paths {
				if strings.HasPrefix(path, skipPath) {
					return true
				}
			}
			return false
		},
	})
}
