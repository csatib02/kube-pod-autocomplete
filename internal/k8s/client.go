package k8s

import (
	"context"
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/csatib02/kube-pod-autocomplete/pkg/common"
)

type Client struct {
	clientset kubernetes.Interface
}

func NewClient() (*Client, error) {
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &Client{clientset: clientset}, nil
}

func (c *Client) ListResource(ctx context.Context, resource common.ResourceType) (common.Resources, error) {
	switch resource {
	case common.PodResourceType:
		return c.listPods(ctx)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resource)
	}
}

func (c *Client) listPods(ctx context.Context) (*v1.PodList, error) {
	pods, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	// Validate whether there are any pods in the cluster
	if pods == nil {
		return nil, errors.New("no pods found in the cluster")
	}

	return pods, nil
}
