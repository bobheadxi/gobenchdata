# Contributing

* [Development](#development)
  * [Testing the Action](#testing-the-action)

## Development

### Testing the Action

The code for the Action is in the `Dockerfile` and `entrypoint.sh`. There is a
symlink of these to [`/actions/gh-pages`](./.actions/gh-pages) since to publish
on the GitHub Marketplace, the action
[must be in the root directory](https://developer.github.com/marketplace/actions/publishing-an-action-in-the-github-marketplace/),
but to use it in a [local workflow](./.github/main.workflow) it seems like it
must be in a subdirectory. Pending more research and a better solution.

To test the action, [`act`](https://github.com/nektos/act) is an awesome tool for
running Actions locally:

```sh
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
act
```
