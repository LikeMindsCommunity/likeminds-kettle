package community

import (
	"fmt"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

type CreateConnectionMemberRequest struct {
	ConnectionType              string `json:"connection_type"`
	AutoAcceptConnectionRequest bool   `json:"connection_request_auto_accepted"`
}

type AcceptRejectConnectionMemberRequest struct {
	Action         string `json:"action"`
	ConnectionType string `json:"connection_type"`
}

func parseCreateConnectionMemberRequest(c *gin.Context) (*CreateConnectionMemberRequest, error) {
	//POST body params
	var ccmr CreateConnectionMemberRequest

	if err := c.ShouldBindJSON(&ccmr); err != nil {
		return nil, err
	}

	return &ccmr, nil
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
	// user id received in path params
	paramUUID := c.Param(ParamUserId)
	memberConnectionEndPoint := fmt.Sprintf(MemberConnectionEndPoint, paramUUID)

	// Send request
	switch method {
	case utils.GETMethod:
		getMemberConnectionInternal(c, userId, memberConnectionEndPoint)

	case utils.POSTMethod:
		createMemberConnectionInternal(c, userId, memberConnectionEndPoint)

	case utils.PatchMethod:
		acceptRejectMemberConnectionInternal(c, userId, memberConnectionEndPoint)
	}
}

func getMemberConnectionInternal(c *gin.Context, userId string, getMemberConnectionEndPoint string) {
	// Params to be sent in the api/chatroom/fetch_all request
	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
		ParamStatus:   c.Query(ParamStatus),
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, getMemberConnectionEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}

func createMemberConnectionInternal(c *gin.Context, userId string, createMemberConnectionEndPoint string) {
	// Body to be sent in the create connection member api internally
	createConnectionMemberRequest, err := parseCreateConnectionMemberRequest(c)
	if err != nil {
		// If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, createMemberConnectionEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createConnectionMemberRequest)
}

func acceptRejectMemberConnectionInternal(c *gin.Context, userId string, acceptRejectMemberConnectionEndPoint string) {

	// Body to be sent in the accept reject member conneciton api internally
	acceptRejectMemberConnectionRequest, err := parseAcceptRejectMemberConnectionRequest(c)
	if err != nil {
		// If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, acceptRejectMemberConnectionEndPoint, utils.PATCHRequest, utils.CreateHeaders(c, userId), nil, acceptRejectMemberConnectionRequest)
}
