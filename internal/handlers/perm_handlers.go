package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-permission/ent"
)

type PermHandler struct {
	Client *ent.Client
}

func (h *PermHandler) GetPerms(c *gin.Context) {
	perms, err := h.Client.Perm.Query().WithGroup().All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, perms)
}

func (h *PermHandler) CreatePerm(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	perm, err := h.Client.Perm.Create().
		SetName(input.Name).
		Save(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, perm)
}

func (h *PermHandler) UpdatePerm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}
	var input struct {
		Name *string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := h.Client.Perm.UpdateOneID(id)
	if input.Name != nil {
		update.SetName(*input.Name)
	}

	perm, err := update.Save(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, perm)
}

func (h *PermHandler) DeletePerm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}
	if err := h.Client.Perm.DeleteOneID(id).Exec(context.Background()); err != nil {
		if ent.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *PermHandler) RegisterRoutes(r *gin.Engine) {
	gr := r.Group("/perms")
	{
		gr.GET("", h.GetPerms)
		gr.POST("", h.CreatePerm)
		gr.PUT(":id", h.UpdatePerm)
		gr.DELETE(":id", h.DeletePerm)
	}
}
