package metrics

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	responseTimeMillisecondsHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "sample_external_url_response_ms",
			Help:    "External url response time latency in milliseconds",
			Buckets: []float64{10, 50, 100, 500, 1000},
		},
		[]string{"url"},
	)

	externalURLUpGuague = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sample_external_url_up",
			Help: "Gauge representing if the url is currently responding with a 200 status code",
		},
		[]string{"url"},
	)
)

// ExternalURLUp sets a gauge for a given `url` label to 1 indicating the external URL is available.
func ExternalURLUp(url string) {
	externalURLUpGuague.With(prometheus.Labels{"url": url}).Set(1)
}

// ExternalURLDown sets a gauge for a given `url` label to 0 indicating the external URL is un-available.
func ExternalURLDown(url string) {
	externalURLUpGuague.With(prometheus.Labels{"url": url}).Set(0)
}

// ExternalURLResponseTime records the ending time for a response from an external URL for a given `url` label in ms.
func ExternalURLResponseTime(url string, startTime time.Time) {
	end := time.Since(startTime)

	responseTimeMillisecondsHistogram.With(prometheus.Labels{"url": url}).Observe(float64(end.Milliseconds()))
}

// StartMetricsServer registers the required metrics and starts a server with a `/metrics` endpoint.
func StartMetricsServer(port int) {

	prometheus.MustRegister(externalURLUpGuague)
	prometheus.MustRegister(responseTimeMillisecondsHistogram)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(formatPort(port), nil))
}

func formatPort(port int) string {
	return fmt.Sprintf(":%d", port)
}
