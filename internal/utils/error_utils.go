package utils

import (
	"github.com/gin-gonic/gin"
)

// RespondWithError sends an error response in JSON format.
func RespondWithError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{"error": err.Error()})
}
