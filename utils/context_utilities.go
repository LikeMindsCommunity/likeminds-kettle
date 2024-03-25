package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/constants"
)

// Exposed utility method to get redis client from context
func GetRedisClientFromContext(c *gin.Context) *redis.Client {

	redisClient, exists := c.Get(constants.ParamRedisClient)
	if !exists {
		return nil
	}
	return redisClient.(*redis.Client)
}

// Exposed utility method to return the user_unique_id of user from LTM token
func GetUserIdFromContext(c *gin.Context) string {

	var userUniqueId string = ""

	// Check if request has LTM token or not
	ltm, ok := c.Get(constants.ParamLTM)
	if !ok {
		// If token is not available
		return ""
	}

	userUniqueId = ltm.(*constants.LoginTokenMeta).UserUniqueID

	return userUniqueId
}
