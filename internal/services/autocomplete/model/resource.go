package model

import v1 "k8s.io/api/core/v1"

// ResourceType serves as an enum for the supported resource types
type ResourceType string

const (
	PodResourceType ResourceType = "Pod"
)

// Resources is an interface that represents the actual resource type
type Resources interface{}

type PodResources = *v1.PodList
