package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubernetesConfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

type Client struct {
	clientset kubernetes.Interface
}

func NewClient() (*Client, error) {
	kubeConfig, err := kubernetesConfig.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &Client{clientset: clientset}, nil
}

func (c *Client) ListPods(ctx context.Context) (*v1.PodList, error) {
	pods, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	// Validate whether there are any pods in the cluster
	if pods == nil {
		return nil, fmt.Errorf("failed to list pods: no pods found")
	}

	return pods, nil
}
