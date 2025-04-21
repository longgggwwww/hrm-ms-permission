package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/longgggwww/hrm-ms-permission/ent"
)

func GetPermsHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		perms, err := client.Perm.Query().WithGroup().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, perms)
	}
}
