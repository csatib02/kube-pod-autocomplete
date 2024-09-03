# Kube Pod Autocomplete documentation

Kube Pod Autocomplete is a service designed to enhance the user experience when navigating resource lists in Kubernetes clusters. It retrieves specific data based on the requested resource type and supported filters.

## Supported resources

| **Resource**     | Available Filters                     |
|------------------|---------------------------------------|
| pods ✅          | namespace, phase, labels, annotations |
| services 🟡      | _Coming Soon_                         |

## API-Spec

Checkout the [OpenAPI](openapi.yaml) file or use the [Swagger spec](swagger.html).

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
kubectl port-forward -n kube-pod-autocomplete "svc/$(kubectl get svc -n kube-pod-autocomplete -o jsonpath='{.items[0].metadata.name}')" 8080:8080 1>/dev/null &
```

Hit the endpoint:

```shell
curl -X GET http://localhost:8080/search/autocomplete/pods
```

## Future improvements

- [ ] Add support for `caching`. (While the current implementation is fast with a small number of pods, but there can be problems in production environments.)
- [ ] `Generate` API specification from `OpenAPI spec`. (The current solution is a really simple POC, if the project is later expanded with additional endpoints code generation from the OpenAPI spec should be utilised.)
- [ ] `Improve End-to-End` Tests: (The existing end-to-end test setup is quite basic, using `cmd.Exec()` and port-forward to access the service is rather limited. Future improvements could include using an ingress controller like NGINX for more robust testing.)
- [ ] Add endpoints that can be called with the received suggestions. (E.g.: `search/:resource/:filters` or get the filters from the body.)
- [ ] Enable deploying with [Garden](https://garden.io/). (Garden helps a lot when it comes to manually testing a project, as new features are added, this should also be implemented.)
