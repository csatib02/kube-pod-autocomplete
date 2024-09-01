package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	services "github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete"
	"github.com/csatib02/kube-pod-autocomplete/pkg/utils"
)

// AutocompleteHandler handles autocomplete requests for Pod resources.
func AutocompleteHandler(c *gin.Context) {
	// TODO: Add support for requestedFilters params
	requestedFilters := validateRequestedFilters(c.Params.ByName("requestedFilters"))

	// Get the AutocompleteService instance
	autocompleteService, err := services.NewAutoCompleteService()
	if err != nil {
		utils.HandleHTTPError(c, err)
		return
	}

	suggestions, err := autocompleteService.GetAutocompleteSuggestions(c, requestedFilters)
	if err != nil {
		utils.HandleHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, suggestions)
}

// validateRequestedFilters validates the requestedFilters parameter.
func validateRequestedFilters(requestedFilters string) []string {
	if requestedFilters == "" {
		return []string{}
	}

	// Filters expected to be formatted as e.g.: namespace,pod,labels,annotations
	requestedFiltersAsSlice := strings.Split(requestedFilters, ",")
	if len(requestedFiltersAsSlice) == 1 && requestedFiltersAsSlice[0] == "" {
		return []string{}
	}

	return requestedFiltersAsSlice
}
