package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Client struct {
	clientset *kubernetes.Clientset
}

func NewClient(clientset *kubernetes.Clientset) *Client {
	return &Client{
		clientset: clientset,
	}
}

func (c *Client) ListPods(ctx context.Context) ([]podInfo, error) {
	pods, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	// Validate whether there are any pods in the cluster
	if pods == nil {
		return nil, fmt.Errorf("failed to list pods: no pods found")
	}

	return extractPodInfos(pods), nil
}

// extractPodInfos extracts pod information from a list of pods
func extractPodInfos(pods *v1.PodList) []podInfo {
	var podInfos []podInfo
	for _, pod := range pods.Items {
		podInfo := podInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Phase:     string(pod.Status.Phase),
			Labels:    pod.Labels,
		}
		podInfos = append(podInfos, podInfo)
	}
	return podInfos
}
