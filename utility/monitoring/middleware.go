package monitoring

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// PrometheusMiddleware created middleware to be used to send metrics
func PrometheusMiddleware(metricService MetricService) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Create HTTP metric record
		httpService := newHTTP(c.FullPath(), c.Request.Method)
		httpService.Started()
		c.Next()
		httpService.Finished()
		httpService.StatusCode = strconv.Itoa(c.Writer.Status())
		//Save http metric record in service
		metricService.SaveHTTP(httpService)
	}
}
