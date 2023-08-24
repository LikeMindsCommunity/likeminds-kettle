package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricService interface {
	SaveResponseTime(responseTimeMetric *ResponseTimeMetric)
	SaveTotalRequest(totalRequestMetric *TotalRequestMetric)
	IncreaseConcurrentRequest()
	DecreaseConcurrentRequest()
}

// PrometheusService implements UseCase interface
type PrometheusService struct {
	responseTimeHistogram  *prometheus.HistogramVec
	totalRequestCounter    *prometheus.CounterVec
	concurrentRequestGauge prometheus.Gauge
}

// SaveResponseTime update response time
func (s *PrometheusService) SaveResponseTime(responseTimeMetric *ResponseTimeMetric) {
	s.responseTimeHistogram.WithLabelValues(responseTimeMetric.Handler, responseTimeMetric.Method, responseTimeMetric.StatusCode).Observe(responseTimeMetric.Duration)
}

// SaveTotalRequest update total no. of requests
func (s *PrometheusService) SaveTotalRequest(totalRequestMetric *TotalRequestMetric) {
	s.totalRequestCounter.WithLabelValues(totalRequestMetric.Handler, totalRequestMetric.Method, totalRequestMetric.StatusCode).Inc()
}

// IncreaseConcurrentRequest update concurrent request
func (s *PrometheusService) IncreaseConcurrentRequest() {
	s.concurrentRequestGauge.Inc()
}

// DecreaseConcurrentRequest update concurrent request
func (s *PrometheusService) DecreaseConcurrentRequest() {
	s.concurrentRequestGauge.Dec()
}

// NewPrometheusService create a new prometheus service
func NewPrometheusService() (*PrometheusService, error) {
	responseTimeHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "kettle",
		Name:      "response_time",
		Help:      "The latency of requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	totalRequestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kettle",
			Name:      "request",
			Help:      "Total no. of requests",
		},
		[]string{"handler", "method", "code"},
	)
	concurrentRequestGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "kettle",
		Name:      "concurrent",
		Help:      "No. of concurrent requests"},
	)

	s := &PrometheusService{
		responseTimeHistogram:  responseTimeHistogram,
		totalRequestCounter:    totalRequestCounter,
		concurrentRequestGauge: concurrentRequestGauge,
	}
	err := prometheus.Register(s.responseTimeHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	err = prometheus.Register(s.totalRequestCounter)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	err = prometheus.Register(s.concurrentRequestGauge)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	return s, nil
}
