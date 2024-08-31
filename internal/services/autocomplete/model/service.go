package model

// AutocompleteService provides methods to get autocomplete suggestions
type AutocompleteService interface {
	GetAutocompleteSuggestions(query string) ([]Suggestion, error)
}
