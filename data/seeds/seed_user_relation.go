package seeds

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect/sql"
	"github.com/longgggwwww/hrm-ms-permission/ent"
	"github.com/longgggwwww/hrm-ms-permission/ent/perm"
	"github.com/longgggwwww/hrm-ms-permission/ent/role"
	"github.com/longgggwwww/hrm-ms-permission/internal/utils"
)

// SeedUserRoles đọc file user_role.csv và seed dữ liệu vào bảng user_role
func SeedUserRoles(ctx context.Context, client *ent.Client) error {
	filePath := filepath.Join("data", "user_role.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return utils.WrapError("opening file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
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
		userIDIdx, userIDExists := headerMap["user_id"]
		roleCodeIdx, roleCodeExists := headerMap["role_code"]
		if !userIDExists || !roleCodeExists {
			log.Printf("Skipping record due to missing required columns: %v", record)
			continue
		}

		userID := record[userIDIdx]
		roleCode := record[roleCodeIdx]

		role, err := client.Role.Query().Where(role.Code(roleCode)).Only(ctx)
		if err != nil {
			log.Printf("Failed to find role %s for user %s: %v", roleCode, userID, err)
			continue
		}

		err = client.UserRole.Create().
			SetUserID(userID).
			SetRoleID(role.ID).
			OnConflict(sql.ConflictColumns("role_id", "user_id")).
			UpdateNewValues().
			Exec(ctx)
		if err != nil {
			log.Printf("Failed to seed user_role for user %s and role %s: %v", userID, roleCode, err)
			continue
		}
		log.Printf("Seeded user_role: user_id=%s, role_code=%s", userID, roleCode)
	}
	return nil
}

// SeedUserPerms đọc file user_perm.csv và seed dữ liệu vào bảng user_perm
func SeedUserPerms(ctx context.Context, client *ent.Client) error {
	filePath := filepath.Join("data", "user_perm.csv")
	file, err := os.Open(filePath)
	if err != nil {
		return utils.WrapError("opening file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
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
		userIDIdx, userIDExists := headerMap["user_id"]
		permCodeIdx, permCodeExists := headerMap["perm_code"]
		if !userIDExists || !permCodeExists {
			log.Printf("Skipping record due to missing required columns: %v", record)
			continue
		}

		userID := record[userIDIdx]
		permCode := record[permCodeIdx]

		perm, err := client.Perm.Query().Where(perm.Code(permCode)).Only(ctx)
		if err != nil {
			log.Printf("Failed to find perm %s for user %s: %v", permCode, userID, err)
			continue
		}

		err = client.UserPerm.Create().
			SetUserID(userID).
			SetPermID(perm.ID).
			OnConflict(sql.ConflictColumns("perm_id", "user_id")).
			UpdateNewValues().
			Exec(ctx)
		if err != nil {
			log.Printf("Failed to seed user_perm for user %s and perm %s: %v", userID, permCode, err)
			continue
		}
		log.Printf("Seeded user_perm: user_id=%s, perm_code=%s", userID, permCode)
	}
	return nil
}
