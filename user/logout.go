package user

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/cache"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

//Logout is used to blacklist LTM and RTM tokens
func Logout(c *gin.Context) {
	client, ok := c.MustGet("redis_client").(*redis.Client)
	if !ok {
		c.JSON(http.StatusInternalServerError,
			utils.AuthenticationResponse{
				Success:      false,
				ErrorMessage: "Something went wrong! Please try after sometime",
			})
		return
	}
	ltm, ok := c.MustGet("ltm").(*token.LoginTokenMeta)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: "Something went wrong! Please try after sometime",
		})
		return
	}
	rtm, ok := c.MustGet("rtm").(*token.RefreshTokenMeta)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: "Something went wrong! Please try after sometime",
		})
		return
	}

	//TODO - call api/user/logout and get response
	success := true
	errorMessage := ""
	if !success {
		c.JSON(http.StatusInternalServerError, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: errorMessage,
		})
		return
	}

	//Blacklist token
	err := cache.BlacklistToken(client, ltm, rtm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, utils.AuthenticationResponse{
		Success: true,
	})
}
