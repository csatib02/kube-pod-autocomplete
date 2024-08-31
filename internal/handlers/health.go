package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles the requests for health checks.
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
