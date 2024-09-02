package podfilter

import (
	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/model"
	"github.com/csatib02/kube-pod-autocomplete/pkg/common"
)

var supportedFilters = map[string]model.FieldFilter{
	"namespace": {
		Type: model.ListFilter,
		Extractor: model.Extractor(func(resource common.Resources) interface{} {
			podResource := resource.(common.PodResources)
			result := make([]string, 0, len(podResource.Items))
			for _, pod := range podResource.Items {
				result = append(result, pod.Namespace)
			}
			return result
		}),
	},
	"phase": {
		Type: model.ListFilter,
		Extractor: model.Extractor(func(resource common.Resources) interface{} {
			podResource := resource.(common.PodResources)
			result := make([]string, 0, len(podResource.Items))
			for _, pod := range podResource.Items {
				result = append(result, string(pod.Status.Phase))
			}
			return result
		}),
	},
	"labels": {
		Type: model.MapFilter,
		Extractor: model.Extractor(func(resource common.Resources) interface{} {
			podResource := resource.(common.PodResources)
			result := make(map[string][]string)
			for _, pod := range podResource.Items {
				for key, value := range pod.Labels {
					result[key] = append(result[key], value)
				}
			}
			return result
		}),
	},
	"annotations": {
		Type: model.MapFilter,
		Extractor: model.Extractor(func(resource common.Resources) interface{} {
			podResource := resource.(common.PodResources)
			result := make(map[string][]string)
			for _, pod := range podResource.Items {
				for key, value := range pod.Annotations {
					result[key] = append(result[key], value)
				}
			}
			return result
		}),
	},
}

// GetFilters returns the supported filters based on the requested filters
// if called with empty requestedFilters or nil, it returns all supported filters
func GetFilters(requestedFilters *[]string) *map[string]model.FieldFilter {
	if requestedFilters == nil || len(*requestedFilters) == 0 {
		return &supportedFilters
	}

	filters := make(map[string]model.FieldFilter)
	for _, requestedFilter := range *requestedFilters {
		if fieldFilter, ok := supportedFilters[requestedFilter]; ok {
			filters[requestedFilter] = fieldFilter
		}
	}

	return &filters
}
