package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-permission/ent"
	"github.com/longgggwwww/hrm-ms-permission/internal/utils"
)

type PermGroupHandler struct {
	Client *ent.Client
}

func (h *PermGroupHandler) GetPermGroups(c *gin.Context) {
	permGroups, err := h.Client.PermGroup.Query().WithPerms().All(c.Request.Context())
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, permGroups)
}

func (h *PermGroupHandler) RegisterRoutes(r *gin.Engine) {
	gr := r.Group("/perm-groups")
	{
		gr.GET("", h.GetPermGroups)
	}
}

func (h *PermGroupHandler) respondWithError(c *gin.Context, statusCode int, err error) {
	utils.RespondWithError(c, statusCode, err)
}
