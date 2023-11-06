package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/token"
)

// Exposed utility method to get redis client from context
func GetRedisClientFromContext(c *gin.Context) *redis.Client {

	redisClient, exists := c.Get(cache.ParamRedisClient)
	if !exists {
		return nil
	}
	return redisClient.(*redis.Client)
}

// Exposed utility method to return the user_unique_id of user from LTM token
func GetUserIdFromContext(c *gin.Context) string {

	var userUniqueId string = ""

	//Check if request has LTM token or not
	ltm, ok := c.Get(token.ParamLTM)
	if !ok {
		//If token is not available
		return ""
	}

	userUniqueId = ltm.(*token.LoginTokenMeta).UserUniqueID

	return userUniqueId
}

func ParseAndFetchProfileWidgets(c *gin.Context, userId string, dataResponse map[string]interface{}) map[string]interface{} {

	if userId == "" {
		return dataResponse
	}

	// Fetch Profile meta configurations and check if widget are enabled
	profileWidgetsEnabled, _ := ProfileWidgetsEnabled(c, userId)
	if profileWidgetsEnabled {
		// If profile widgets are enabled
		widgetIds := GetWidgetIdsFromDataResponse(dataResponse)

		if len(widgetIds) > 0 {
			widgets, _ := GetWidgetsFromWidgetIds(CreateHeaders(c, userId), widgetIds)

			if dataResponse["widgets"] == nil {
				dataResponse["widgets"] = map[string]interface{}{}
			}

			for _, value := range widgets {
				widgetId := value.ID
				dataResponse["widgets"].(map[string]interface{})[widgetId] = value
			}
		}

	}

	return dataResponse
}
