
# gobenchdata

[![Build Status](https://dev.azure.com/bobheadxi/bobheadxi/_apis/build/status/bobheadxi.gobenchdata?branchName=master)](https://dev.azure.com/bobheadxi/bobheadxi/_build/latest?definitionId=7&branchName=master)
[![View Action](https://img.shields.io/badge/view-github%20action-yellow.svg)](https://bobheadxi.dev/r/gobenchdata)
[![GoDoc](https://img.shields.io/badge/go.pkg.dev-reference-5272B4)](https://pkg.go.dev/go.bobheadxi.dev/gobenchdata)
[![Demo](https://img.shields.io/website/https/gobenchdata.bobheadxi.dev.svg?down_color=grey&down_message=offline&label=demo&up_message=live)](https://gobenchdata.bobheadxi.dev/)
[![Demo Benchmarks](https://github.com/bobheadxi/gobenchdata/workflows/gobenchdata%20demo/badge.svg)](https://github.com/bobheadxi/gobenchdata/blob/master/.github/workflows/push.yml)

a tool for inspecting `go test -bench` data, and a
[GitHub Action](https://github.com/features/actions) for continuous benchmarking.

<a href="https://gobenchdata.bobheadxi.dev/" target="_blank">
  <img align="right" width="500" src="./.static/demo-chart.png" alt="example">
</a>

- [About](#about)
- [GitHub Action](#github-action)
  - [Setup](#setup)
  - [Configuration](#configuration)
    - [`inputs`](#inputs)
    - [`env`](#env)
  - [Visualisation](#visualisation)
- [`gobenchdata` CLI](#gobenchdata-cli)
- [Development and Contributions](#development-and-contributions)

<br />

## About

`gobenchdata` was inspired by the [deno.land continuous benchmarks](https://deno.land/benchmarks.html),
which aims to display performance improvements and regressions on a continuous basis.

It is available as a GitHub action or a command-line application.

<br />

## GitHub Action

`gobenchdata` can be used as GitHub Action for uploading Go benchmark data as
JSON to `gh-pages` and visualizing it with a generated web app (using `gobenchdata-web`)
or your own web application.

### Setup

For example, in `.github/workflows/push.yml`, using [the new YAML syntax for workflows](https://help.github.com/en/articles/workflow-syntax-for-github-actions):

```yml
name: Benchmark
on:
  push:
    branches: [ master ]

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
    - name: checkout
      uses: actions/checkout@v1
      with:
        fetch-depth: 1
    - name: gobenchdata to gh-pages
      uses: bobheadxi/gobenchdata@v0.3.0
      with:
        PRUNE_COUNT: 30
        GO_TEST_FLAGS: -cpu 1,2
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

Learn more about GitHub Actions in the [official documentation](https://github.com/features/actions).

### Configuration

#### `inputs`

Input variables are configured using
[`jobs.<job_id>.steps.with`](https://help.github.com/en/articles/workflow-syntax-for-github-actions#jobsjob_idstepswith).

| Variable             | Default                   | Purpose
| -------------------- | ------------------------- | -------
| `SUBDIRECTORY`       | `.`                       | subdirectory of project to run commands from
| `GO_BENCHMARKS`      | `.`                       | benchmarks to run (argument for `-bench`)
| `GO_TEST_FLAGS`      |                           | additional flags for `go test`
| `GO_TEST_PKGS`       | `./...`                   | packages to test (argument for `go test`)
| `BENCHMARKS_OUT`     | `benchmarks.json`         | destination path of benchmark data
| `PRUNE_COUNT`        | `0`                       | number of past runs to keep (`0` keeps everything)
| `GIT_COMMIT_MESSAGE` | `"add new benchmark run"` | the commit message for the benchmark update

#### `env`

Environment variables are configured using
[`jobs.<job_id>.steps.env`](https://help.github.com/en/articles/workflow-syntax-for-github-actions#jobsjob_idstepsenv).

| Variable             | Recommended                   | Purpose
| -------------------- | ----------------------------- | -------
| `GITHUB_TOKEN`       | `${{ secrets.GITHUB_TOKEN }}` | token to provide access to repository
| `GITHUB_ACTOR`       | set by GitHub                 | the user to make commits as

### Visualisation

The `gobenchdata` GitHub action eventually generates a JSON file with past benchmarks.
You can visualize these continuous benchmarks by creating a web app that reads
from the JSON benchmarks file, or by using `gobenchdata-web`:

```sh
go get -u go.bobheadxi.dev/gobenchdata/x/gobenchdata-web
git checkout gh-pages
gobenchdata-web --title "my benchmarks" # generates a web app in your working directory
```

The generator offers a variety of customization options documented under
`gobenchdata-web help` to configure descriptions, charts, and more. The easiest
way to use `gobenchdata-web` is to set up a Makefile in your `gh-pages` branch
to update the web app using the latest version of `gobenchdata-web` with your
desired configuration - [for example](https://github.com/bobheadxi/gobenchdata/blob/gh-pages/Makefile):

```makefile
all: build
	git commit -a -m "regenerate web app"

build:
	gobenchdata-web --title "gobenchdata continuous benchmark demo" --desc "This is a demo for gobenchdata"
```

You can test the web app locally using a tool like [serve](https://www.npmjs.com/package/serve):

```
serve .
```

The web application generator is a work in progress. An example site published
by this repository is available at [gobenchdata.bobheaxi.dev](https://gobenchdata.bobheadxi.dev/)
([configuration](https://github.com/angeleneleow/gobenchdata/blob/master/.github/workflows/push.yml)).

Other examples:

* [`bobheadxi/zapx`](https://zapx.bobheadxi.dev/benchmarks/)
* [`benchx.temporal.cloud`](https://benchx.temporal.cloud/) by [@RTradeLtd](https://github.com/RTradeLtd/)

<br />

## `gobenchdata` CLI

`gobenchdata`, which the GitHub Action leverages to manage benchmark data,
is also available as a CLI:

```
go get -u go.bobheadxi.dev/gobenchdata
gobenchdata help
```

The easiest way to use the CLI is by piping the output of `go test -bench` to
it:

```
go test -bench . -benchmem ./... | gobenchdata --json bench.json
```

More detailed usage documentation and examples can be found in the
[godocs](https://godoc.org/go.bobheadxi.dev/gobenchdata) or by running
`gobenchdata help`.

<br />

## Development and Contributions

Please report bugs and requests in the [repository issues](https://go.bobheadxi.dev/gobenchdata)!

See [CONTRIBUTING.md](./CONTRIBUTING.md) for more detailed development documentation.
