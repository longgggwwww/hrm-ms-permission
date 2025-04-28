package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-permission/ent"
)

type PermGroupHandler struct {
	Client *ent.Client
}

func (h *PermGroupHandler) GetPermGroups(c *gin.Context) {
	permGroups, err := h.Client.PermGroup.Query().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, permGroups)
}

func (h *PermGroupHandler) CreatePermGroup(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permGroup, err := h.Client.PermGroup.Create().
		SetName(input.Name).
		Save(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, permGroup)
}

func (h *PermGroupHandler) UpdatePermGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission group ID"})
		return
	}
	var input struct {
		Name *string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := h.Client.PermGroup.UpdateOneID(id)
	if input.Name != nil {
		update.SetName(*input.Name)
	}

	permGroup, err := update.Save(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permGroup)
}

func (h *PermGroupHandler) DeletePermGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission group ID"})
		return
	}
	if err := h.Client.PermGroup.DeleteOneID(id).Exec(context.Background()); err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permission group not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *PermGroupHandler) RegisterRoutes(r *gin.Engine) {
	gr := r.Group("/perm-groups")
	{
		gr.GET("", h.GetPermGroups)
		gr.POST("", h.CreatePermGroup)
		gr.PUT(":id", h.UpdatePermGroup)
		gr.DELETE(":id", h.DeletePermGroup)
	}
}
