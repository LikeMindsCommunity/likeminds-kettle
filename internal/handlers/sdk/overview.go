package sdk

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// MauOverviewEndPoint | togther service mau overview endpoint
const MauOverviewEndPoint = "/api/sdk/mau_overview"

func GetMauOverview(c *gin.Context){
	userId := user.GetRequestingUserId(c)

	params := map[string]string{
		ParamNoOfMonths: c.Query(ParamNoOfMonths),
	}

	utils.SendRequest(c, utils.CoreService, MauOverviewEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
