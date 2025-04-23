package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/longgggwww/hrm-ms-permission/data/seeds"
	"github.com/longgggwww/hrm-ms-permission/ent"
)

func main() {
	log.Println("Starting seeding process...")

	// Load environment variables
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Initialize database client
	client, err := initDBClient()
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}
	defer client.Close()

	// Run seeders
	if err := runSeeders(context.Background(), client); err != nil {
		log.Fatalf("Seeding process failed: %v", err)
	}

	log.Println("Seeding process completed successfully.")
}

// loadEnv loads environment variables from the .env file.
func loadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}

// initDBClient initializes and returns an Ent database client.
func initDBClient() (*ent.Client, error) {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		return nil, errMissingEnv("DB_URL")
	}

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// runSeeders executes all seed functions.
func runSeeders(ctx context.Context, client *ent.Client) error {
	if err := seeds.SeedPermGroups(ctx, client); err != nil {
		return wrapError("PermGroup", err)
	}

	if err := seeds.SeedPerms(ctx, client); err != nil {
		return wrapError("Perm", err)
	}

	// Add seeding for admin role
	if err := seeds.SeedAdminRole(ctx, client); err != nil {
		return wrapError("AdminRole", err)
	}

	if err := seeds.SeedRoles(ctx, client); err != nil {
		return wrapError("Role", err)
	}

	return nil
}

// errMissingEnv creates an error for missing environment variables.
func errMissingEnv(varName string) error {
	return fmt.Errorf("environment variable %s is not set", varName)
}

// wrapError wraps a seeding error with additional context.
func wrapError(entity string, err error) error {
	return fmt.Errorf("failed to seed %s: %w", entity, err)
}
