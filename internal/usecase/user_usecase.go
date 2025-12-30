package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mehrnoosh-hk/devnorth-back/internal/domain"
)

// userUseCase implements domain.UserUseCase
// It orchestrates user-related business operations using repository, password hasher, and token generator
type userUseCase struct {
	userRepo       domain.UserRepository
	passwordHasher domain.PasswordHasher
	tokenGenerator domain.TokenGenerator
}

// NewUserUseCase creates a new user use case instance
// Dependencies are injected following the Dependency Inversion Principle
func NewUserUseCase(
	userRepo domain.UserRepository,
	passwordHasher domain.PasswordHasher,
	tokenGenerator domain.TokenGenerator,
) (domain.UserUseCase, error) {
	// Nil-check the injected dependencies
	if userRepo == nil {
		return nil, fmt.Errorf("user repository cannot be nil")
	}
	if passwordHasher == nil {
		return nil, fmt.Errorf("password hasher cannot be nil")
	}
	if tokenGenerator == nil {
		return nil, fmt.Errorf("token generator cannot be nil")
	}
	return &userUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}, nil
}

// Register creates a new user account
// Business logic flow:
// 1. Validate email (basic validation for POC)
// 2. Validate password (basic validation for POC)
// 3. Check if email already exists
// 4. Hash password
// 5. Create user in repository
func (uc *userUseCase) Register(ctx context.Context, email, password string) (*domain.User, error) {
	// Normalize email by trimming spaces
	email = strings.TrimSpace(email)

	// Step 1: Validate email (basic validation for POC)
	if err := uc.validateEmail(email); err != nil {
		return nil, err
	}

	// Step 2: Validate password (basic validation for POC)
	if err := uc.validatePassword(password); err != nil {
		return nil, err
	}

	// Step 3: Check if email already exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	// Step 4: Hash password
	hashedPassword, err := uc.passwordHasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Step 5: Create user with default role (USER)
	user, err := uc.userRepo.Create(ctx, email, hashedPassword, domain.UserRoleUSER)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// validateEmail performs basic email validation
// For POC: minimal validation - just check it's not empty and contains @
// Production: use proper email validation library or regex
func (uc *userUseCase) validateEmail(email string) error {
	email = strings.TrimSpace(email)

	// Basic validation for POC
	if email == "" {
		return domain.ErrInvalidEmail
	}

	// Minimal check: must contain @
	if !strings.Contains(email, "@") {
		return domain.ErrInvalidEmail
	}

	// TODO for production:
	// - Use proper email validation regex or library
	// - Check email length limits
	// - Validate domain exists (optional)

	return nil
}

// validatePassword performs basic password validation
// For POC: minimal validation - just check it's not empty and has minimum length
// Production: enforce complexity requirements
func (uc *userUseCase) validatePassword(password string) error {
	// Basic validation for POC
	if password == "" || utf8.RuneCountInString(password) < 8 || len(password) > 72 {
		return domain.ErrInvalidPassword
	}

	// TODO for production:
	// - Require mix of uppercase, lowercase, numbers, special characters
	// - Check against common password lists
	// - Check password strength score

	return nil
}

// Login authenticates a user and returns a JWT token
// Business logic flow:
// 1. Get user by email
// 2. Verify user exists
// 3. Compare password hash
// 4. Generate JWT token
// 5. Return token + user (without password)
func (uc *userUseCase) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
	// Step 1 & 2: Get user by email
	email = strings.TrimSpace(email)
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Fake hash compare to prevent timing attack
		dummyHash := "$2y$10$QqDjvtHjrzwxwjQJrIwGFuLOPKhVOe0.67k2/Hl2IdOYjnQobQh/i" // bycript hash of "dummy"
		_ = uc.passwordHasher.Compare("someFakePassword", dummyHash)
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", nil, domain.ErrInvalidCredentials
		}
		return "", nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Step 3: Compare password hash
	err = uc.passwordHasher.Compare(user.HashedPassword, password)
	if err != nil {
		// Password doesn't match - return generic error
		return "", nil, domain.ErrInvalidCredentials
	}

	// Step 4: Generate JWT token
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	token, err := uc.tokenGenerator.Generate(ctxWithTimeout, user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Step 5: Return token and user (password is already hashed, but good practice to not return it)
	return token, user, nil
}
