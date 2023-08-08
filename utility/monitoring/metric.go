package monitoring

import "time"

// HTTP application
type HTTP struct {
	Handler    string
	Method     string
	StatusCode string
	StartedAt  time.Time
	FinishedAt time.Time
	Duration   float64
}

// newHTTP create a new HTTP app
func newHTTP(handler string, method string) *HTTP {
	return &HTTP{
		Handler: handler,
		Method:  method,
	}
}

// Started start monitoring the app
func (http *HTTP) Started() {
	http.StartedAt = time.Now()
}

// Finished app finished
func (http *HTTP) Finished() {
	http.FinishedAt = time.Now()
	http.Duration = time.Since(http.StartedAt).Seconds()
}

type MetricService interface {
	SaveHTTP(h *HTTP)
}
