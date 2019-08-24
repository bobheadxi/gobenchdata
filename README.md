
# 📉 gobenchdata

[![Build Status](https://dev.azure.com/bobheadxi/bobheadxi/_apis/build/status/bobheadxi.gobenchdata?branchName=master)](https://dev.azure.com/bobheadxi/bobheadxi/_build/latest?definitionId=7&branchName=master)
[![View Action](https://img.shields.io/badge/view-action-yellow.svg)](https://github.com/marketplace/actions/gobenchdata-to-gh-pages)
[![GoDoc](https://godoc.org/go.bobheadxi.dev/gobenchdata?status.svg)](https://godoc.org/go.bobheadxi.dev/gobenchdata)
[![Demo](https://img.shields.io/website/https/gobenchdata.bobheadxi.dev.svg?down_color=grey&down_message=offline&label=demo&up_message=live)](https://gobenchdata.bobheadxi.dev/)

a tool for inspecting `go test -bench` data, and a
[GitHub Action](https://github.com/features/actions) for continuous benchmarking.

<a href="https://gobenchdata.bobheadxi.dev/" target="_blank">
  <img align="right" width="500" src="./.static/demo-chart.png" alt="example">
</a>

* [About](#about)
* [GitHub Action](#github-action)
  * [Setup](#setup)
  * [Configuration](#configuration)
  * [Visualisation](#visualisation)
* [`gobenchdata` CLI](#gobenchdata-cli)
* [Development and Contributions](#development-and-contributions)

<br />

## About

`gobenchdata` was inspired by the [deno.land](https://deno.land/benchmarks.html)
continuous benchmarks, which aims to display performance improvements/regressions
on a continuous basis.

It is available as a GitHub action or a command-line application.

## GitHub Action

`gobenchdata` can be used as GitHub Action for uploading Go benchmark data as
JSON to `gh-pages` and visualizing it with a generated web app (using `gobenchdata-web`)
or your own web application.

### Setup

For example, in `.github/workflows/push.yml`:

```yml
on: push
name: gobenchdata benchmarking
jobs:
  filter:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: filter
      uses: actions/bin/filter@master
      with:
        args: branch master
    - name: gobenchdata to gh-pages
      uses: bobheadxi/gobenchdata@v0.2.0
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        PRUNE: "20"
```

Learn more about GitHub Actions in the [official documentation](https://github.com/features/actions).

### Configuration

| Variable             | Default                   | Purpose
| -------------------- | ------------------------- | -------
| `GITHUB_TOKEN`       | set by GitHub             | token to provide access to repository
| `GITHUB_ACTOR`       | set by GitHub             | the user to make commits as
| `GIT_COMMIT_MESSAGE` | `"add new benchmark run"` | the commit message for the benchmark update
| `GO_BENCHMARKS`      | `.`                       | benchmarks to run (argument for `-bench`)
| `GO_BENCHMARK_FLAGS` |                           | additional flags for `go test`
| `GO_BENCHMARK_PKGS`  | `./...`                   | packages to test (argument for `go test`)
| `FINAL_OUTPUT`       | `benchmarks.json`         | destination path of benchmark data
| `PRUNE_COUNT`        | `0`                       | number of past runs to keep (`0` keeps everything)

### Visualisation

The `gobenchdata` GitHub action eventually generates a JSON file with past benchmarks.
You can visualize these continuous benchmarks by creating a web app that reads
from the JSON benchmarks file, or by using `gobenchdata-web`:

```
go get -u go.bobheadxi.dev/gobenchdata/x/gobenchdata-web
git checkout gh-pages
gobenchdata-web --title "my benchmarks" # generates a web app in your working directory
```

You test the web app locally using a tool like [serve](https://www.npmjs.com/package/serve):

```
serve .
```

This feature is a work in progress. An example site published by this repository is
available at [gobenchdata.bobheaxi.dev](https://gobenchdata.bobheadxi.dev/).

## `gobenchdata` CLI

`gobenchdata`, which the GitHub Action leverages to manage benchmark data,
is also available as a CLI:

```
go get -u go.bobheadxi.dev/gobenchdata
gobenchdata help
```

Usage documentation can be found in the
[godocs](https://godoc.org/go.bobheadxi.dev/gobenchdata).

## Development and Contributions

Please report bugs and requests in the [repository issues](https://go.bobheadxi.dev/gobenchdata)!

See [CONTRIBUTING.md](./CONTRIBUTING.md) for more detailed development documentation.