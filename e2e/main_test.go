//go:build e2e

package e2e

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/support/kind"
	"sigs.k8s.io/e2e-framework/third_party/helm"
)

const defaultTimeout = 2 * time.Minute

var testenv env.Environment

type reverseFinishEnvironment struct {
	env.Environment

	finishFuncs []env.Func
}

// Run launches the test suite from within a TestMain.
func (e *reverseFinishEnvironment) Run(m *testing.M) int {
	e.Environment.Finish(e.finishFuncs...)

	return e.Environment.Run(m)
}

// Finish registers funcs that are executed at the end of the test suite in a reverse order.
func (e *reverseFinishEnvironment) Finish(f ...env.Func) env.Environment {
	e.finishFuncs = append(f[:], e.finishFuncs...)

	return e
}

func TestMain(m *testing.M) {
	testenv = &reverseFinishEnvironment{Environment: env.New()}

	flags := flag.NewFlagSet("", flag.ContinueOnError)
	klog.InitFlags(flags)
	flags.Parse([]string{"-v", "4"})
	log.SetLogger(klog.NewKlogr())

	clusterName := envconf.RandomName("kube-pod-autocomplete-test", 32)
	kindCluster := kind.NewProvider()
	if v := os.Getenv("KIND_K8S_VERSION"); v != "" {
		kindCluster.WithOpts(kind.WithImage("kindest/node:" + v))
		testenv.Setup(envfuncs.CreateClusterWithConfig(kindCluster, clusterName, "kind.yaml"))
	} else {
		testenv.Setup(envfuncs.CreateCluster(kindCluster, clusterName))
	}
	testenv.Finish(envfuncs.DestroyCluster(clusterName))

	if image := os.Getenv("LOAD_IMAGE"); image != "" {
		testenv.Setup(envfuncs.LoadDockerImageToCluster(clusterName, image))
	}

	if imageArchive := os.Getenv("LOAD_IMAGE_ARCHIVE"); imageArchive != "" {
		testenv.Setup(envfuncs.LoadImageArchiveToCluster(clusterName, imageArchive))
	}

	testenv.Setup(envfuncs.CreateNamespace("kube-pod-autocomplete"), envfuncs.CreateNamespace("prod"), envfuncs.CreateNamespace("staging"), installKPA)
	testenv.Finish(uninstallKPA)

	os.Exit(testenv.Run(m))
}

func installKPA(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
	chart := "../deploy/charts/kube-pod-autocomplete"
	if v := os.Getenv("HELM_CHART"); v != "" {
		chart = v
	}

	version := "latest"
	if v := os.Getenv("VERSION"); v != "" {
		version = v
	}

	err := helm.New(cfg.KubeconfigFile()).RunInstall(
		helm.WithName("kube-pod-autocomplete"),
		helm.WithChart(chart),
		helm.WithNamespace("kube-pod-autocomplete"),
		helm.WithArgs("-f", "deploy/kube-pod-autocomplete/values.yaml", "--set", "image.tag="+version),
		helm.WithWait(),
		helm.WithTimeout(defaultTimeout.String()),
	)
	if err != nil {
		return ctx, fmt.Errorf("installing kube-pod-autocomplete: %w", err)
	}

	return ctx, nil
}

func uninstallKPA(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
	err := helm.New(cfg.KubeconfigFile()).RunUninstall(
		helm.WithName("kube-pod-autocomplete"),
		helm.WithNamespace("kube-pod-autocomplete"),
		helm.WithTimeout(defaultTimeout.String()),
	)
	if err != nil {
		return ctx, fmt.Errorf("uninstalling kube-pod-autocomplete: %w", err)
	}

	return ctx, nil
}
