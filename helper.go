package ginprometheus

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func computeRequestSize(request *http.Request) float64 {
	return float64(computeRequestSummarySize(request) + computeRequestBodySize(request))
}

func computeRequestSummarySize(request *http.Request) int {
	s := 0
	if request.URL != nil {
		s = len(request.URL.String())
	}

	s += len(request.Method)
	s += len(request.Proto)
	for name, values := range request.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(request.Host)

	if request.ContentLength != -1 {
		s += int(request.ContentLength)
	}
	return s
}

func computeRequestBodySize(request *http.Request) int {
	body, _ := ioutil.ReadAll(request.Body)
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return len(body)
}

func urlMapping(c *gin.Context) string {
	url := c.Request.URL.Path

	for _, p := range c.Params {
		url = strings.Replace(url, "/"+p.Value, "/:"+p.Key, 1)
	}

	return url
}
