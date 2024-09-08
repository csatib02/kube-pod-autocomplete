package autocomplete

import (
	"context"
	"errors"
	"fmt"

	"github.com/csatib02/kube-pod-autocomplete/internal/k8s"
	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/filter"
	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/model"
	"github.com/csatib02/kube-pod-autocomplete/pkg/common"
)

type Service struct {
	k8sClient k8s.Client
}

func NewAutoCompleteService() (*Service, error) {
	client, err := k8s.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %w", err)
	}

	return &Service{
		k8sClient: *client,
	}, nil
}

// GetAutocompleteSuggestions returns a list of suggestions (for the given query)
func (s *Service) GetAutocompleteSuggestions(ctx context.Context, req model.AutoCompleteRequest) (*model.AutocompleteSuggestions, error) {
	filters, err := filter.NewFieldFilters(req.ResourceType, &req.Filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get field filters: %w", err)
	}

	resources, err := s.k8sClient.ListResource(ctx, req.ResourceType)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	return s.extractSuggestions(resources, filters)
}

// extractSuggestions extracts suggestions from the given resources based on the requested filters
func (s *Service) extractSuggestions(resources common.Resources, filters *map[string]model.FieldFilter) (*model.AutocompleteSuggestions, error) {
	suggestions := make([]model.Suggestion, 0, len(*filters))
	for fieldName, fieldFilter := range *filters {
		extractedData := fieldFilter.Extractor.Extract(resources)
		switch fieldFilter.Type {
		case model.ListFilter:
			listData, ok := extractedData.([]string)
			if !ok {
				return nil, errors.New("invalid data type for ListFilter")
			}

			suggestions = append(suggestions, model.Suggestion{
				Key:    fieldName,
				Values: listData,
			})

		case model.MapFilter:
			mapData, ok := extractedData.(map[string][]string)
			if !ok {
				return nil, errors.New("invalid data type for MapFilter")
			}

			s.processMapSuggestion(&suggestions, fieldName, &mapData)

		default:
			return nil, fmt.Errorf("unsupported filter type: %v", fieldFilter.Type)
		}
	}

	// These should be options on the UI
	filterOptions := filter.Options{}

	filterOptions.RemoveDuplicateValues(&suggestions)
	filterOptions.RemoveIgnoredKeys(&suggestions)

	return &model.AutocompleteSuggestions{Suggestions: suggestions}, nil
}

func (s *Service) processMapSuggestion(suggestions *[]model.Suggestion, filterName string, mapData *map[string][]string) {
	for key, value := range *mapData {
		*suggestions = append(*suggestions, model.Suggestion{
			Key:    fmt.Sprintf("%s:%s", filterName, key),
			Values: value,
		})
	}
}
