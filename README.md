[![Sourcegraph](https://sourcegraph.com/github.com/BiteBit/ginprom/-/badge.svg?style=flat-square)](https://sourcegraph.com/github.com/BiteBit/ginprom?badge)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/BiteBit/ginprom)
[![Go Report Card](https://goreportcard.com/badge/github.com/BiteBit/ginprom?style=flat-square)](https://goreportcard.com/report/github.com/BiteBit/ginprom)
[![Build Status](http://img.shields.io/travis/BiteBit/ginprom.svg?style=flat-square)](https://travis-ci.org/BiteBit/ginprom)
[![Codecov](https://img.shields.io/codecov/c/github/BiteBit/ginprom.svg?style=flat-square)](https://codecov.io/gh/BiteBit/ginprom)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/BiteBit/ginprom/master/LICENSE)


# Installl

```sh
$ go get github.com/BiteBit/ginprom
```

# Usage

```go
package main

import (
	"github.com/BiteBit/ginprom"
	"github.com/gin-gonic/gin"
)

var (
	prom = ginprom.New("namespace", "module")
)

func main() {
	router := gin.New()

	router.Use(prom.Handler())
	router.GET("/metrics", prom.Metrics())
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello world!")
	})

	router.Run(":8080")
}
```