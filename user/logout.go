package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/token"
	"net/http"
)

func Logout(c *gin.Context) {
	client, ok := c.MustGet("redis_client").(*redis.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError, "Something went wrong! Please try after sometime")
		return
	}
	ltm, err := token.ExtractLoginTokenMeta(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	//TODO - call api/user/logout and get response
	success := true
	errorMessage := ""
	if !success {
		c.JSON(http.StatusInternalServerError, errorMessage)
		return
	}

	//Blacklist token
	err = cache.BlacklistToken(client, ltm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, "ok")
}
