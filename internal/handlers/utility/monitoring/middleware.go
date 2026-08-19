package monitoring

import (
	"strconv"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/logging"
	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware created middleware to be used to send metrics
func PrometheusMiddleware(metricService MetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		responseTimeMetric := NewResponseTimeMetric(c.FullPath(), c.Request.Method)
		responseTimeMetric.Started()

		c.Next()

		statusCode := strconv.Itoa(c.Writer.Status())
		responseTimeMetric.StatusCode = statusCode
		responseTimeMetric.Finished()

		metricService.SaveResponseTime(responseTimeMetric)
		metricService.SaveTotalRequest(NewTotalRequestMetric(c.FullPath(), c.Request.Method, statusCode))
	}
}

// getPrometheusMetricService returns prometheus metrics service
func GetPrometheusMetricService() *PrometheusService {
	prometheusService, err := NewPrometheusService()
	if err != nil {
		logging.Fatal(err.Error())
		return nil
	}
	return prometheusService
}
