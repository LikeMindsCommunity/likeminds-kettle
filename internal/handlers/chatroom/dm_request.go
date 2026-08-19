package chatroom

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

type InitiateDMRequest struct {
	ChatroomID                  interface{}            `json:"chatroom_id"`
	ChatRequestState            int                    `json:"chat_request_state"`
	Text                        string                 `json:"text"`
	MemberID                    interface{}            `json:"member_id"`
	ShouldStreamChatbotResponse bool                   `json:"should_stream_chatbot_response,omitempty"`
	Metadata                    map[string]interface{} `json:"metadata,omitempty"`
	TemporaryID                 string                 `json:"temporary_id,omitempty"`
}

func parseInitiateDMRequest(c *gin.Context) (*InitiateDMRequest, error) {
	//POST body params
	var idmr InitiateDMRequest
	if err := c.ShouldBindJSON(&idmr); err != nil {
		return nil, err
	}

	if idmr.ChatroomID != nil {
		idmr.ChatroomID = utils.ParseInterfaceToString(idmr.ChatroomID)
	}

	if idmr.MemberID != nil {
		idmr.MemberID = utils.ParseInterfaceToString(idmr.MemberID)
	}

	return &idmr, nil
}

func InitiatingDMRequest(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the intiate dm request internally
	initiateDMRequest, err := parseInitiateDMRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Get request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, InitiateDMEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, initiateDMRequest)
	if respBytes == nil {
		return
	}

	utils.ParseResponse(c, respBytes, statusCode, true)
}
