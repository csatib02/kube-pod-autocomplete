package model

import "github.com/csatib02/kube-pod-autocomplete/pkg/common"

type FieldType int

const (
	ListFilter FieldType = iota
	MapFilter
)

type FieldFilter struct {
	ResourceType common.ResourceType
	Type         FieldType
	Extractor    FieldExtractor
}
