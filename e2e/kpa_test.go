//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

const (
	pollInterval    = 5 * time.Second
	healthURL       = "http://%s:%d/health"
	autocompleteURL = "http://%s:%d/search/autocomplete/%s"
)

func TestKPAEndpoints(t *testing.T) {
	endpoints := applyResource(features.New("validate endpoint functionality")).
		Assess("pods are available", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			// check if test pods are running
			err := wait.PollUntilContextTimeout(ctx, pollInterval, defaultTimeout, true, func(ctx context.Context) (bool, error) {
				// get all pods with label: team=test
				pods := &corev1.PodList{}
				err := cfg.Client().Resources().List(ctx, pods, func(opts *metav1.ListOptions) {
					opts.LabelSelector = labels.Set{"team": "test"}.String()
				})
				require.NoError(t, err)

				for _, pod := range pods.Items {
					if pod.Status.Phase != corev1.PodRunning {
						return false, nil
					}
				}

				return true, nil
			})
			require.NoError(t, err)

			return ctx
		}).Assess("check KPA response", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {

		// start port forwarding
		killProcess := startPortForwardingToService(t, "kube-pod-autocomplete-service", "kube-pod-autocomplete", "8080:8080")
		defer killProcess()

		// hit the /health endpoint
		resp, err := http.Get(fmt.Sprintf(healthURL, "localhost", 8080))
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		// hit the /search/autocomplete/pods endpoint
		autocompleteUrl := fmt.Sprintf(autocompleteURL, "localhost", 8080, "pods")
		resp, err = http.Get(autocompleteUrl)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		// check if the response body contains the expected filters
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Contains(t, string(body), "namespace")
		require.Contains(t, string(body), "phase")
		require.Contains(t, string(body), "labels")
		require.Contains(t, string(body), "annotations")

		return ctx
	}).Feature()

	testenv.Test(t, endpoints)
}

func applyResource(builder *features.FeatureBuilder) *features.FeatureBuilder {
	return builder.
		Setup(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			err := decoder.DecodeEachFile(
				ctx, os.DirFS("test"), "*",
				decoder.CreateHandler(cfg.Client().Resources()),
			)
			require.NoError(t, err)

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			err := decoder.DecodeEachFile(
				ctx, os.DirFS("test"), "*",
				decoder.DeleteHandler(cfg.Client().Resources()),
			)
			require.NoError(t, err)

			return ctx
		})
}

func startPortForwardingToService(t *testing.T, svcName, ns, portMapping string) func() {
	args := []string{"port-forward", fmt.Sprintf("svc/%s", svcName), portMapping, "-n", ns}
	cmd := exec.Command("kubectl", args...)
	cmd.Stderr = os.Stderr // redirect stderr to test output
	err := cmd.Start()
	require.NoError(t, err)

	// Wait for port forwarding to be established
	time.Sleep(pollInterval)

	return func() {
		_ = cmd.Process.Kill()
	}
}
