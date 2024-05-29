package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// FetchMemberState is used to fetch member state in a community
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

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, MemberStateEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode, true)

}

// Exposed method to fetch member role
func FetchMemberRole(c *gin.Context) {

	// If x-accept-version header is not present, then add v1 as default
	if c.GetHeader(utils.HeadersAcceptVersion) == "" {
		c.Request.Header.Add(utils.HeadersAcceptVersion, utils.ApiRevampV1)
	}

	// Send request to api/members_state
	FetchMemberState(c)
}
