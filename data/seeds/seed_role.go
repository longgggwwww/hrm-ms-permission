package seeds

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwww/hrm-ms-permission/ent"
	"github.com/longgggwww/hrm-ms-permission/ent/perm"
	"github.com/longgggwww/hrm-ms-permission/internal/utils"
)

func SeedRoles(ctx context.Context, client *ent.Client) error {
	filePath := filepath.Join("data", "role.csv")
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
		colorIdx, colorExists := headerMap["color"]
		permCodesIdx, permCodesExists := headerMap["perm_codes"]
		descIdx, descExists := headerMap["description"]

		if !codeExists || !nameExists || !colorExists || !permCodesExists || !descExists {
			log.Printf("Skipping record due to missing required columns: %v", record)
			continue
		}

		log.Printf("Seeding role: %s - %s", record[codeIdx], record[nameIdx])

		roleCreate := client.Role.Create().
			SetCode(record[codeIdx]).
			SetName(record[nameIdx]).
			SetColor(record[colorIdx]).
			SetDescription(record[descIdx])

		roleUpsert := roleCreate.OnConflict(sql.ConflictColumns("code")).
			UpdateNewValues()

		if err := roleUpsert.Exec(ctx); err != nil {
			log.Printf("Failed to seed role %s: %v", record[codeIdx], err)
			continue
		}

		// Handle perm_codes
		// Split elements by comma
		permCodes := record[permCodesIdx]
		if permCodes != "" {
			permCodesList := strings.Split(permCodes, ",")
			for _, permCode := range permCodesList {
				// Find the permission by code
				perm, err := client.Perm.Query().
					Where(perm.Code(permCode)).
					Only(ctx)
				if err != nil {
					log.Printf("Failed to find permission %s for role %s: %v", permCode, record[codeIdx], err)
					continue
				}

				// Add the permission to the role
				if err := roleCreate.AddPerms(perm).Exec(ctx); err != nil {
					log.Printf("Failed to add permission %s to role %s: %v", permCode, record[codeIdx], err)
					continue
				}
			}
		}
	}

	return nil
}

func SeedAdminRole(ctx context.Context, client *ent.Client) error {
	log.Println("Seeding admin role")

	// Create the admin role
	adm := client.Role.Create().
		SetCode("admin").
		SetName("Quản trị viên").
		SetColor("#FF0000").
		SetDescription("tài khoản quản trị viên").
		OnConflict(sql.ConflictColumns("code")).
		Ignore()
	if err := adm.Exec(ctx); err != nil {
		return utils.WrapError("creating admin role", err)
	}

	// Fetch all perms from the database
	perms, err := client.Perm.Query().All(ctx)
	if err != nil {
		return utils.WrapError("fetching permissions", err)
	}

	admID, err := adm.ID(ctx)
	if err != nil {
		return utils.WrapError("getting admin role ID", err)
	}

	// Edge the admin role with all permissions
	for _, perm := range perms {
		if err := client.Role.UpdateOneID(admID).
			AddPerms(perm).
			Exec(ctx); err != nil {
			return utils.WrapError("adding permission to admin role", err)
		}
	}

	log.Println("Admin role seeded successfully")
	return nil
}
