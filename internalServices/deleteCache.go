package internalServices

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/utils"
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

// External method for API to delete cache from key patterns
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

	// Delete all keys matching the pattern
	for _, keyPattern := range dcr.KeyPatterns {

		// Get all keys matching the pattern
		keys, err := cache.Keys(redisClient, keyPattern)
		if err != nil {
			utils.GeneralAPIError(c, err.Error())
			return
		}

		// Delete all keys
		for _, key := range keys {
			err = cache.Delete(redisClient, key)
			if err != nil {
				utils.GeneralAPIError(c, err.Error())
				return
			}
		}
	}

	utils.GenerateResponse(c, nil, false)
}
