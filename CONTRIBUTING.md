# Contributing

[![pipeline](https://github.com/bobheadxi/gobenchdata/workflows/pipeline/badge.svg)](https://github.com/bobheadxi/gobenchdata/actions?workflow=pipeline)
[![codecov](https://codecov.io/gh/bobheadxi/gobenchdata/branch/master/graph/badge.svg)](https://codecov.io/gh/bobheadxi/gobenchdata)
[![pkg.go.dev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/go.bobheadxi.dev/gobenchdata?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/go.bobheadxi.dev/gobenchdata)](https://goreportcard.com/report/go.bobheadxi.dev/gobenchdata)

- [Action](#action)
- [CLI](#cli)
- [Web](#web)

## Action

The code for the Action is in the `Dockerfile` and `entrypoint.sh`. It is
continuously tested by the [demo workflow](https://github.com/bobheadxi/gobenchdata/blob/master/.github/workflows/push.yml).

[`act`](https://github.com/nektos/act) is a great tool to test Actions locally.

## CLI

The `gobenchdata` CLI and its associated utilities are written in [Golang](https://golang.org/).
To get started, clone the repository and enable [Go Modules](https://github.com/golang/go/wiki/Modules):

```sh
export GO111MODULE=on
go mod download
make # install binary
```

Code generation tasks should be able to be triggered by [`go generate`](https://blog.golang.org/generate),
but some tasks don't seem to work with it, so run the following to run all `go generate`
tasks as well as any other code generation scripts:

```sh
make generate
```

The example benchmarks can be run using `make bench`.

## Web

The web app and the web app generator are both in [/web](./web). Assets are compiled
using [`fileb0x`](https://github.com/UnnoTed/fileb0x) (see previous section). The app is built using [Vue.js](https://vuejs.org/).

To test the web app, add a `benchmarks.json` (for example, the demo data available
in [`gh-pages`](https://go.bobheadxi.dev/gobenchdata/blob/gh-pages/benchmarks.json))
to the `web/public` directory, and run the following in `web`:

```sh
npm install
npm run serve
```

An example configuration is provided in [`web/public/gobenchdata-web.yml`](./web/public/gobenchdata-web.yml)
that should allow you to test most of the app's features.

To generate benchmarks from scratch, run:

```sh
make demo
```

This can be run repeatedly to make very large `benchmark.json` run histories.
