package monitoring

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// PrometheusMiddleware created middleware to be used to send metrics
func PrometheusMiddleware(metricService MetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		responseTimeMetric := newResponseTimeMetric(c.FullPath(), c.Request.Method)
		//Start recording response time
		responseTimeMetric.Started()
		//Increase concurrency
		metricService.IncreaseConcurrentRequest()
		//Serve request
		c.Next()
		statusCode := strconv.Itoa(c.Writer.Status())
		responseTimeMetric.StatusCode = statusCode
		//Finish recording response time
		responseTimeMetric.Finished()
		//Save response time
		metricService.SaveResponseTime(responseTimeMetric)
		//Save total no. of request
		metricService.SaveTotalRequest(newTotalRequestMetric(c.FullPath(), c.Request.Method, statusCode))
		//Decrease concurrency
		metricService.DecreaseConcurrentRequest()
	}
}
