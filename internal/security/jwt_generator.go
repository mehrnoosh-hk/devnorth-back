package security

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
)

// jwtGenerator implements domain.TokenGenerator using JWT
type jwtGenerator struct {
	keys          map[string]string
	currentKey    string
	tokenDuration time.Duration
	logger        *slog.Logger
	method        jwt.SigningMethod
	minKeyLength  int
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
func NewJWTGenerator(kid string, Keys map[string]string, durationMinutes int, l *slog.Logger) (domain.TokenGenerator, error) {
	// Check all the keys value in Keys map is acceptable
	// Trim all keys
	for k := range Keys {
		Keys[k] = strings.TrimSpace(Keys[k])
		if len(Keys[k]) < 32 {
			return nil, ErrKeyIsTooShort
		}
	}

	if durationMinutes <= 0 {
		return nil, ErrDurationMustBePositive
	}
	if l == nil {
		return nil, ErrLoggerCanNotBeNil
	}
	return &jwtGenerator{
		keys:          Keys,
		currentKey:    kid,
		tokenDuration: time.Duration(durationMinutes) * time.Minute,
		logger:        l,
		method:        jwt.SigningMethodHS256, // TODO: make it configurable
		minKeyLength:  32, // TODO: make it configurable
	}, nil
}

// Generate creates a JWT token for the given user
func (g *jwtGenerator) Generate(ctx context.Context, user *domain.User) (string, error) {
	if user == nil {
		g.logger.Error("user can not be nil")
		return "", ErrUserCanNotBeNil
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
	token := jwt.NewWithClaims(g.method, claims)
	token.Header["kid"] = g.currentKey

	// Sign token with secret key
	signedToken, err := token.SignedString([]byte(g.keys[g.currentKey]))
	if err != nil {
		g.logger.Error("failed to sign token", "error", err)
		return "", fmt.Errorf("%w: %w", ErrFailedToSignToken, err)
	}

	return signedToken, nil
}

// Validate verifies a JWT token and returns the user claims
func (g *jwtGenerator) Validate(ctx context.Context, tokenString string) (*domain.User, error) {
	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			g.logger.Error("unexpected signing method", "method", token.Header["alg"])
			return nil, ErrUnexpectedSigningMethod
		}
		// Verify kid
		if token.Header["kid"] != g.currentKey {
			g.logger.Error("unexpected key ID", "kid", token.Header["kid"])
			return nil, ErrUnexpectedKeyID
		}
		return []byte(g.keys[g.currentKey]), nil
	})

	if err != nil {
		g.logger.Error("invalid token", "error", err)
		return nil, fmt.Errorf("%w: %w", domain.ErrInvalidToken, err)
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		g.logger.Error("invalid token claims")
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
