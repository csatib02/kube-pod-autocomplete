package handlers

import "github.com/gin-gonic/gin"

// SetupRoutes registers all the handlers on the router
func SetupRoutes(router *gin.Engine) {
	router.GET("/search/autocomplete/:resource", AutocompleteHandler)
	router.GET("/health", HealthHandler)
}
