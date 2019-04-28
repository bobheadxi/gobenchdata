# Contributing

* [Development](#development)
  * [GitHub Action](#github-action)
  * [CLI](#cli)

## Development

### GitHub Action

The code for the Action is in the `Dockerfile` and `entrypoint.sh`.

To test the action, [`act`](https://github.com/nektos/act) is an awesome tool for
running Actions locally:

```sh
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
act
```

### CLI

The `gobenchdata` CLI and its associated utilities are written in [Golang](https://golang.org/).
To get started, clone the repository and enable [Go Modules](https://github.com/golang/go/wiki/Modules):

```
export GO111MODULE=on
go mod download
```

Utilities like `gobenchdata-web` are developed in subdirectories under [`/x`](./x).

Code generation tasks should be able to be triggered by [go generate](https://blog.golang.org/generate):

```
go generate ./...
```
