package model

import "github.com/csatib02/kube-pod-autocomplete/pkg/common"

// FieldExtractor interface defines the method for extracting field values from a PodList
// NOTE: There is no actual difference between ListExtractor and MapExtractor,
// since when processing the extracted data, we can always check the type of the underlying data structure
// via FieldFilter.Type, but for the sake of clarity, I have defined two separate types.
type FieldExtractor interface {
	Extract(common.Resources) any
}

type Extractor func(resource common.Resources) any

func (e Extractor) Extract(resource common.Resources) any {
	return e(resource)
}
