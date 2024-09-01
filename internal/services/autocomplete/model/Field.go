package model

type FieldType int

const (
	ListFilter FieldType = iota
	MapFilter
)

type FieldFilter struct {
	ResourceType ResourceType
	Type         FieldType
	Extractor    FieldExtractor
}
