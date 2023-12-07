package community

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type AcceptRejectConnectionMemberRequest struct {
	Action string `json:"action"`
}

func parseAcceptRejectMemberConnectionRequest(c *gin.Context) (*AcceptRejectConnectionMemberRequest, error) {
	//POST body params
	var arcmr AcceptRejectConnectionMemberRequest

	if err := c.ShouldBindJSON(&arcmr); err != nil {
		return nil, err
	}

	return &arcmr, nil
}

// CreateMemberConnection is used to create a new member connection request
func CreateMemberConnection(c *gin.Context) {
	MemberConnection(c, utils.POSTMethod)
}

// AcceptRejectMemberConnection is used to accept or reject a member connection request
func AcceptRejectMemberConnection(c *gin.Context) {
	MemberConnection(c, utils.PatchMethod)
}

// GetMemberConnection is used to get a member connections
func GetMemberConnection(c *gin.Context) {
	MemberConnection(c, utils.GETMethod)
}

// MemberConnection handles all connection related objects
func MemberConnection(c *gin.Context, method int) {
	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Send request
	switch method {
	case utils.GETMethod:
		getMemberConnectionInternal(c, userId)

	case utils.POSTMethod:
		createMemberConnectionInternal(c, userId)

	case utils.PatchMethod:
		acceptRejectMemberConnectionInternal(c, userId)
	}
}

func createMemberConnectionInternal(c *gin.Context, userId string) {
	// user id received in path params
	paramUserId := c.Param(ParamUserId)
	createMemberConnectionEndPoint := fmt.Sprintf(MemberConnectionEndPoint, paramUserId)

	// Send Request
	utils.SendRequest(c, utils.CoreService, createMemberConnectionEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, nil)
}

func acceptRejectMemberConnectionInternal(c *gin.Context, userId string) {
	// user id received in path params
	paramUserId := c.Param(ParamUserId)
	acceptRejectMemberConnectionEndPoint := fmt.Sprintf(MemberConnectionEndPoint, paramUserId)

	// Body to be sent in the accept reject member conneciton api internally
	acceptRejectMemberConnectionRequest, err := parseAcceptRejectMemberConnectionRequest(c)
	if err != nil {
		// If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, acceptRejectMemberConnectionEndPoint, utils.PATCHRequest, utils.CreateHeaders(c, userId), nil, acceptRejectMemberConnectionRequest)
}

func getMemberConnectionInternal(c *gin.Context, userId string) {
	// user id received in path params
	paramUserId := c.Param(ParamUserId)
	getMemberConnectionEndPoint := fmt.Sprintf(MemberConnectionEndPoint, paramUserId)

	// Params to be sent in the api/chatroom/fetch_all request
	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
		ParamStatus:   c.Query(ParamStatus),
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, getMemberConnectionEndPoint, utils.PATCHRequest, utils.CreateHeaders(c, userId), params, nil)
}
