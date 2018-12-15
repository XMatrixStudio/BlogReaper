# BlogReaper

A blog reaper, not a blog system.

[![Build Status](https://travis-ci.com/XMatrixStudio/BlogReaper.svg?branch=master)](https://travis-ci.com/XMatrixStudio/BlogReaper)
[![Coverage Status](https://coveralls.io/repos/github/XMatrixStudio/BlogReaper/badge.svg)](https://coveralls.io/github/XMatrixStudio/BlogReaper)
[![CodeFactor](https://www.codefactor.io/repository/github/xmatrixstudio/blogreaper/badge)](https://www.codefactor.io/repository/github/xmatrixstudio/blogreaper)
[![GoDoc](https://godoc.org/github.com/XMatrixStudio/BlogReaper?status.svg)](https://godoc.org/github.com/XMatrixStudio/BlogReaper)

[[Schema](https://github.com/XMatrixStudio/BlogReaper/blob/master/graphql/schema.graphql)]

## Quick Start

BlogReaper uses [Violet](https://oauth.xmatrix.studio/) as user system.

In order to build and run your BlogReaper, you need your own Violet application id and key. **But we are sorry that Violet v2 doesn't provide any way to register an application for others outsides our studio.**

Maybe Violet v3 supports it.

## Installation

Install BlogReaper.

```sh
$ go get -u -v github.com/XMatrixStudio/BlogReaper
```

Copy the `config/` folder into `$GOPATH/bin/`, and rename `config.example.yaml` as `config.yaml`.

For configure file, you need to input your Violet application id and key.

And run BlogReaper.

```sh
$GOPATH/bin/BlogReaper
```

## Development

BlogReaper uses `gqlgen` to generate GraphQL code.

```sh
$ go run ./scripts/gqlgen.go -v
```

## License

For v1, none.

