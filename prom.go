package ginprom

import (
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type (
	// Prom is prometheus helper
	Prom struct {
		namespace           string
		subsystem           string
		pushTargetURL       string
		pushJobName         string
		pushInterval        time.Duration
		requestURLMappingFn func(*gin.Context) string
		reqCounter          *prometheus.CounterVec
		reqDurationSummary  *prometheus.SummaryVec
		reqSizeSummary      *prometheus.SummaryVec
		resSizeSummary      *prometheus.SummaryVec
	}
)

var (
	labelNames = []string{"code", "errcode", "method", "url", "handler"}
)

// New a prom instance
func New(namesapce, subSystem string) *Prom {
	prom := &Prom{
		namespace: namesapce,
		subsystem: subSystem,
	}

	prom.requestURLMappingFn = urlMapping

	prom.reqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: prom.namespace,
			Subsystem: prom.subsystem,
			Name:      "http_request_total",
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		labelNames,
	)

	prom.reqDurationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: prom.namespace,
			Subsystem: prom.subsystem,
			Name:      "http_request_duration_seconds",
			Help:      "The HTTP request latencies in seconds.",
		},
		labelNames,
	)

	prom.reqSizeSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: prom.namespace,
			Subsystem: prom.subsystem,
			Name:      "http_request_size_bytes",
			Help:      "The HTTP request sizes in bytes.",
		},
		labelNames,
	)

	prom.resSizeSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: prom.namespace,
			Subsystem: prom.subsystem,
			Name:      "http_response_size_bytes",
			Help:      "The HTTP response sizes in bytes.",
		},
		labelNames,
	)

	prometheus.MustRegister(prom.reqCounter)
	prometheus.MustRegister(prom.reqDurationSummary)
	prometheus.MustRegister(prom.reqSizeSummary)
	prometheus.MustRegister(prom.resSizeSummary)

	return prom
}

// SetPushGateway set up push gateway configure
// pushTargetURL is push gateway server url
// pushJobName is job name
// pushInterval is second
func (prom *Prom) SetPushGateway(pushTargetURL, pushJobName string, pushInterval int) *Prom {
	prom.pushInterval = time.Duration(pushInterval) * time.Second
	prom.pushJobName = pushJobName
	prom.pushTargetURL = pushTargetURL

	prom.startPusher()

	return prom
}

// SetRequestURLMappingFn set up url mapping func
// default is:
// func urlMapping(c *gin.Context) string {
// 	url := c.Request.URL.Path
// 	for _, p := range c.Params {
// 		url = strings.Replace(url, "/"+p.Value, "/:"+p.Key, 1)
// 	}
// 	return url
// }
func (prom *Prom) SetRequestURLMappingFn(newFn func(*gin.Context) string) *Prom {
	prom.requestURLMappingFn = newFn

	return prom
}

// Handler is prometheus middleware of gin
func (prom *Prom) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		reqSize := computeRequestSize(ctx.Request)

		ctx.Next()

		end := time.Now()
		elapsed := end.Sub(start)
		status := strconv.Itoa(ctx.Writer.Status())
		resSize := float64(ctx.Writer.Size())
		url := prom.requestURLMappingFn(ctx)
		errcode := strconv.Itoa(ctx.GetInt("errcode"))
		labels := []string{status, errcode, ctx.Request.Method, url, ctx.HandlerName()}

		prom.reqCounter.WithLabelValues(labels...).Inc()
		prom.reqSizeSummary.WithLabelValues(labels...).Observe(reqSize)
		prom.resSizeSummary.WithLabelValues(labels...).Observe(resSize)
		prom.reqDurationSummary.WithLabelValues(labels...).Observe(elapsed.Seconds())
	}
}

// Metrics is prometheus metrics middleware of gin
func (prom *Prom) Metrics() gin.HandlerFunc {
	return gin.WrapH(prometheus.Handler())
}

func (prom *Prom) startPusher() {
	if prom.pushTargetURL == "" || prom.pushJobName == "" {
		return
	}

	hostname, _ := os.Hostname()

	go func() {
		ticker := time.NewTicker(prom.pushInterval)
		defer ticker.Stop()

		pusher := push.
			New(prom.pushTargetURL, prom.pushJobName).
			Gatherer(prometheus.DefaultRegisterer.(prometheus.Gatherer)).
			Grouping("instance", hostname)

		for {
			select {
			case <-ticker.C:
				pusher.Add()
			}
		}
	}()
}
