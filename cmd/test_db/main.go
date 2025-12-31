package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/mehrnoosh-hk/devnorth-back/db/sqlc"
)

func main() {
	// Get database URL from environment
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ping to verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("✓ Database connection successful")

	// Create queries instance
	queries := sqlc.New(db)
	ctx := context.Background()

	// Generate test data
	testEmail := fmt.Sprintf("test-%d@example.com", time.Now().Unix())

	// Measure CreateUser query duration
	fmt.Println("\nTesting CreateUser query...")
	startTime := time.Now()

	user, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:          testEmail,
		HashedPassword: "$2a$10$abcdefghijklmnopqrstuvwxyz1234567890",
		Role:           sqlc.UserRoleUSER,
	})

	duration := time.Since(startTime)

	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Print results
	fmt.Println("✓ User created successfully")
	fmt.Printf("\n📊 Query Performance:\n")
	fmt.Printf("   Duration: %v (%d ms)\n", duration, duration.Milliseconds())
	fmt.Printf("\n📝 Created User:\n")
	fmt.Printf("   ID: %d\n", user.ID)
	fmt.Printf("   Email: %s\n", user.Email)
	fmt.Printf("   Role: %s\n", user.Role)
	fmt.Printf("   Created At: %s\n", user.CreatedAt.Format(time.RFC3339))

	// Test GetUserByEmail query
	fmt.Println("\nTesting GetUserByEmail query...")
	startTime = time.Now()

	retrievedUser, err := queries.GetUserByEmail(ctx, testEmail)

	duration = time.Since(startTime)

	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}

	fmt.Println("✓ User retrieved successfully")
	fmt.Printf("\n📊 Query Performance:\n")
	fmt.Printf("   Duration: %v (%d ms)\n", duration, duration.Milliseconds())
	fmt.Printf("\n📝 Retrieved User:\n")
	fmt.Printf("   ID: %d\n", retrievedUser.ID)
	fmt.Printf("   Email: %s\n", retrievedUser.Email)
	fmt.Printf("   Matches created user: %v\n", retrievedUser.ID == user.ID)
}
