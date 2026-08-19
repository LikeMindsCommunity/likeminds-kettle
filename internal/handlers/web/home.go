package web

import (
	"net/http"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// Home is used to validate health check and host / path
func Home(c *gin.Context) {
	//Send response with success as true
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
	})
}
