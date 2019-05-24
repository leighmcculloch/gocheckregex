# gocheckregex

[![Build Status](https://img.shields.io/travis/leighmcculloch/gocheckregex.svg)](https://travis-ci.org/leighmcculloch/gocheckregex)
[![Codecov](https://img.shields.io/codecov/c/github/leighmcculloch/gocheckregex.svg)](https://codecov.io/gh/leighmcculloch/gocheckregex)
[![Go Report Card](https://goreportcard.com/badge/github.com/leighmcculloch/gocheckregex)](https://goreportcard.com/report/github.com/leighmcculloch/gocheckregex)

Check that regular expressions in global `MustCompile` calls are valid.

## Why

It's a common practice to compile regexs on application start by using the `regexp.MustCompile` function to build the regular expression and store it in a global variable. The `MustCompile` function will panic then on startup if the regular expression can't be compiled. It would be great if we found out sooner without running the application if the regex we've written is invalid.

## Install

```
go get 4d63.com/gocheckregex
```

## Usage

```
gocheckregex
```

```
gocheckregex ./...
```

```
gocheckregex [path] [path] [path] [etc]
```

Note: Paths are only inspected recursively if the Go `/...` recursive path suffix is appended to the path.

## Limitations

- Does not pickup on regexp usage if package `regexp` is aliased.
- Only looks at MustCompile calls that's argument is a string literal.
