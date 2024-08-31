package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	services "github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete"
	"github.com/csatib02/kube-pod-autocomplete/pkg/utils"
)

// AutocompleteHandler handles autocomplete requests for Pod resources.
func AutocompleteHandler(c *gin.Context) {
	// Get the query from the request
	query := c.Query("q")
	// TODO: Enable query params with validation
	// var req QueryRequest
	// if err := c.BindJSON(&req); err != nil {
	// 	utils.HandleHTTPError(c, err)
	// 	return
	// }

	// Get the AutocompleteService instance
	autocompleteService, err := services.NewAutoCompleteService()
	if err != nil {
		utils.HandleHTTPError(c, err)
		return
	}

	suggestions, err := autocompleteService.GetAutocompleteSuggestions(c, query)
	if err != nil {
		utils.HandleHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}
