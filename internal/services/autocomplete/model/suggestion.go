package model

type AutocompleteSuggestions struct {
	Suggestions []Suggestion `json:"suggestions"`
}

type Suggestion struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}
