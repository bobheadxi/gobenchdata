# Contributing

* [Development](#development)
  * [Testing the Action](#testing-the-action)

## Development

### Testing the Action

The code for the Action is in [`actions/gh-pages`](./actions/gh-pages). There is
a symlink from the `Dockerfile` and `entrypoint.sh` since to publish on the
GitHub Marketplace, the action
[must be in the root directory](https://developer.github.com/marketplace/actions/publishing-an-action-in-the-github-marketplace/).

To test the action, [`act`](https://github.com/nektos/act) is an awesome tool for
running Actions locally:

```sh
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
act
```
