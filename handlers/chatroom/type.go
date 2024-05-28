package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type ChatroomTypeRequest struct {
	ChatroomID interface{} `json:"chatroom_id" binding:"required"`
	IsSecret   *bool       `json:"is_secret" binding:"required"`
}

// ChangeChatroomType is used to change chatroom type
func ChangeChatroomType(c *gin.Context) {
	ChatroomType(c, utils.PUTMethod)
}

// GetChatroomTypeStatus is used to get chatroom conversion status
func GetChatroomTypeStatus(c *gin.Context) {
	ChatroomType(c, utils.GETMethod)
}

// ChatroomType is used to change the type of chatroom
func ChatroomType(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Send request
	switch method {
	case utils.GETMethod:

		//Params to be sent in the api/chatroom/change_type GET request
		params := map[string]string{
			ParamChatroomId: c.Query(ParamChatroomId),
		}

		//Params Validation
		if params[ParamChatroomId] == "" {
			//If GET params are missing
			utils.GETQueryParamsMissingError(c)
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, ChatroomTypeEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	case utils.PUTMethod:

		//Body to be sent in the api/chatroom/change_type POST request
		chatroomTypeRequest, err := parseChatroomTypeRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralBadRequestError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, ChatroomTypeEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, chatroomTypeRequest)
	}
}

func parseChatroomTypeRequest(c *gin.Context) (*ChatroomTypeRequest, error) {
	//POST body params
	var ctr ChatroomTypeRequest

	if err := c.ShouldBindJSON(&ctr); err != nil {
		return nil, err
	}

	ctr.ChatroomID = utils.ParseInterfaceToString(ctr.ChatroomID)

	return &ctr, nil
}
