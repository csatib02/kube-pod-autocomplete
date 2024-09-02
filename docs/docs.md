# Kube Pod Autocomplete documentation

Kube Pod Autocomplete is a service designed to enhance the user experience when navigating resource lists in Kubernetes clusters. It can retrieve specific data based on the requested resource type and supported filters.

## Supported resources

| **Resource**     | available Filters                     |
|------------------|---------------------------------------|
| pods ✅          | namespace, phase, labels, annotations |
| services 🟡      | _Next up_                             |

## API-Spec

Checkout the [openapi](openapi.yaml) file, or use the [swagger spec](swagger.html).

## Quick-start

Create a Cluster:

```shell
make up
```

Deploy Kube Pod Autocomplete:

```shell
make deploy
```

Deploy some pods:

```shell
make deploy-testdata
```

Port-forward to Kube Pod Autocomplete:

```shell
kubectl port-forward -n kube-pod-autocomplete $(kubectl get pods -n kube-pod-autocomplete -o jsonpath='{.items[0].metadata.name}') 8080:8080 1>/dev/null &
```

Hit the endpoint:

```shell
curl -X GET http://localhost:8080/search/autocomplete/pods
```

## Future improvements

- [ ] Add support for caching. (Based on tests with the current implementation with small amount of pods, the response speed is blazingly fast, but there can be problems in production environments.)
- [ ] Generate API specification from OpenAPI spec. (The current implementation is a really simple POC, if the project is later expanded with additional endpoints code generation should be utilised.)
