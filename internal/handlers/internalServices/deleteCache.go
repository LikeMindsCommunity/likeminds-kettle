package internalServices

import (
	"fmt"
	"strings"

	"github.com/go-redis/redis/v7"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/cache"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/logging"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
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
	result, err := cache.DeleteCacheByKeyPatterns(redisClient, keyPatterns)
	if err != nil {
		return err
	}

	//To delete chatroom_total_participants_<> key
	utils.SafeGo(func() {
		// Loop through the result map.
		for keyPattern, keys := range result {
			// If the key pattern matches chatroom participants, handle total participants cache deletion.
			if isChatroomParticipantKey(keyPattern) {
				if len(keys) == 0 {
					// If no keys were found => chatroom_participants_<> doesn't exist then in that case also we need to delete chatroom_total_participants_<>
					keys = append(keys, keyPattern)
				}

				for _, key := range keys {
					chatroomID := extractChatroomIDFromKey(key) // Extract the chatroom ID from the key.
					err = deleteChatroomTotalParticipantsCache(redisClient, chatroomID)
					if err != nil {
						logging.Error(fmt.Sprintf("Error deleting total participants key for chatroom %s: %v", chatroomID, err))
					}
				}
			}
		}
	})

	return nil
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
