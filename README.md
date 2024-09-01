# kube-pod-autocomplete

Kube Pod Autocomplete is a Go-based backend service designed to enhance the user experience when navigating resource lists in Kubernetes clusters.

## TODO

- Add caching.
- Add search pods by label/ns/phase endpoint as a possible use-case.
- Add Unit-tests.
- Add e2e-tests.
- Consider moving main.go to cmd.
- Consider adding garden config to simplify testing.
- Create a Helm-Chart for Kube Pod Autocomplete

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
curl http://localhost:8080/search/autocomplete/pods
```

## Development

Make sure Docker is installed.

Fetch required tools:

```shell
make deps
```

Deploy development environment:

```shell
make up
```

Deploy Kube Pod Autocomplete:

```shell
make deploy
```

Run the test suite:

```shell
make test
make test-e2e
```

Run linters:

```shell
make lint # pass -j option to run them in parallel
```

Some linter violations can automatically be fixed:

```shell
make fmt
```

Build artifacts locally:

```shell
make container-image
```

Once you are done, you can tear down the development environment:

```shell
make down
```

## License

The project is licensed under the [Apache 2.0 License](LICENSE).
