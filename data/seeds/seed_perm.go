package seeds

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwww/hrm-ms-permission/ent"
	"github.com/longgggwww/hrm-ms-permission/ent/permgroup"
	"github.com/longgggwww/hrm-ms-permission/internal/utils"
)

func SeedPerms(ctx context.Context, client *ent.Client) error {
	filePath := filepath.Join("data", "perm.csv")
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
		groupCodeIdx, groupCodeExists := headerMap["group_code"]
		nameIdx, nameExists := headerMap["name"]
		descIdx, descExists := headerMap["description"]

		if !codeExists || !groupCodeExists || !nameExists || !descExists {
			log.Printf("Skipping record due to missing required columns: %v", record)
			continue
		}

		log.Printf("Seeding permission: %s - %s", record[codeIdx], record[groupCodeIdx])

		group, err := client.PermGroup.
			Query().
			Where(permgroup.Code(record[groupCodeIdx])).
			Only(ctx)
		if err != nil {
			log.Printf("Failed to find group for permission %s: %v", record[codeIdx], err)
			continue
		}

		perm := client.Perm.Create().
			SetCode(record[codeIdx]).
			SetName(record[nameIdx]).
			SetDescription(record[descIdx]).
			SetGroup(group).
			OnConflict(sql.ConflictColumns("code")).
			UpdateNewValues()

		if err := perm.Exec(ctx); err != nil {
			log.Printf("Failed to seed permission %s: %v", record[codeIdx], err)
		}
	}

	return nil
}
