package monitoring

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// PrometheusMiddleware created middleware to be used to send metrics
func PrometheusMiddleware(metricService MetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		responseTimeMetric := newResponseTimeMetric(c.FullPath(), c.Request.Method)
		responseTimeMetric.Started()
		c.Next()
		responseTimeMetric.Finished()
		responseTimeMetric.StatusCode = strconv.Itoa(c.Writer.Status())
		//Save response time
		metricService.SaveResponseTime(responseTimeMetric)
		//Save total no. of request
		metricService.SaveTotalRequest(newTotalRequestMetric(c.FullPath(), c.Request.Method, strconv.Itoa(c.Writer.Status())))
	}
}
