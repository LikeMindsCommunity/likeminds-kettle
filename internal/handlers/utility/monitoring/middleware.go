package monitoring

import (
	"github.com/gin-gonic/gin"
	"strconv"
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
