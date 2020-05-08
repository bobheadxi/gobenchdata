# gobenchdata

> This is a work-in-progress branch for `bobheadxi/gobenchdata@v1`

[![Build Status](https://dev.azure.com/bobheadxi/bobheadxi/_apis/build/status/bobheadxi.gobenchdata?branchName=master)](https://dev.azure.com/bobheadxi/bobheadxi/_build/latest?definitionId=7&branchName=master)
[![View Action](https://img.shields.io/badge/view-github%20action-yellow.svg)](https://bobheadxi.dev/r/gobenchdata)
[![GoDoc](https://img.shields.io/badge/go.pkg.dev-reference-5272B4)](https://pkg.go.dev/go.bobheadxi.dev/gobenchdata)
[![Demo](https://img.shields.io/website/https/gobenchdata.bobheadxi.dev.svg?down_color=grey&down_message=offline&label=demo&up_message=live)](https://gobenchdata.bobheadxi.dev/)
[![Demo Benchmarks](https://github.com/bobheadxi/gobenchdata/workflows/gobenchdata%20demo/badge.svg)](https://github.com/bobheadxi/gobenchdata/blob/master/.github/workflows/push.yml)

`gobenchdata`  is a tool for inspecting `go test -bench` data, and a [GitHub Action](https://github.com/features/actions) for continuous benchmarking. was inspired by the [deno.land continuous benchmarks](https://deno.land/benchmarks.html), which aims to display performance improvements and regressions on a continuous basis.

> :wave: I am currently working on `gobenchdata@v1`, which will bring significant changes to the GitHub Action and web visualization - if you currently use `gobenchdata`, I'd love to hear from you over at the [`gobenchdata@v1` tracking issue](https://github.com/bobheadxi/gobenchdata/issues/36)!

<a href="https://gobenchdata.bobheadxi.dev/" target="_blank">
  <img align="right" width="500" src="./.static/demo-chart.png" alt="example">
</a>

- [GitHub Action](#github-action)
  - [Setup](#setup)
  - [Configuration](#configuration)
    - [`inputs`](#inputs)
      - [Publishing](#publishing)
      - [Checks](#checks)
    - [`env`](#env)
  - [Pull Request Checks](#pull-request-checks)
  - [Visualisation](#visualisation)
- [Command Line Interface](#command-line-interface)
- [Development and Contributions](#development-and-contributions)

<br />

<br />

## GitHub Action

`gobenchdata` can be used as GitHub Action for uploading Go benchmark data as
JSON to `gh-pages` and visualizing it with a generated web app or your own web application.

### Setup

For example, in `.github/workflows/push.yml`, using [the new YAML syntax for workflows](https://help.github.com/en/articles/workflow-syntax-for-github-actions):

```yml
name: gobenchdata publish
on: push
jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
    - name: checkout
      uses: actions/checkout@v2
    - name: gobenchdata publish
      uses: ./
      with:
        PRUNE_COUNT: 30
        GO_TEST_FLAGS: -cpu 1,2
        PUBLISH: true
        PUBLISH_BRANCH: gh-pages
      env:
        GITHUB_TOKEN: ${{ secrets.ACCESS_TOKEN }}
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

##### Publishing

The following `inputs` enable publishing - this merges and publishes benchmark results to a
repository and branch of your choice. This is most useful in conjunction with the [`gobenchdata` web application](#visualisation).

| Variable             | Default                   | Purpose
| -------------------- | ------------------------- | -------
| `PUBLISH`            | `false`                   | if `true`, publishes results
| `PUBLISH_REPO`       |                           | an alternative repository to publish to
| `PUBLISH_BRANCH`     | `gh-pages`                | branch to publish to
| `PRUNE_COUNT`        | `0`                       | number of past runs to keep (`0` keeps everything)
| `GIT_COMMIT_MESSAGE` | `"add new benchmark run"` | the commit message for the benchmark update
| `BENCHMARKS_OUT`     | `benchmarks.json`         | destination path of benchmark data

##### Checks

The following `inputs` are for enabling [Pull Request Checks](#pull-request-checks), which allow
you to watch for performance regressions in your pull requests.

| Variable             | Default                   | Purpose
| -------------------- | ------------------------- | -------
| `CHECKS`             | `false`                   | if `true`, runs checks and sets JSON results to `checks-results`
| `CHECKS_CONFIG`      | `gobenchdata-checks.yml`  | path to checks configuration
| `PUBLISH_REPO`       |                           | repository of benchmark data to check against
| `PUBLISH_BRANCH`     | `gh-pages`                | branch of benchmark data to check against
| `BENCHMARKS_OUT`     | `benchmarks.json`         | path to benchmark data to check against

#### `env`

Environment variables are configured using
[`jobs.<job_id>.steps.env`](https://help.github.com/en/articles/workflow-syntax-for-github-actions#jobsjob_idstepsenv).

| Variable             | Recommended                   | Purpose
| -------------------- | ----------------------------- | -------
| `GITHUB_TOKEN`       | `${{ secrets.GITHUB_TOKEN }}` | token to provide access to repository
| `GITHUB_ACTOR`       | set by GitHub                 | the user to make commits as

Note that for `GITHUB_TOKEN`, it seems that pushes to `gh-pages` made by the default
`secrets.GITHUB_TOKEN` might not trigger page builds. This issue can be resolved by using
a [personal access token](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line)
instead.

### Pull Request Checks

Instead of publishing results, benchmark output can be used to pass and fail pull requests
using `CHECKS: true`. To get started, set up the checks configuration:

```sh
go get -u go.bobheadxi.dev/gobenchdata
gobenchdata checks generate
```

This will generate a file, `gobenchdata-checks.yml`, where you can configure what checks are
executed. The checks are run against any benchmarks that match given `package` and `benchmarks`
values, which should be provided as [regular expressions](https://regexr.com/).

<details>
<summary>Simple Example</summary>
<p>

<!-- copy from './gobenchdata checks generate --checks.config ./tmp/gobenchdata-checks.yml' -->

```yml
checks:
- name: My Check
  description: |-
    Define a check here - in this example, we caculate % difference for NsPerOp in the diff function.
    diff is a function where you receive two parameters, current and base, and in general this function
    should return a negative value for an improvement and a positive value for a regression.
  package: .
  benchmarks: []
  diff: (current.NsPerOp - base.NsPerOp) / base.NsPerOp * 100
  thresholds:
    max: 10
```

</p>
</details>

### Visualisation

The `gobenchdata` GitHub action eventually generates a JSON file with past benchmarks.
You can visualize these continuous benchmarks by creating a web app that reads
from the JSON benchmarks file, or by using `gobenchdata`. An easy way to get started is:

```sh
go get -u go.bobheadxi.dev/gobenchdata
gobenchdata web generate --web.config-only .
gobenchdata web serve # opens visualization in browser
```

You can configure the web application using `gobenchdata-web.yml`. The configuration allows
you to define groups of charts, where each group can be used to compare a set of benchmarks.
Benchmarks are selected with [regular expressions](https://regexr.com/) by package and benchmark
names provided in the configuration.

Note that in each set of compared benchmarks, every metric will get its own chart. You can
select which metrics to display using the `metrics` option.

<details>
<summary>Example</summary>
<p>

<!-- copy from './web/public/gobenchdata-web.yml' -->

```yml
title: gobenchdata web
description: Benchmarks generated using 'gobenchdata'
repository: https://github.com/bobheadxi/gobenchdata
benchmarksFile: benchmarks.json
chartGroups:
  - name: Demo Benchmarks
    description: |
      This is a demo for gobenchdata, a tool and GitHub action for setting up simple continuous
      benchmarks to monitor performance improvements and regressions in your Golang benchmarks!
    charts:
      - name: specify charts by package
        package: go.bobheadxi.dev\/gobenchdata\/demo
      - name: match on specific benchmarks across packages with glob patterns
        benchmarks: [ 'BenchmarkFib.' ]
  - name: More Demo Benchmarks
    description: Create multiple groups of benchmarks
    charts:
      - name: match by a combination of package and benchmarks
        package: go.bobheadxi.dev\/gobenchdata\/.
        benchmarks: [ 'BenchmarkPizzas.', '.FibSlow.' ]
```

</p>
</details>

You can output the entire web application (to publish to Github pages, for example) using:

```sh
gobenchdata web generate ./app
```

<br />

## Command Line Interface

`gobenchdata`, which the GitHub Action leverages to manage benchmark data,
is also available as a CLI:

```sh
go get -u go.bobheadxi.dev/gobenchdata
gobenchdata help
```

The easiest way to use the CLI is by piping the output of `go test -bench` to
it - `gobenchdata` will consume the output and generate a JSON report for you.

```sh
go test -bench . -benchmem ./... | gobenchdata --json bench.json
```

You can use this report to create your own charts, or just use the [built-in web application](#visualisation):

```sh
gobenchdata web serve
```

`gobenchdata` can also execute checks for you to help you ensure performance
regressions don't happen:

```sh
gobenchdata checks generate
gobenchdata checks eval ${base benchmarks} ${current benchmarks} --checks.pretty
```

For more details on how to use checks, see the [pull request checks documentation](#pull-request-checks).

<br />

## Development and Contributions

Please report bugs and requests in the [repository issues](https://go.bobheadxi.dev/gobenchdata)!

See [CONTRIBUTING.md](./CONTRIBUTING.md) for more detailed development documentation.
