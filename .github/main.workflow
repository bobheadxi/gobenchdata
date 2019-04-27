workflow "New workflow" {
  on = "push"
  resolves = ["GitHub Action for Uploading Go Benchmark Data to gh-pages"]
}

action "Filters for GitHub Actions" {
  uses = "actions/bin/filter@3c0b4f0e63ea54ea5df2914b4fabf383368cd0da"
  args = "branch master"
}

action "GitHub Action for Uploading Go Benchmark Data to gh-pages" {
  uses = "./actions/ghpages"
  needs = ["Filters for GitHub Actions"]
}
