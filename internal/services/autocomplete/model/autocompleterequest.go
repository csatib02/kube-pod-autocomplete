package model

import "github.com/csatib02/kube-pod-autocomplete/pkg/common"

type AutoCompleteRequest struct {
	ResourceType common.ResourceType `json:"resourceType"`
	Filters      []string            `json:"filters"`
	Query        string              `json:"query"` // Currently unused
}
