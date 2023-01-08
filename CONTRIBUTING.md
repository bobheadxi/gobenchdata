# Contributing

[![pipeline](https://github.com/bobheadxi/gobenchdata/workflows/pipeline/badge.svg)](https://github.com/bobheadxi/gobenchdata/actions?workflow=pipeline)
[![codecov](https://codecov.io/gh/bobheadxi/gobenchdata/branch/master/graph/badge.svg)](https://codecov.io/gh/bobheadxi/gobenchdata)
[![pkg.go.dev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/go.bobheadxi.dev/gobenchdata?tab=doc)
[![Sourcegraph](https://img.shields.io/badge/view%20on-sourcegraph-A112FE?logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAAEZklEQVRoQ+2aXWgUZxSG3292sxtNN43BhBakFPyhxSujRSxiU1pr7SaGXqgUxOIEW0IFkeYighYUxAuLUlq0lrq2iCDpjWtmFVtoG6QVNOCFVShVLyxIk0DVjZLMxt3xTGTccd2ZOd/8JBHci0CY9zvnPPN+/7sCIXwKavOwAcy2QgngQiIztDSE0OwQlDPYR1ebiaH6J5kZChyfW12gRG4QVgGTBfMchMbFP9Sn5nlZL2D0JjLD6710lc+z0NfqSGTXQRQ4bX07Mq423yoBL3OSyHSvUxirMuaEvgbJWrdcvkHMoJwxYuq4INUhyuWvQa1jvdMGxAvCxJlyEC9XOBCWL04wwRzpbDoDQ7wfZJzIQLi5Eggk6DiRhZgWIAbE3NrM4A3LPT8Q7UgqAqLqTmLSHLGPkyzG/qXEczhd0q6RH+zaSBfaUoc4iQx19pIClIscrTkNZzG6gd7qMY6eC2Hqyo705ZfTf+eqJmhMzcSbYtQpOXc92ZsZjLVAL4YNUQbJ5Ttg4CQrQdGYj44Xr9m1XJCzmZusFDJOWNpHjmh5x624a2ZFtOKDVL+uNo2TuXE3bZQQZUf8gtgqP31uI94Z/rMqix+IGiRfWw3xN9dCgVx+L3WrHm4Dju6PXz/EkjuXJ6R+IGgyOE1TbZqTq9y1eo0EZo7oMo1ktPu3xjHvuiLT5AFNszUyDULtWpzE2/fEsey8O5TbWuGWwxrs5rS7nFNMWJrNh2No74s9Ec4vRNmRRzPXMP19fBMSVsGcOJ98G8N3Wl2gXcbTjbX7vUBxLaeASDQCm5Cu/0E2tvtb0Ea+BowtskFD0wvlc6Rf2M+Jx7dTu7ubFr2dnKDRaMQe2v/tcIrNB7FH0O50AcrBaApmRDVwFO31ql3pD8QW4dP0feNwl/Q+kFEtRyIGyaWXnpy1OO0qNJWHo1y6iCmAGkBb/Ru+HenDWIF2mo4r8G+tRRzoniSn2uqFLxANhe9LKHVyTbz6egk9+x5w5fK6ulSNNMhZ/Feno+GebLZV6isTTa6k5qNl5RnZ5u56Ib6SBvFzaWBBVFZzvnERWlt/Cg4l27XChLCqFyLekjhy6xJyoytgjPf7opIB8QPx7sYFiMXHPGt76m741MhCKMZfng0nBOIjmoJPsLqWHwgFpe6V6qtfcopxveR2Oy+J0ntIN/zCWkf8QNAJ7y6d8Bq4lxLc2/qJl5K7t432XwcqX5CrI34gzATWuYILQtdQPyePDK3iuOekCR3Efjhig1B1Uq5UoXEEoZX7d1q535J5S9VOeFyYyEBku5XTMXXKQTToX5Rg7OI44nbW5oKYeYK4EniMeF0YFNSmb+grhc84LyRCEP1/OurOcipCQbKxDeK2V5FcVyIDMQvsgz5gwFhcWWwKyRlvQ3gv29RwWoDYAbIofNyBxI9eDlQ+n3YgsgCWnr4MStGXQXmv9pF2La/k3OccV54JEBM4yp9EsXa/3LfO0dGPcYq0Y7DfZB8nJzZw2rppHgKgVHs8L5wvRwAAAABJRU5ErkJggg==)](https://sourcegraph.com/github.com/bobheadxi/gobenchdata)
[![Go Report Card](https://goreportcard.com/badge/go.bobheadxi.dev/gobenchdata)](https://goreportcard.com/report/go.bobheadxi.dev/gobenchdata)

This document outlines some useful information for developing and contributing to `gobenchdata`.

The core components of `gobenchdata` are:

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
