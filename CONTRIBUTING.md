# Contributing

* [Development](#development)
  * [GitHub Action](#github-action)
  * [CLI](#cli)
  * [Web App](#web-app)

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

### Web App

The web app is in [x/gobenchdata-web/web](./x/gobenchdata-web/web), and the
generator is in [x/gobenchdata-web](./x/gobenchdata-web). Assets are compiled
using [`fileb0x`](https://github.com/UnnoTed/fileb0x) (see previous section).

The web app should remain as simple as possible - right now it only consists of
3 files (the base HTML, a JavaScript app, and a plain CSS stylesheet), and
ideally it'll stay that way.

To test the web app, add a `benchmarks.json` (for example, the demo data available
in [`gh-pages`](https://github.com/bobheadxi/gobenchdata/blob/gh-pages/benchmarks.json))
to the `web` directory, and run:

```
make serve
```

This requires [serve](https://www.npmjs.com/package/serve) to be installed.
