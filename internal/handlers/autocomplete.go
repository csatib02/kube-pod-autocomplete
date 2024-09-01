package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	services "github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete"
	"github.com/csatib02/kube-pod-autocomplete/pkg/utils"
)

// AutocompleteHandler handles autocomplete requests for Pod resources
func AutocompleteHandler(c *gin.Context) {
	slog.Info("Received autocomplete request")

	// TODO: Add support for requestedFilters params
	requestedFilters := validateRequestedFilters(c.Params.ByName("requestedFilters"))

	autocompleteService, err := services.NewAutoCompleteService()
	if err != nil {
		slog.Error(fmt.Errorf("failed to create autocomplete service: %w", err).Error())
		utils.HandleHTTPError(c, err)
		return
	}

	suggestions, err := autocompleteService.GetAutocompleteSuggestions(c, requestedFilters)
	if err != nil {
		slog.Error(fmt.Errorf("failed to get autocomplete suggestions: %w", err).Error())

		utils.HandleHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, suggestions)
}

// validateRequestedFilters validates the requestedFilters parameter
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
