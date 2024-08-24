package sdk

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

func GetMauOverview(c *gin.Context){
	
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	params := map[string]string{
		ParamNoOfMonths: c.Query(ParamNoOfMonths),
	}

	utils.SendRequest(c, utils.CoreService, MauOverviewEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
