package filter

import (
	"strings"

	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/model"
)

// FilterOptions defines additional options for filtering suggestions
type Options struct{}

var ignoredKeys = map[string][]string{
	"annotations": {"annotations:kubectl.kubernetes.io/last-applied-configuration"}, // Contains sensitive information
}

// RemoveDuplicateValues removes duplicate values from the suggestions
func (o *Options) RemoveDuplicateValues(suggestions *[]model.Suggestion) {
	for i, suggestion := range *suggestions {
		valueMap := make(map[string]bool)
		uniqueValues := []string{}

		for _, value := range suggestion.Values {
			// If the value is not seen before
			if !valueMap[value] {
				valueMap[value] = true
				uniqueValues = append(uniqueValues, value)
			}
		}

		// Replace the original Values slice with the uniqueValues slice
		(*suggestions)[i].Values = uniqueValues
	}
}

// RemoveIgnoredKeys removes the ignored keys from the suggestions
// NOTE: IgnoreKeys should be retrieved from request parameters
func (o *Options) RemoveIgnoredKeys(suggestions *[]model.Suggestion) {
	filteredSuggestions := make([]model.Suggestion, 0, len(*suggestions))
	for _, suggestion := range *suggestions {
		ignored := false
		for _, ignoredKey := range ignoredKeys[strings.Split(suggestion.Key, ":")[0]] {
			if suggestion.Key == ignoredKey {
				ignored = true
				break
			}
		}

		if !ignored {
			filteredSuggestions = append(filteredSuggestions, suggestion)
		}
	}

	*suggestions = filteredSuggestions
}
