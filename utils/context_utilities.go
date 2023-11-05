package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
)

// Exposed utility method to get redis client from context
func GetRedisClientFromContext(c *gin.Context) *redis.Client {
	redisClient, exists := c.Get(cache.ParamRedisClient)
	if !exists {
		return nil
	}
	return redisClient.(*redis.Client)
}
