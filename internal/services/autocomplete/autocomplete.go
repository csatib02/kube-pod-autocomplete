package autocomplete

import (
	"context"
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"

	"github.com/csatib02/kube-pod-autocomplete/internal/k8s"
	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/filter"
	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/model"
)

type Service struct {
	k8sClient k8s.Client
	// TODO: Enable caching
	// useCache          bool
	// cacheUpdatePeriod time.Duration
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

// GetAutocompleteSuggestions returns a list of suggestions for the given query
func (s *Service) GetAutocompleteSuggestions(ctx context.Context, requestedFilters []string) (*model.AutocompleteSuggestions, error) {
	// If no filters are requested, use all supported filters
	if len(requestedFilters) == 0 {
		requestedFilters = filter.GetSupportedFilters()
	}

	pods, err := s.k8sClient.ListPods(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	return s.extractSuggestions(pods, requestedFilters)
}

// extractSuggestions extracts suggestions from the given pods based on the requested filters
func (s *Service) extractSuggestions(pods *v1.PodList, requestedFilters []string) (*model.AutocompleteSuggestions, error) {
	filterInfos, err := filter.ParseFilters(requestedFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to parse filters: %w", err)
	}

	suggestions := make([]model.Suggestion, 0, len(*filterInfos))
	for _, filterInfo := range *filterInfos {
		extractedData := filterInfo.Extractor.Extract(pods)

		switch filterInfo.Type {
		case filter.ListFilter:
			listData, ok := extractedData.([]string)
			if !ok {
				return nil, errors.New("invalid data type for ListFilter")
			}

			suggestions = append(suggestions, model.Suggestion{
				Key:    filterInfo.ForField,
				Values: listData,
			})

		case filter.MapFilter:
			mapData, ok := extractedData.(map[string][]string)
			if !ok {
				return nil, errors.New("invalid data type for MapFilter")
			}

			s.processMapSuggestion(&suggestions, filterInfo.ForField, &mapData)

		default:
			return nil, fmt.Errorf("unsupported filter type: %v", filterInfo.Type)
		}
	}

	return &model.AutocompleteSuggestions{Suggestions: suggestions}, nil
}

func (s *Service) processMapSuggestion(suggestions *[]model.Suggestion, filter string, mapData *map[string][]string) {
	for key, value := range *mapData {
		*suggestions = append(*suggestions, model.Suggestion{
			Key:    fmt.Sprintf("%s:%s", filter, key),
			Values: value,
		})
	}
}
