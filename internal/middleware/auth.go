package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// GinJWTMiddleware returns Gin middleware for JWT authentication
func (am *AuthMiddleware) GinJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("DEBUG: GinJWTMiddleware invoked for", c.Request.URL.Path)
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			log.Println("DEBUG: Missing Authorization header, aborting with 401")
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "missing authorization header",
				Code:    http.StatusUnauthorized,
			})
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(am.config.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			log.Printf("DEBUG: Invalid JWT token: %v, aborting with 401", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "invalid or missing jwt token",
				Code:    http.StatusUnauthorized,
			})
			return
		}
		if claims, ok := token.Claims.(*models.JWTClaims); ok {
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_name", claims.Name)
		}
		log.Println("DEBUG: JWT valid, proceeding to next handler")
		c.Next()
	}
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

	// Safely extract claims with type assertions
	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'email' claim in Google ID token")
	}
	emailVerified, ok := payload.Claims["email_verified"].(bool)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'email_verified' claim in Google ID token")
	}

	// Optional fields - use safe extraction with fallback to empty string
	name, _ := payload.Claims["name"].(string)
	givenName, _ := payload.Claims["given_name"].(string)
	familyName, _ := payload.Claims["family_name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	locale, _ := payload.Claims["locale"].(string)

	// Extract user information from payload
	userInfo := &models.GoogleUserInfo{
		Sub:           payload.Subject,
		Email:         email,
		EmailVerified: emailVerified,
		Name:          name,
		GivenName:     givenName,
		FamilyName:    familyName,
		Picture:       picture,
		Locale:        locale,
	}

	// Verify email is verified
	if !userInfo.EmailVerified {
		return nil, fmt.Errorf("Google account email not verified")
	}

	return userInfo, nil
}

// RequireAuth returns Gin JWT middleware
func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return am.GinJWTMiddleware()
}

// GetUserFromContext extracts user information from Gin context
func GetUserFromContext(c *gin.Context) (userID, email, name string) {
	if val, exists := c.Get("user_id"); exists {
		if uid, ok := val.(string); ok {
			userID = uid
		}
	}
	if val, exists := c.Get("user_email"); exists {
		if em, ok := val.(string); ok {
			email = em
		}
	}
	if val, exists := c.Get("user_name"); exists {
		if nm, ok := val.(string); ok {
			name = nm
		}
	}
	return userID, email, name
}

// SkipAuth is not implemented for Gin. Use route group configuration for public endpoints.
