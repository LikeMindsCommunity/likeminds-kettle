package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

//FetchMemberState is used to fetch member state in a community
func FetchMemberState(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the member state api internally
	params := map[string]string{
		ParamMemberId: userId,
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, MemberStateEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
