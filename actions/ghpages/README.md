# gobenchdata to gh-pages

A GitHub Action for uploading Go benchmark data to `gh-pages` using `gobenchdata`.

## Example

In `main.workflow`:

```
workflow "Benchmark Demo" {
  on = "push"
  resolves = ["gobenchdata to gh-pages"]
}

action "filter" {
  uses = "actions/bin/filter@master"
  args = "branch master"
}

action "gobenchdata to gh-pages" {
  uses = "./actions/ghpages"
  needs = ["filter"]
  secrets = ["GITHUB_TOKEN"]
}
```

## Configuration

| Variable            | Default           | Purpose
| ------------------- | ----------------- | -------
| `GITHUB_TOKEN`      | set by GitHub     | token to provide access to repository
| `GITHUB_ACTOR`      | set by GitHub     | the user to make commits as
| `GO_BENCHMARKS`     | `.`               | benchmarks to run (argument for `-bench`)
| `GO_BENCHMARK_PKGS` | `./...`           | packages to test (argument for `go test`)
| `FINAL_OUTPUT`      | `benchmarks.json` | destination path of benchmark data
