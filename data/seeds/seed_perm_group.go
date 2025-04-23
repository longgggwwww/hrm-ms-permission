package seeds

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwww/hrm-ms-permission/ent"
	"github.com/longgggwww/hrm-ms-permission/internal/utils"
)

func SeedPermGroups(ctx context.Context, client *ent.Client) error {
	filePath := filepath.Join("data", "perm_group.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return utils.WrapError("opening file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Read the header row to map column names to indices
	header, err := reader.Read()
	if err != nil {
		return utils.WrapError("reading header row", err)
	}

	headerMap := make(map[string]int)
	for i, col := range header {
		headerMap[col] = i
	}

	records, err := reader.ReadAll()
	if err != nil {
		return utils.WrapError("reading CSV records", err)
	}

	for _, record := range records {
		// Validate required columns
		codeIdx, codeExists := headerMap["code"]
		nameIdx, nameExists := headerMap["name"]

		if !codeExists || !nameExists {
			log.Printf("Skipping record due to missing required columns: %v", record)
			continue
		}

		log.Printf("Seeding PermGroup: %s - %s", record[codeIdx], record[nameIdx])

		permGroup := client.PermGroup.Create().
			SetCode(record[codeIdx]).
			SetName(record[nameIdx]).
			OnConflict(sql.ConflictColumns("code")).
			UpdateNewValues()

		if err := permGroup.Exec(ctx); err != nil {
			log.Printf("Failed to seed PermGroup %s: %v", record[codeIdx], err)
		}
	}

	return nil
}
