package common

import v1 "k8s.io/api/core/v1"

// ResourceType serves as an enum for the supported resource types
type ResourceType string

const (
	PodResourceType ResourceType = "pods"
)

// Resources represents the actual resource type
type Resources any

type PodResources = *v1.PodList

var resourceTypeMap = map[ResourceType]Resources{
	PodResourceType: &v1.PodList{},
}

func IsValidResourceType(resourceType ResourceType) bool {
	_, ok := resourceTypeMap[resourceType]
	return ok
}
