package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/utils/api_client"
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
		return userUniqueId
	}

	userUniqueId = ltm.(*constants.LoginTokenMeta).UserUniqueID

	return userUniqueId
}

func GeneralAPIError(c *gin.Context, errorMessage string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success:      false,
		ErrorMessage: errorMessage,
	})
}

func APIClientError(c *gin.Context, apiCTR api_client.APIClientResponse) {
	c.JSON(http.StatusInternalServerError, apiCTR)
}

func GETQueryParamsMissingError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Success:      false,
		ErrorMessage: ErrorQueryParamsMissing,
	})
}

func POSTBodyParamsMissingError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Success:      false,
		ErrorMessage: ErrorBodyParamsMissing,
	})
}

func TokenAuthError(c *gin.Context, errorMessage string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success:      false,
		ErrorMessage: errorMessage,
	})
}

func MemberAccessFailError(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Success:      false,
		ErrorMessage: ErrorMemeberAccessFail,
	})
}

func AccessForbiddenError(c *gin.Context, errorMessage string) {
	c.JSON(http.StatusForbidden, Response{
		Success:      false,
		ErrorMessage: errorMessage,
	})
}

func GeneralBadRequestError(c *gin.Context, errorMessage string) {
	c.JSON(http.StatusBadRequest, Response{
		Success:      false,
		ErrorMessage: errorMessage,
	})
}

func APIError(c *gin.Context, httpStatus int, response Response) {
	c.JSON(httpStatus, response)
}

func RateLimitError(c *gin.Context, errMessage string) {
	response := Response{
		Success:      false,
		ErrorMessage: errMessage,
	}
	c.AbortWithStatusJSON(http.StatusTooManyRequests, response)
}
