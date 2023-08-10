package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricService interface {
	SaveResponseTime(responseTimeMetric *ResponseTimeMetric)
	SaveTotalRequest(totalRequestMetric *TotalRequestMetric)
}

// PrometheusService implements UseCase interface
type PrometheusService struct {
	responseTimeHistogram *prometheus.HistogramVec
	totalRequestCounter   *prometheus.CounterVec
}

// SaveResponseTime update response time
func (s *PrometheusService) SaveResponseTime(responseTimeMetric *ResponseTimeMetric) {
	s.responseTimeHistogram.WithLabelValues(responseTimeMetric.Handler, responseTimeMetric.Method, responseTimeMetric.StatusCode).Observe(responseTimeMetric.Duration)
}

// SaveTotalRequest update total no. of requests
func (s *PrometheusService) SaveTotalRequest(responseTimeMetric *TotalRequestMetric) {
	s.totalRequestCounter.WithLabelValues(responseTimeMetric.Handler, responseTimeMetric.Method, responseTimeMetric.StatusCode).Inc()
}

// NewPrometheusService create a new prometheus service
func NewPrometheusService() (*PrometheusService, error) {
	responseTimeHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the ResponseTimeMetric requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	var totalRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "http",
			Name:      "request_total",
			Help:      "Number of get request",
		},
		[]string{"path"},
	)

	s := &PrometheusService{
		responseTimeHistogram: responseTimeHistogram,
		totalRequestCounter:   totalRequestCounter,
	}
	err := prometheus.Register(s.responseTimeHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	err = prometheus.Register(s.totalRequestCounter)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	return s, nil
}
