package monitoring

import "time"

// ResponseTimeMetric to record response time metric
type ResponseTimeMetric struct {
	Handler    string
	Method     string
	StatusCode string
	StartedAt  time.Time
	FinishedAt time.Time
	Duration   float64
}

// NewResponseTimeMetric create a new ResponseTimeMetric
func NewResponseTimeMetric(handler string, method string) *ResponseTimeMetric {
	return &ResponseTimeMetric{
		Handler: handler,
		Method:  method,
	}
}

// Started start recording response time
func (http *ResponseTimeMetric) Started() {
	http.StartedAt = time.Now()
}

// Finished response time recorded
func (http *ResponseTimeMetric) Finished() {
	http.FinishedAt = time.Now()
	http.Duration = time.Since(http.StartedAt).Seconds()
}

// TotalRequestMetric to record total no. of requests
type TotalRequestMetric struct {
	Handler    string
	Method     string
	StatusCode string
}

// NewTotalRequestMetric create a new TotalRequestMetric
func NewTotalRequestMetric(handler string, method string, statusCode string) *TotalRequestMetric {
	return &TotalRequestMetric{
		Handler:    handler,
		Method:     method,
		StatusCode: statusCode,
	}
}
