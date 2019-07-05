# approve-bot

[![CircleCI](https://circleci.com/gh/d-kuro/approve-bot.svg?style=svg)](https://circleci.com/gh/d-kuro/approve-bot) [![Go Report Card](https://goreportcard.com/badge/github.com/d-kuro/approve-bot)](https://goreportcard.com/report/github.com/d-kuro/approve-bot)

## About approve-bot

The approve-bot defines an owner and a regular expression pattern and approves Pull Requests based on it.
approve-bot config is defined in YAML. By default, it is read from the file `.approve.yaml`.

Sample config `.approve.yaml`:

```yaml
repo: github.com/d-kuro/approve-bot
owners:
  - name: d-kuro
    patterns:
      - path/to/file
      - regexp
      - cmd/approve.go
      - cmd/[a-z]+.go
  - name: d-kuro-kuro
    patterns:
      - path/to/file
      - regexp
      - cmd/version.go
      - cmd/.+
```

A pattern is selected that matches the owner of the pull request.
If all the modified files match the defined pattern, approve-bot will approve.

approve-bot does not accept webhooks and expects to run in CI or locally.
To run, you need GitHub repo scope token. Specify with `--token` option or `GITHUB_TOKEN` environment variable.

Support CI services:

* [x] CircleCI
* [x] Travis CI
* [ ] GitHub Actions

## Build and Run

Build approve-bot:

```shell
$ make build
```

and run approve-bot:

```shell
$ ./dist/approve-bot --token < your GitHub token for repo scope > --pr https://github.com/d-kuro/approve-bot/pull/1
```

For more information about available options run:

```shell
$ dist/approve-bot -h
Approve Pull Request of the file owner.

Usage:
  approve-bot [flags]
  approve-bot [command]

Examples:

$ approve-bot --token <your GitHub token for repo scope> --pr https://github.com/d-kuro/approve-bot/pull/1

.approve.yaml:
---
owners:
  - name: d-kuro
    patterns:
      - path/to/file
      - regexp
      - cmd/approve.go
      - cmd/[a-z]+.go
---

# Or specify a Pull Request number. "repo" of config is required, when using Pull Request number.
$ approve-bot --token <your GitHub token for repo scope > --prnum 1

.approve.yaml:
---
repo: github.com/d-kuro/approve-bot
owners:
  - name: d-kuro
    patterns:
      - path/to/file
      - regexp
      - cmd/approve.go
      - cmd/[a-z]+.go
---


Available Commands:
  help        Help about any command
  version     Show version

Flags:
      --config string   Config YAML file path. (default ".approve.yaml")
  -h, --help            help for approve-bot
      --pr string       Pull Request URL. Or environment variable ("CIRCLE_PULL_REQUEST")
      --prnum int       Pull Request Number. Or environment variable ("TRAVIS_PULL_REQUEST")
      --token string    GitHub token. Or environment variable ("GITHUB_TOKEN")

Use "approve-bot [command] --help" for more information about a command.
```

## CircleCI example

```yaml
version: 2.1

executors:
  golang:
    docker:
      - image: circleci/golang:1.12.6
        environment:
          GO111MODULE: "on"

workflows:
  workflow:
    jobs:
      - approve

jobs:
  approve:
    executor:
      name: golang
    steps:
      - checkout
      - run:
          name: install approve-bot
          command: go get github.com/d-kuro/approve-bot
      - run:
          name: run approve-bot
          command: approve-bot
```
