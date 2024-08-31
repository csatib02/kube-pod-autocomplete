package k8s

type podInfo struct {
	Name      string
	Namespace string
	Phase     string
	Labels    map[string]string
}
