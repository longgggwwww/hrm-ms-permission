package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pb "github.com/huynhthanhthao/hrm_user_service/generated"
	"github.com/longgggwwww/hrm-ms-permission/ent"
	"github.com/longgggwwww/hrm-ms-permission/ent/role"
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
		Name *string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := h.Client.Role.UpdateOneID(id)
	if input.Name != nil {
		update.SetName(*input.Name)
	}

	role, err := update.Save(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
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

func (h *RoleHandler) AssignRoleToUsers(c *gin.Context) {
	// Parse the UUID from the URL parameter
	roleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for role ID"})
		return
	}

	exists, err := h.Client.Role.Query().Where(role.IDEQ(roleID)).Exist(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check role existence"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	var input struct {
		UserIDs []string `json:"user_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := (*h.UserClient).ListUsers(context.Background(), &pb.ListUsersRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	fmt.Println("Fetched users:", users)

	// Check if user IDs exist
	userIDsMap := make(map[string]bool)
	for _, user := range users.Users {
		userIDsMap[user.Id] = true
	}
	for _, userID := range input.UserIDs {
		if _, exists := userIDsMap[userID]; !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("User ID %s does not exist", userID)})
			return
		}
	}

	// Print userIDsMap for debugging
	fmt.Println("User IDs map:", userIDsMap)

	// // Assign role to users
	// for _, userID := range input.UserIDs {
	// 	_, err := h.Client.UserRole.Create().
	// 		SetRoleID(roleID).
	// 		SetUserID(userID).Save(context.Background())

	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to assign role to user %s: %v", userID, err)})
	// 		return
	// 	}
	// }

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned to users successfully"})
}

func (h *RoleHandler) RegisterRoutes(r *gin.Engine) {
	gr := r.Group("/roles")
	{
		gr.GET("", h.GetRoles)
		gr.POST("", h.CreateRole)
		gr.PUT(":id", h.UpdateRole)
		gr.DELETE(":id", h.DeleteRole)
		gr.POST(":id/assign", h.AssignRoleToUsers)
	}
}
