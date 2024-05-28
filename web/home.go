package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// Home is used to validate health check and host / path
func Home(c *gin.Context) {
	//Send response with success as true
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
	})
}
