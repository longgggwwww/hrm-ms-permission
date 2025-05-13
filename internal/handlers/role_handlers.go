package handlers

import (
	"context"
	"fmt"
	"net/http"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pb "github.com/huynhthanhthao/hrm_user_service/generated"
	"github.com/longgggwwww/hrm-ms-permission/ent"
	"github.com/longgggwwww/hrm-ms-permission/ent/role"
	"github.com/longgggwwww/hrm-ms-permission/ent/userrole"
	"github.com/longgggwwww/hrm-ms-permission/internal/utils"
)

type RoleHandler struct {
	Client     *ent.Client
	UserClient *pb.UserServiceClient
}

func NewRoleHandler(client *ent.Client, userClient *pb.UserServiceClient) *RoleHandler {
	return &RoleHandler{
		Client:     client,
		UserClient: userClient,
	}
}

func (h *RoleHandler) GetRoles(c *gin.Context) {
	roles, err := h.Client.Role.Query().WithPerms().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var input struct {
		Code        string  `json:"code" binding:"required"`
		Name        string  `json:"name" binding:"required"`
		Color       *string `json:"color"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleCreate := h.Client.Role.Create().
		SetCode(input.Code).
		SetName(input.Name)
	if input.Color != nil {
		roleCreate.SetColor(*input.Color)
	}
	if input.Description != nil {
		roleCreate.SetDescription(*input.Description)
	}

	role, err := roleCreate.Save(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	// Parse the UUID from the URL parameter
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var input struct {
		Name        *string      `json:"name"`
		Color       *string      `json:"color"`
		Description *string      `json:"description"`
		PermIDs     []*uuid.UUID `json:"perm_ids"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := h.Client.Role.UpdateOneID(id)
	if input.Name != nil {
		update.SetName(*input.Name)
	}
	if input.Color != nil {
		update.SetColor(*input.Color)
	}
	if input.Description != nil {
		update.SetDescription(*input.Description)
	}
	if len(input.PermIDs) > 0 {
		update.ClearPerms()
		for _, permID := range input.PermIDs {
			update.AddPermIDs(*permID)
		}
	}

	updated, err := update.Save(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedRole, err := h.Client.Role.Query().Where(role.IDEQ(updated.ID)).WithPerms().Only(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated role with permissions"})
		return
	}

	c.JSON(http.StatusOK, updatedRole)
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	// Parse the UUID from the URL parameter
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	if err := h.Client.Role.DeleteOneID(id).Exec(context.Background()); err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *RoleHandler) DeleteMultipleRoles(c *gin.Context) {
	var input struct {
		IDs []uuid.UUID `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	if len(input.IDs) == 0 {
		h.handleError(c, http.StatusBadRequest, "No role IDs provided")
		return
	}

	tx, err := h.Client.Tx(context.Background())
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	for _, id := range input.IDs {
		if err := tx.Role.DeleteOneID(id).Exec(context.Background()); err != nil {
			tx.Rollback()
			if ent.IsNotFound(err) {
				h.handleError(c, http.StatusNotFound, fmt.Sprintf("Role with ID %s not found", id))
			} else {
				h.handleError(c, http.StatusInternalServerError, err.Error())
			}
			return
		}
	}

	if err := tx.Commit(); err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *RoleHandler) checkUserIDsExist(userIDs []string, users *pb.GetUsersByIDsResponse) error {
	userIDsMap := make(map[string]bool)
	for _, user := range users.Users {
		userIDsMap[user.Id] = true
	}
	for _, userID := range userIDs {
		if _, exists := userIDsMap[userID]; !exists {
			return fmt.Errorf("user ID %s does not exist", userID)
		}
	}
	return nil
}

func (h *RoleHandler) assignRoleToUsers(roleID uuid.UUID, userIDs []string) error {
	for _, userID := range userIDs {
		err := h.Client.UserRole.Create().
			SetRoleID(roleID).
			SetUserID(userID).
			OnConflict(sql.ConflictColumns("role_id", "user_id")).
			UpdateNewValues().
			Exec(context.Background())

		if err != nil {
			return fmt.Errorf("failed to assign role to user %s: %v", userID, err)
		}
	}
	return nil
}

func (h *RoleHandler) handleError(c *gin.Context, statusCode int, message string) {
	utils.RespondWithError(c, statusCode, fmt.Errorf("%s", message))
}

func (h *RoleHandler) AssignRoleToUsers(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid UUID format for role ID")
		return
	}

	exists, err := h.Client.Role.Query().Where(role.IDEQ(roleID)).Exist(context.Background())
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to check role existence")
		return
	}
	if !exists {
		h.handleError(c, http.StatusNotFound, "Role not found")
		return
	}

	var input struct {
		UserIDs []string `json:"user_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	users, err := (*h.UserClient).GetUsersByIDs(context.Background(), &pb.GetUsersByIDsRequest{
		Ids: input.UserIDs,
	})
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	if err := h.checkUserIDsExist(input.UserIDs, users); err != nil {
		h.handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.assignRoleToUsers(roleID, input.UserIDs); err != nil {
		h.handleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned to users successfully"})
}

func (h *RoleHandler) GetUsersByRole(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid UUID format for role ID")
		return
	}

	userRoles, err := h.Client.UserRole.Query().Where(userrole.RoleIDEQ(roleID)).All(context.Background())
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to fetch users by role")
		return
	}
	if len(userRoles) == 0 {
		h.handleError(c, http.StatusNotFound, "No users found for the given role")
		return
	}

	// Get user IDs from userRoles
	var userIDs []string
	for _, userRole := range userRoles {
		userIDs = append(userIDs, userRole.UserID)
	}

	// Fetch user details from the UserService
	respb, err := (*h.UserClient).GetUsersByIDs(context.Background(), &pb.GetUsersByIDsRequest{
		Ids: userIDs,
	})
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	if len(respb.GetUsers()) == 0 {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	c.JSON(http.StatusOK, respb.GetUsers())
}

func (h *RoleHandler) GetRolesByUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		h.handleError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	userRoles, err := h.Client.UserRole.Query().Where(userrole.UserIDEQ(userID)).WithRole().All(context.Background())
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, "Failed to fetch roles for the user")
		return
	}

	fmt.Println("User Roles:", userRoles)

	var roles []*ent.Role = []*ent.Role{}
	for _, userRole := range userRoles {
		role := userRole.Edges.Role
		if role == nil {
			h.handleError(c, http.StatusInternalServerError, "Role not found for user role")
			return
		}
		roles = append(roles, role)
	}

	c.JSON(http.StatusOK, roles)
}

func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid UUID format for role ID")
		return
	}

	role, err := h.Client.Role.Query().Where(role.IDEQ(id)).WithPerms().Only(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			h.handleError(c, http.StatusNotFound, "Role not found")
		} else {
			h.handleError(c, http.StatusInternalServerError, "Failed to fetch role")
		}
		return
	}

	c.JSON(http.StatusOK, role)
}

func (h *RoleHandler) RegisterRoutes(r *gin.Engine) {
	gr := r.Group("/roles")
	{
		gr.GET("", h.GetRoles)
		gr.POST("", h.CreateRole)
		gr.PATCH(":id", h.UpdateRole)
		gr.DELETE(":id", h.DeleteRole)
		gr.DELETE("", h.DeleteMultipleRoles)
		gr.POST(":id/assign", h.AssignRoleToUsers)
		gr.GET(":id/users", h.GetUsersByRole)
		gr.GET(":id", h.GetRoleByID)
	}

	ur := r.Group("/users")
	{
		ur.GET(":user_id/roles", h.GetRolesByUser)
	}
}
