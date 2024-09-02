package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	services "github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete"
	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/model"
	"github.com/csatib02/kube-pod-autocomplete/pkg/utils"
)

// AutocompleteHandler handles autocomplete requests for Pod resources
func AutocompleteHandler(c *gin.Context) {
	slog.Info("Received autocomplete request")

	// TODO: Add support for query parameters
	// var req model.AutoCompleteRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	slog.Error(fmt.Errorf("failed to bind request: %w", err).Error())
	// 	utils.HandleHTTPError(c, errors.New("failed to bind request"))
	// 	return
	// }

	resourceParam := c.Param("resource")
	if resourceParam == "" {
		slog.Error("resource parameter is missing")
		utils.HandleHTTPError(c, http.StatusBadRequest, errors.New("resource parameter is missing"))
		return
	}

	resourceType := model.ResourceType(resourceParam)
	if !model.IsValidResourceType(resourceType) {
		slog.Error(fmt.Sprintf("resource type: %s not supported", resourceType))
		utils.HandleHTTPError(c, http.StatusBadRequest, errors.New("invalid resource type"))
		return
	}

	// For now, hardcoding the request
	req := model.AutoCompleteRequest{
		ResourceType: resourceType,
		Filters:      []string{"namespace", "phase", "labels", "annotations"},
	}

	validFilters, err := validateRequestedFilters(req.Filters)
	if err != nil {
		slog.Error(fmt.Errorf("failed to validate requested filters: %w", err).Error())
		utils.HandleHTTPError(c, http.StatusBadRequest, err)
		return
	}
	req.Filters = validFilters

	autocompleteService, err := services.NewAutoCompleteService()
	if err != nil {
		slog.Error(fmt.Errorf("failed to create autocomplete service: %w", err).Error())
		utils.HandleHTTPError(c, http.StatusBadRequest, err)
		return
	}

	suggestions, err := autocompleteService.GetAutocompleteSuggestions(c, req)
	if err != nil {
		slog.Error(fmt.Errorf("failed to get autocomplete suggestions: %w", err).Error())
		utils.HandleHTTPError(c, http.StatusBadRequest, err)
		return
	}

	// Pretty-print the JSON response
	prettyJSON, err := json.MarshalIndent(suggestions, "", "  ")
	if err != nil {
		// Log the error and return the response as is
		slog.Error(fmt.Errorf("failed to pretty-print JSON response: %w", err).Error())
		c.JSON(http.StatusOK, suggestions)
		return
	}

	c.Data(http.StatusOK, "application/json", prettyJSON)
}

// validateRequestedFilters validates the requestedFilters parameter
func validateRequestedFilters(requestedFilters []string) ([]string, error) {
	validFilters := make([]string, 0, len(requestedFilters))
	for _, filter := range requestedFilters {
		if strings.TrimSpace(filter) == "" {
			continue
		}

		validFilters = append(validFilters, filter)
	}

	if len(validFilters) == 0 {
		return nil, errors.New("no valid filters provided")
	}

	return validFilters, nil
}
