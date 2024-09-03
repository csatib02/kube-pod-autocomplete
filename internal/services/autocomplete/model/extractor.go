package model

import "github.com/csatib02/kube-pod-autocomplete/pkg/common"

// FieldExtractor interface defines the method for extracting field values from resources
type FieldExtractor interface {
	Extract(common.Resources) any
}

type Extractor func(resource common.Resources) any

func (e Extractor) Extract(resource common.Resources) any {
	return e(resource)
}
