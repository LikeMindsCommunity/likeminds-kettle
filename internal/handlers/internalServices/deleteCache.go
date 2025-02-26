package internalServices

import (
	"fmt"
	"github.com/go-redis/redis/v7"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/cache"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type deleteCacheRequest struct {
	KeyPatterns []string `json:"key_patterns" binding:"required"`
}

func parseDeleteCacheRequest(c *gin.Context) (*deleteCacheRequest, error) {

	var dcr deleteCacheRequest
	if err := c.ShouldBindJSON(&dcr); err != nil {
		return &dcr, err
	}

	return &dcr, nil
}

// DeleteCache External method for API to delete cache from key patterns
func DeleteCache(c *gin.Context) {

	// Parse request
	dcr, err := parseDeleteCacheRequest(c)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// get redis client from context
	redisClient := utils.GetRedisClientFromContext(c)
	if redisClient == nil {
		utils.GeneralAPIError(c, "Redis client not found")
		return
	}

	// Call the abstracted function to delete cache keys based on patterns
	err = deleteCacheByKeyPatternsInternal(redisClient, dcr.KeyPatterns)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	utils.GenerateResponse(c, nil, false)
}

// deleteCacheByKeyPatternsInternal Abstracted function to delete keys from Redis based on key patterns.
func deleteCacheByKeyPatternsInternal(redisClient *redis.Client, keyPatterns []string) error {
	// Use DeleteCacheByKeyPatterns to get the map of key patterns and matched keys.
	_, err := cache.DeleteCacheByKeyPatterns(redisClient, keyPatterns)
	if err != nil {
		return err
	}
	return nil
}

func DeleteChatroomCache(c *gin.Context, chatroomID interface{}) {
	// get redis client from context
	redisClient := utils.GetRedisClientFromContext(c)
	if redisClient == nil {
		logging.Error("Redis client not found")
		return
	}

	// Cache key for the chatroom data
	cacheKey := fmt.Sprintf(cache.ChatroomKey, chatroomID)
	// Convert the single cacheKey into a slice of strings
	cacheKeys := []string{cacheKey}

	// Call the abstracted function to delete cache keys based on patterns
	err := deleteCacheByKeyPatternsInternal(redisClient, cacheKeys)
	if err != nil {
		logging.Error(fmt.Sprintf("Error deleting key %s: %v", cacheKey, err))
		return
	}
}
