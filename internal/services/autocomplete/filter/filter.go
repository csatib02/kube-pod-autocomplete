package filter

import (
	"errors"

	v1 "k8s.io/api/core/v1"

	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/model"
)

type FilterType int

const (
	ListFilter FilterType = iota
	MapFilter
)

type FilterInfo struct {
	Type      FilterType
	ForField  string
	Extractor model.FieldExtractor
}

var supportedFilters = map[string]FilterInfo{
	"namespace": {
		Type:     ListFilter,
		ForField: "namespace",
		Extractor: model.ListExtractor(func(pods *v1.PodList) interface{} {
			result := make([]string, 0, len(pods.Items))
			for _, pod := range pods.Items {
				result = append(result, pod.Namespace)
			}
			return result
		}),
	},
	"phase": {
		Type:     ListFilter,
		ForField: "phase",
		Extractor: model.ListExtractor(func(pods *v1.PodList) interface{} {
			result := make([]string, 0, len(pods.Items))
			for _, pod := range pods.Items {
				result = append(result, string(pod.Status.Phase))
			}
			return result
		}),
	},
	"labels": {
		Type:     MapFilter,
		ForField: "labels",
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
		Type:     MapFilter,
		ForField: "annotations",
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

func ParseFilters(filter []string) (*[]FilterInfo, error) {
	var filters []FilterInfo
	for _, f := range filter {
		if info, ok := supportedFilters[f]; ok {
			filters = append(filters, info)
		}
	}

	if len(filters) == 0 {
		return nil, errors.New("no supported filters found")
	}

	return &filters, nil
}
