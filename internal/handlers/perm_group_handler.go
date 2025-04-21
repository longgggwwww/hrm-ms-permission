package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/longgggwww/hrm-ms-permission/ent"
)

func GetPermGroupsHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		permGroups, err := client.PermGroup.Query().All(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, permGroups)
	}
}
