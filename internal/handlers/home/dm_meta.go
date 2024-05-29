package home

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

func DMHome(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, UserDMHomeEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
