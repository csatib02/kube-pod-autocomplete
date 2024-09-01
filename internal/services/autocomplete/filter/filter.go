package filter

import (
	"errors"

	v1 "k8s.io/api/core/v1"

	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/model"
)

type FieldType int

const (
	ListFilter FieldType = iota
	MapFilter
)

type FieldFilter struct {
	Type      FieldType
	FieldName string
	Extractor model.FieldExtractor
}

var supportedFilters = map[string]FieldFilter{
	"namespace": {
		Type:      ListFilter,
		FieldName: "namespace",
		Extractor: model.ListExtractor(func(pods *v1.PodList) interface{} {
			result := make([]string, 0, len(pods.Items))
			for _, pod := range pods.Items {
				result = append(result, pod.Namespace)
			}
			return result
		}),
	},
	"phase": {
		Type:      ListFilter,
		FieldName: "phase",
		Extractor: model.ListExtractor(func(pods *v1.PodList) interface{} {
			result := make([]string, 0, len(pods.Items))
			for _, pod := range pods.Items {
				result = append(result, string(pod.Status.Phase))
			}
			return result
		}),
	},
	"labels": {
		Type:      MapFilter,
		FieldName: "labels",
		Extractor: model.MapExtractor(func(pods *v1.PodList) interface{} {
			result := make(map[string][]string)
			for _, pod := range pods.Items {
				for key, value := range pod.Labels {
					result[key] = append(result[key], value)
				}
			}
			return result
		}),
	},
	"annotations": {
		Type:      MapFilter,
		FieldName: "annotations",
		Extractor: model.MapExtractor(func(pods *v1.PodList) interface{} {
			result := make(map[string][]string)
			for _, pod := range pods.Items {
				for key, value := range pod.Annotations {
					result[key] = append(result[key], value)
				}
			}
			return result
		}),
	},
}

func GetSupportedFilters() []string {
	var filters []string
	for filter := range supportedFilters {
		filters = append(filters, filter)
	}

	return filters
}

func ParseFilters(filter []string) (*[]FieldFilter, error) {
	var filters []FieldFilter
	for _, f := range filter {
		if fieldFilter, ok := supportedFilters[f]; ok {
			filters = append(filters, fieldFilter)
		}
	}

	if len(filters) == 0 {
		return nil, errors.New("no supported filters found")
	}

	return &filters, nil
}
