package main

import (
	"log"

	"context"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/longgggwww/hrm-ms-permission/data/seeds"
	"github.com/longgggwww/hrm-ms-permission/ent"
)

func main() {
	log.Println("Starting seeding process...")

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
		return
	}

	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("DB_URL environment variable is not set")
		return
	}

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	if err := seeds.SeedPermGroups(ctx, client); err != nil {
		log.Fatal("Failed to seed PermGroup:", err)
		return
	}

	if err := seeds.SeedPerms(ctx, client); err != nil {
		log.Fatal("Failed to seed Perm:", err)
		return
	}
}
