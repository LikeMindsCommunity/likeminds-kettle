package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusService implements UseCase interface
type PrometheusService struct {
	httpHistogram *prometheus.HistogramVec
}

// NewPrometheusService create a new prometheus service
func NewPrometheusService() (*PrometheusService, error) {
	httpHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	s := &PrometheusService{
		httpHistogram: httpHistogram,
	}
	err := prometheus.Register(s.httpHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	return s, nil
}

// SaveHTTP send metrics to server
func (s *PrometheusService) SaveHTTP(h *HTTP) {
	s.httpHistogram.WithLabelValues(h.Handler, h.Method, h.StatusCode).Observe(h.Duration)
}
