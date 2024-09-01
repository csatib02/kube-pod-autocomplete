package model

import (
	v1 "k8s.io/api/core/v1"
)

// FieldExtractor interface defines the method for extracting field values from a PodList
// NOTE: There is no actual difference between ListExtractor and MapExtractor,
// since when processing the extracted data, we can always check the type of the underlying data structure
// via FilterInfo.FilterType, but for the sake of clarity, I have defined two separate types.
type FieldExtractor interface {
	Extract(*v1.PodList) any
}

type ListExtractor func(*v1.PodList) interface{}

func (e ListExtractor) Extract(pods *v1.PodList) interface{} {
	return e(pods)
}

type MapExtractor func(*v1.PodList) interface{}

func (e MapExtractor) Extract(pods *v1.PodList) interface{} {
	return e(pods)
}
