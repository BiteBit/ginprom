# Installl

```sh
$ go get github.com/BiteBit/ginprom
```

# Usage

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/BiteBit/ginprom"
)

var (
  prom = ginprom.New("namespace", "module")
)

func main() {
	router := gin.New()

  router.Use(prom.Handler())

	router.Get("/metrics", prom.Metrics())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello world!")
	})

	router.Run(":8080")
}
```