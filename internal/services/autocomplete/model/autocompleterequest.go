package model

type AutoCompleteRequest struct {
	ResourceType ResourceType `json:"resourceType"`
	Filters      []string     `json:"filters"`
	Query        string       `json:"query"` // Currently not used
}
