package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Driver Postgres
	"github.com/longgggwwww/hrm-ms-permission/data/seeds"
	"github.com/longgggwwww/hrm-ms-permission/ent"
)

func main() {
	log.Println("Starting seeding process...")

	client, err := initDBClient()
	if err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}
	defer client.Close()

	seeders := []struct {
		name string
		fn   func(context.Context, *ent.Client) error
	}{
		{"perm groups", seeds.SeedPermGroups},
		{"perms", seeds.SeedPerms},
		{"admin role", seeds.SeedAdminRole},
		{"roles", seeds.SeedRoles},
		{"user roles", seeds.SeedUserRoles},
		{"user perms", seeds.SeedUserPerms},
	}

	ctx := context.Background()
	for _, seeder := range seeders {
		log.Printf("Seeding %s...", seeder.name)
		if err := seeder.fn(ctx, client); err != nil {
			log.Fatalf("Failed to seed %s: %v", seeder.name, err)
		}
	}

	log.Println("Seeding process completed successfully.")
}

func initDBClient() (*ent.Client, error) {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		return nil, fmt.Errorf("missing DB_URL environment variable")
	}

	client, err := ent.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return client, nil
}
