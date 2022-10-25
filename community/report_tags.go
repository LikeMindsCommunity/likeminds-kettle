package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// GetReportTags is used to fetch report tags in a community
func GetReportTags(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, FetchReportTagsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
