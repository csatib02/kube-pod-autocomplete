package handlers

import "github.com/gin-gonic/gin"

// SetupRouter initializes the router and registers all the handlers
func SetupRouter(router *gin.Engine) {
	router.GET("/search/autocomplete/:resource", AutocompleteHandler)
	router.GET("/health", HealthHandler)
}
