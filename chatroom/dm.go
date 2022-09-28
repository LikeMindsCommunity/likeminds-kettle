package chatroom

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type InitiateDMRequest struct {
	ChatroomID       int    `json:"chatroom_id"`
	ChatRequestState int    `json:"chat_request_state"`
	Text             string `json:"text"`
	RequestDM        bool   `json:"request_dm"`
	MemberID         int    `json:"member_id"`
}

func InitiatingDMRequest(c *gin.Context) {
	DMChatroom(c, utils.POSTMethod)
}

func ListDMChatrooms(c *gin.Context) {
	DMChatroom(c, utils.GETMethod)
}

func parseInitiateDMRequest(c *gin.Context) (*InitiateDMRequest, error) {
	//POST body params
	var idmr InitiateDMRequest
	if err := c.ShouldBindJSON(&idmr); err != nil {
		return nil, err
	}

	return &idmr, nil
}

func DMChatroom(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	switch method {
	case utils.GETMethod:

		isDMLimit := c.Query(ParamDMLimit)
		var requestURL string
		var requestParams map[string]string

		if isDMLimit == "true" {
			log.Println(isDMLimit)
			requestURL = RequestDMLimitEndPoint

			// Params to be sent in the api/community_member/request_dm_limit request
			requestParams = map[string]string{
				ParamMemberId: c.Query(ParamMemberId),
			}

		} else {
			requestURL = FetchDMChatroomsEndPoint

			// Params to be sent in the api/community_member/fetch_dm_chatrooms request
			requestParams = map[string]string{
				ParamPage: c.Query(ParamPage),
			}
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, requestURL, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)

	case utils.POSTMethod:

		//Params to be sent in the schedule follow request internally
		initiateDMRequest, err := parseInitiateDMRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		request_dm := initiateDMRequest.RequestDM

		if request_dm {

			//Send Request
			utils.SendRequest(c, utils.CoreService, InitiateDMEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, initiateDMRequest)
		} else {

			//Send Request
			utils.SendRequest(c, utils.CoreService, CreateDMEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, initiateDMRequest)
		}
	}
}
