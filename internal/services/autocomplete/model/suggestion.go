package model

// Suggestion represents a suggestion for autocompletion
type Suggestion struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}
