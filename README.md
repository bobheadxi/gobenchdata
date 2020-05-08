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

- [About](#about)
- [GitHub Action](#github-action)
  - [Setup](#setup)
  - [Configuration](#configuration)
    - [`inputs`](#inputs)
      - [Publishing](#publishing)
      - [Checks](#checks)
    - [`env`](#env)
  - [Pull Request Checks](#pull-request-checks)
  - [Visualisation](#visualisation)
- [`gobenchdata` CLI](#gobenchdata-cli)
- [Development and Contributions](#development-and-contributions)

<br />

## About



It is available as a GitHub action or a command-line application.

<br />

## GitHub Action

`gobenchdata` can be used as GitHub Action for uploading Go benchmark data as
JSON to `gh-pages` and visualizing it with a generated web app or your own web application.

### Setup

For example, in `.github/workflows/push.yml`, using [the new YAML syntax for workflows](https://help.github.com/en/articles/workflow-syntax-for-github-actions):

```yml
TODO
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

The default behaviour of the `gobenchdata` Action is to commit and publish to your repository's `gh-pages` branch.

| Variable             | Default                   | Purpose
| -------------------- | ------------------------- | -------
| `PUBLISH`            | `false`                   | if `true`, publishes results
| `PUBLISH_REPO`       |                           | an alternative repository to publish to
| `PUBLISH_BRANCH`     | `gh-pages`                | branch to publish to
| `PRUNE_COUNT`        | `0`                       | number of past runs to keep (`0` keeps everything)
| `GIT_COMMIT_MESSAGE` | `"add new benchmark run"` | the commit message for the benchmark update
| `BENCHMARKS_OUT`     | `benchmarks.json`         | destination path of benchmark data

##### Checks

The following `inputs` are for enabling [Pull Request Checks](#pull-request-checks):

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
`secrets.GITHUB_TOKEN` do not trigger page builds. This issue can be resolved by using
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
executed.

<details>
<summary>Simple Example</summary>
<p>

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

You can configure the web application using `gobenchdata-web.json`. TODO: documentation

You can output the entire web application (to commit to Github pages, for example) using:

```sh
gobenchdata web generate .
```

TODO - configure:
* chart groups: group of charts
* chart: a set of metrics to render (by name), each metric gets its own chart, matching on packages and benchmark names (using regexp)

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

```sh
go test -bench . -benchmem ./... | gobenchdata --json bench.json
```

You can then visualize the benchmark using the built-in web application:

```sh
gobenchdata web serve
```

You can further customize the visualization by generating a configuration file:

```sh
gobenchdata web generate --web.config .
```

More detailed usage documentation and examples can be found in the
[godocs](https://godoc.org/go.bobheadxi.dev/gobenchdata) or by running
`gobenchdata help`.

<br />

## Development and Contributions

Please report bugs and requests in the [repository issues](https://go.bobheadxi.dev/gobenchdata)!

See [CONTRIBUTING.md](./CONTRIBUTING.md) for more detailed development documentation.
