# kube-pod-autocomplete

*Kube Pod Autocomplete is a Go-based backend service designed to enhance the user experience when navigating resource lists in Kubernetes clusters.*

## Getting started

- Kube Pod Autocomplete is designed to be used in Kubernetes environment.
- Take a look at the [documentation](./docs/docs.md).

## TODO

- Add e2e-tests.
- Consider adding garden config to simplify testing.

- Add search pods by label/ns/phase endpoint as a possible use-case.

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
