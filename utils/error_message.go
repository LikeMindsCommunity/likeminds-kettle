package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"net/http"
)

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
		ErrorMessage: "Query params missing!",
	})
}

func POSTBodyParamsMissingError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Success:      false,
		ErrorMessage: "Body params missing!",
	})
}

func SomethingWentWrongError(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{
		Success:      false,
		ErrorMessage: "Something went wrong! Please try after sometime",
	})
}

