package internalServices

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"strings"

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
	err = DeleteCacheByKeyPatterns(redisClient, dcr.KeyPatterns)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	utils.GenerateResponse(c, nil, false)
}

// DeleteCacheByKeyPatterns Abstracted function to delete keys from Redis based on key patterns
func DeleteCacheByKeyPatterns(redisClient *redis.Client, keyPatterns []string) error {
	for _, keyPattern := range keyPatterns {
		// Get all keys matching the pattern
		keys, err := cache.GetKeys(redisClient, keyPattern)
		if err != nil {
			logging.Error(fmt.Sprintf("Error fetching keys for pattern %s: %v", keyPattern, err))
			return err
		}

		if len(keys) == 0 {
			logging.Info(fmt.Sprintf("No cache keys found for pattern: %s", keyPattern))
		} else {
			// Delete all keys
			for _, key := range keys {
				err = cache.Delete(redisClient, key)
				if err != nil {
					logging.Error(fmt.Sprintf("Error deleting key %s: %v", key, err))
					return err
				}
				logging.Info(fmt.Sprintf("Successfully deleted cache for key: %s", key))
			}
		}

		// If deleting participants key, also delete the total participants key
		if isChatroomParticipantKey(keyPattern) {
			//If not keys found against ChatroomParticipantsKey. This can happen since we are saving in ChatroomParticipantsKey only in case of secret chatroom / DM
			if len(keys) == 0 {
				keys = append(keys, keyPattern)
			}
			for _, key := range keys {
				chatroomID := extractChatroomIDFromKey(key) // function to extract chatroomID
				err = deleteChatroomTotalParticipantsCache(redisClient, chatroomID)
				if err != nil {
					logging.Error(fmt.Sprintf("Error deleting total participants key for chatroom %s: %v", chatroomID, err))
					return err
				}
			}
		}
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
	err := DeleteCacheByKeyPatterns(redisClient, cacheKeys)
	if err != nil {
		logging.Error(fmt.Sprintf("Error deleting key %s: %v", cacheKey, err))
		return
	}
}

// deleteChatroomTotalParticipantsCache deletes the chatroom_total_participants_%s key from Redis
func deleteChatroomTotalParticipantsCache(redisClient *redis.Client, chatroomID interface{}) error {
	// Cache key for the total participants data
	totalParticipantsCacheKey := fmt.Sprintf(cache.ChatroomTotalParticipantsKey, chatroomID)

	// Delete the cache key
	err := cache.Delete(redisClient, totalParticipantsCacheKey)
	if err != nil {
		logging.Error(fmt.Sprintf("Error deleting chatroom total participants cache for key %s: %v", totalParticipantsCacheKey, err))
		return err
	}
	logging.Info(fmt.Sprintf("Successfully deleted chatroom total participants cache for key: %s", totalParticipantsCacheKey))
	return nil
}

// Helper function to extract chatroomID from cache key
func extractChatroomIDFromKey(key string) string {
	// Assuming the key format is "chatroom_participants_<chatroom_id>"
	var chatroomID string
	_, err := fmt.Sscanf(key, cache.ChatroomParticipantsKey, &chatroomID)
	if err != nil {
		return ""
	}
	return chatroomID
}

// Check if the key matches the expected prefix
func isChatroomParticipantKey(key string) bool {
	return strings.HasPrefix(key, cache.ChatroomParticipantsPrefix)
}
