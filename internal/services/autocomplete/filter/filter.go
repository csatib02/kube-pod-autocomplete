package filter

import (
	"errors"

	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/filter/podfilter"
	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/model"
)

// NewFieldFilters returns the supported filters for the requested resource type
func NewFieldFilters(resourceType model.ResourceType, requestedFilters *[]string) (*map[string]model.FieldFilter, error) {
	switch resourceType {
	case model.PodResourceType:
		return podfilter.GetFilters(requestedFilters), nil
	// Add cases for other resource types here
	default:
		return nil, errors.New("unsupported resource type")
	}
}
