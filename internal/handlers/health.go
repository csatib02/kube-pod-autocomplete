package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles the requests for health checks
func HealthHandler(c *gin.Context) {
	slog.Info("Received health check request")

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
