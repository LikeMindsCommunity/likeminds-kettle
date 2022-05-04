package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"net/http"
)

const ErrorSomethingWentWrong = "Something went wrong! Please try again after sometime"
const ErrorDeviceLoggedOut = "Device logged out! Please login again"
const ErrorInvalidRequest = "Invalid request!"
const ErrorQueryParamsMissing = "Query params missing!"
const ErrorBodyParamsMissing = "Body params missing!"
const ErrorInvalidLTM = "Invalid LTM!"
const ErrorInvalidRTM = "Invalid RTM!"
const ErrorRedisFailed = "Unable to initialize Redis!"

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
