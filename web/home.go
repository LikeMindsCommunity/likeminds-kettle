package web

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

// Home is used to validate health check and host / path
func Home(c *gin.Context) {
	//Send response with success as true
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
	})
}
