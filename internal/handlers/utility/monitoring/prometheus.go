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
func (s *PrometheusService) SaveTotalRequest(totalRequestMetric *TotalRequestMetric) {
	s.totalRequestCounter.WithLabelValues(totalRequestMetric.Handler, totalRequestMetric.Method, totalRequestMetric.StatusCode).Inc()
}

// NewPrometheusService create a new prometheus service
func NewPrometheusService() (*PrometheusService, error) {
	responseTimeHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "kettle",
		Name:      "response_time",
		Help:      "The latency of requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	totalRequestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "kettle",
		Name:      "request",
		Help:      "Total no. of requests",
	},
		[]string{"handler", "method", "code"},
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
