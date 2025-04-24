package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/longgggwwww/hrm-ms-permission/ent"
)

func GetRolesHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := client.Role.Query().WithPerms().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, roles)
	}
}
