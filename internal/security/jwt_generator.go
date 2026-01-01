package security

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
)

// jwtGenerator implements domain.TokenGenerator using JWT
type jwtGenerator struct {
	keys          map[string]string
	currentKey    string
	tokenDuration time.Duration
}

// Claims represents the JWT token claims
type Claims struct {
	UserID int32           `json:"user_id"`
	Email  string          `json:"email"`
	Role   domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// NewJWTGenerator creates a new JWT-based token generator
// currentKey: Which key to use for new tokens
// durationMinutes: token expiration time in minutes
func NewJWTGenerator(kid string, Keys map[string]string, durationMinutes int) (domain.TokenGenerator, error) {
	// Check all the keys value in Keys map is acceptable
	for k := range Keys {
		if len(Keys[k]) < 32 {
			return nil, fmt.Errorf("key %s is too short", k)
		}
	}
	if durationMinutes <= 0 {
		return nil, fmt.Errorf("durationMinutes must be positive")
	}
	return &jwtGenerator{
		keys:          Keys,
		currentKey:    kid,
		tokenDuration: time.Duration(durationMinutes) * time.Minute,
	}, nil
}

// Generate creates a JWT token for the given user
func (g *jwtGenerator) Generate(ctx context.Context, user *domain.User) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user cannot be nil")
	}
	now := time.Now()
	expiresAt := now.Add(g.tokenDuration)

	// Create claims with user information
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = g.currentKey

	// Sign token with secret key
	signedToken, err := token.SignedString([]byte(g.keys[g.currentKey]))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// Validate verifies a JWT token and returns the user claims
func (g *jwtGenerator) Validate(ctx context.Context, tokenString string) (*domain.User, error) {
	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Verify kid
		if token.Header["kid"] != g.currentKey {
			return nil, fmt.Errorf("unexpected key id: %v", token.Header["kid"])
		}
		return []byte(g.keys[g.currentKey]), nil
	})

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	// Reconstruct user from claims
	// Note: This is a partial user object from token claims
	// For full user data, query the repository
	user := &domain.User{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	}

	return user, nil
}
