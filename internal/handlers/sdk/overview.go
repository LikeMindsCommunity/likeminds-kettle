package sdk

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetMauOverview(c *gin.Context) {

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
