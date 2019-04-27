# Contributing

* [Development](#development)
  * [GitHub Action](#github-action)

## Development

### GitHub Action

The code for the Action is in the `Dockerfile` and `entrypoint.sh`.

To test the action, [`act`](https://github.com/nektos/act) is an awesome tool for
running Actions locally:

```sh
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
act
```
