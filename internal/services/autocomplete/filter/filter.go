package filter

import (
	"errors"

	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/filter/podfilter"
	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/model"
	"github.com/csatib02/kube-pod-autocomplete/pkg/common"
)

// NewFieldFilters returns the supported filters for the requested resource type
func NewFieldFilters(resourceType common.ResourceType, requestedFilters *[]string) (*map[string]model.FieldFilter, error) {
	switch resourceType {
	case common.PodResourceType:
		return podfilter.GetFilters(requestedFilters), nil
	default:
		return nil, errors.New("unsupported resource type")
	}
}
