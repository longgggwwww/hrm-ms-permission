package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-permission/ent"
)

type PermGroupHandler struct {
	Client *ent.Client
}

func NewPermGroupHandler(c *ent.Client) *PermGroupHandler {
	return &PermGroupHandler{
		Client: c,
	}
}

func (h *PermGroupHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/perm-groups", h.List)
}

func (h *PermGroupHandler) List(c *gin.Context) {
	out, err := h.Client.PermGroup.Query().
		WithPerms().
		All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}
