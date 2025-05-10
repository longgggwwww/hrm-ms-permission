package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-permission/ent"
	"github.com/longgggwwww/hrm-ms-permission/internal/utils"
)

type PermHandler struct {
	Client *ent.Client
}

// GetPerms handles the retrieval of permissions.
func (h *PermHandler) GetPerms(c *gin.Context) {
	perms, err := h.Client.Perm.Query().WithGroup().All(c.Request.Context())
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, perms)
}

// RegisterRoutes registers the permission-related routes.
func (h *PermHandler) RegisterRoutes(r *gin.Engine) {
	gr := r.Group("/perms")
	{
		gr.GET("", h.GetPerms)
	}
}

// respondWithError sends an error response in JSON format.
func (h *PermHandler) respondWithError(c *gin.Context, statusCode int, err error) {
	utils.RespondWithError(c, statusCode, err)
}
