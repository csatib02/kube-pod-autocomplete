package autocomplete

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"

	"github.com/csatib02/kube-pod-autocomplete/internal/k8s"
	"github.com/csatib02/kube-pod-autocomplete/internal/services/autocomplete/model"
)

type Service struct {
	k8sClient k8s.Client
	// TODO: Enable caching
	// useCache          bool
	// cacheUpdatePeriod time.Duration
}

func NewAutoCompleteService() (*Service, error) {
	client, err := k8s.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %w", err)
	}

	return &Service{
		k8sClient: *client,
	}, nil
}

// GetAutocompleteSuggestions returns a list of suggestions for the given query
func (s *Service) GetAutocompleteSuggestions(ctx context.Context, _ string) ([]model.Suggestion, error) {
	// TODO: Implement the logic to fetch suggestions for the given query

	pods, err := s.k8sClient.ListPods(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	namespaces, phases, labels := s.extractPodData(pods)

	suggestions := []model.Suggestion{
		{Key: "namespace", Values: namespaces.ToSlice()},
		{Key: "phase", Values: phases.ToSlice()},
	}
	suggestions = append(suggestions, s.formatLabelSuggestions(labels)...)

	return suggestions, nil
}

// extractPodData extracts namespaces, phases, and labels information from the pods
func (s *Service) extractPodData(pods *v1.PodList) (
	*model.Set[string], *model.Set[string], map[string]*model.Set[string],
) {
	namespaces := model.NewSet[string]()
	phases := model.NewSet[string]()
	labels := make(map[string]*model.Set[string])

	for _, pod := range pods.Items {
		namespaces.Add(pod.Namespace)
		phases.Add(string(pod.Status.Phase))
		for key, value := range pod.Labels {
			if labels[key] == nil {
				labels[key] = model.NewSet[string]()
			}
			labels[key].Add(value)
		}
	}

	return namespaces, phases, labels
}

func (s *Service) formatLabelSuggestions(labels map[string]*model.Set[string]) []model.Suggestion {
	labelSuggestions := make([]model.Suggestion, 0, len(labels))
	for key, labelSet := range labels {
		labelSuggestions = append(labelSuggestions, model.Suggestion{
			Key:    fmt.Sprintf("labels:%s", key),
			Values: labelSet.ToSlice(),
		})
	}

	return labelSuggestions
}
