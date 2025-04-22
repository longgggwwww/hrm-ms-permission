package seeds

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwww/hrm-ms-permission/ent"
)

func SeedPermGroups(ctx context.Context, client *ent.Client) error {
	filePath := filepath.Join("data", "perm_group.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip the header row
	if _, err := reader.Read(); err != nil {
		return err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record) < 2 {
			log.Printf("Invalid record: %v", record)
			continue
		}

		log.Printf("Seeding PermGroup: %s - %s", record[0], record[1])

		err := client.PermGroup.Create().
			SetCode(record[0]).
			SetName(record[1]).
			OnConflict(sql.ConflictColumns("code")).
			UpdateNewValues().
			Exec(ctx)
		if err != nil {
			log.Printf("Failed to seed PermGroup: %s - %s", record[0], err)
		}
	}

	return nil
}
