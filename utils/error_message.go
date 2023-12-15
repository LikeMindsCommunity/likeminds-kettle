package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
)

const ErrorDeviceLoggedOut = "Device logged out! Please login again"
const ErrorInvalidRequest = "Invalid request!"
const ErrorQueryParamsMissing = "Query params missing!"
const ErrorBodyParamsMissing = "Body params missing!"
const ErrorInvalidLTM = "Invalid LTM!"
const ErrorInvalidRTM = "Invalid RTM!"
const ErrorRedisFailed = "Unable to initialize Redis!"
const ErrorRedisClientFailed = "Unable to get Redis client!"
const ErrorInvalidAPIKey = "Invalid API key!"
const ErrorGuestAccessDenied = "Login required!"
const ErrorMemeberAccessFail = "You are not authorized to perform this operation!"
const ErrorUserCannotDm = "User cannot DM"
const ErrorInvalidChannelType = "Invalid channel type!"
const ErrorFetchingUserData = "Error while fetching user data!"
const ErrorInvalidUserId = "Invalid user_id!"
const ErrorCommunityConfigurationsNotFound = "Community configurations not found!"
const ErrorApiKeyNotFound = "Api key not found!"
const ErrorInvalidNotificationType = "Invalid notification type sent"
const ErrorInvalidNotificationSettings = "Invalid notification settings sent"

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

func GeneralBadRequestError(c *gin.Context, errorMessage string) {
	c.JSON(http.StatusBadRequest, Response{
		Success:      false,
		ErrorMessage: errorMessage,
	})
}

func APIError(c *gin.Context, httpStatus int, response Response) {
	c.JSON(httpStatus, response)
}
