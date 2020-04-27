## Golang log library

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/CodyGuo/glog) [![build](https://img.shields.io/travis/op/go-logging.svg?style=flat)](https://travis-ci.org/CodyGuo/glog)

Package glog implements a log infrastructure for Go.

## Example

Let's have a look at an [example](examples/main.go) which demonstrates most
of the features found in this library.

```go
package main

import (
	"github.com/CodyGuo/glog"
	"os"
)

func main() {
	glog.Debug("hello debug")
	glog.Info("hello info")

	customLog := glog.New(os.Stdout,
		glog.WithLevel(glog.DEBUG),
		glog.WithFlags(glog.LstdFlags),
		glog.WithPrefix("[customLog] "))

	customLog.Debug("hello debug")
	customLog.Info("hello info")
}
```

## Installing

### Using *go get*
    $ go get github.com/CodyGuo/glog

You can use `go get -u` to update the package.

## Documentation

For docs, see http://godoc.org/github.com/CodyGuo/glog or run:

    $ godoc github.com/CodyGuo/glog
