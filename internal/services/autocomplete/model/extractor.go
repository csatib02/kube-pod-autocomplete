package model

// FieldExtractor interface defines the method for extracting field values from a PodList
// NOTE: There is no actual difference between ListExtractor and MapExtractor,
// since when processing the extracted data, we can always check the type of the underlying data structure
// via FieldFilter.FieldType, but for the sake of clarity, I have defined two separate types.
type FieldExtractor interface {
	Extract(Resource) any
}

type ListExtractor func(resource Resource) interface{}

func (e ListExtractor) Extract(resource Resource) interface{} {
	return e(resource)
}

type MapExtractor func(resource Resource) interface{}

func (e MapExtractor) Extract(resource Resource) interface{} {
	return e(resource)
}
