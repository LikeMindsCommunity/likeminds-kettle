package internalServices

import (
	"fmt"

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
		keys, err := cache.GetKeys(redisClient, keyPattern)
		if err != nil {
			utils.GeneralAPIError(c, err.Error())
			return
		}

		if len(keys) == 0 {
			logging.Info(fmt.Sprint("No Cache keys found for pattern: ", keyPattern))
			continue
		}

		// Delete all keys
		for _, key := range keys {
			err = cache.Delete(redisClient, key)
			if err != nil {
				utils.GeneralAPIError(c, err.Error())
				return
			}
			logging.Info(fmt.Sprint("Successfully deleted cache for key: ", key))
		}
	}

	utils.GenerateResponse(c, nil, false)
}
