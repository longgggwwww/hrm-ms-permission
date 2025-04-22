package seeds

import (
	ctx "context"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwww/hrm-ms-permission/ent"
	"github.com/longgggwww/hrm-ms-permission/ent/permgroup"
)

func SeedPerms(ctx ctx.Context, client *ent.Client) error {
	filePath := filepath.Join("data", "perm.csv")
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
		// Assuming record[0] is the permission name and record[1] is the description
		log.Printf("Seeding permission: %s - %s", record[0], record[1])

		// Find the group by code
		group, err := client.PermGroup.
			Query().
			Where(permgroup.Code(record[1])).
			Only(ctx)
		if err != nil {
			log.Printf("Failed to find group for permission %s: %s", record[0], err)
			continue
		}

		err = client.Perm.Create().
			SetCode(record[0]).
			SetName(record[2]).
			SetDescription(record[3]).
			SetGroupID(group.ID).
			OnConflict(sql.ConflictColumns("code")).
			UpdateNewValues().
			Exec(ctx)
		if err != nil {
			log.Printf("Failed to seed permission: %s - %s", record[0], err)
		}
	}

	return nil
}
