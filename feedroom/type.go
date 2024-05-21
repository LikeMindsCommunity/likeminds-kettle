package feedroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type FeedroomTypeRequest struct {
	FeedroomID interface{} `json:"feedroom_id" binding:"required"`
	IsSecret   *bool       `json:"is_secret" binding:"required"`
}

// ChangeFeedroomType is used to change feedroom type
func ChangeFeedroomType(c *gin.Context) {
	FeedroomType(c, utils.PUTMethod)
}

// GetFeedroomTypeStatus is used to get feedroom conversion status
func GetFeedroomTypeStatus(c *gin.Context) {
	FeedroomType(c, utils.GETMethod)
}

// FeedroomType is used to change the type of feedroom
func FeedroomType(c *gin.Context, method int) {

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
			chatroom.ParamChatroomId: c.Query(ParamFeedroomId),
		}

		//Params Validation
		if params[chatroom.ParamChatroomId] == "" {
			//If GET params are missing
			utils.GETQueryParamsMissingError(c)
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, chatroom.ChatroomTypeEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	case utils.PUTMethod:

		//Body to be sent in the api/chatroom/change_type POST request
		feedroomTypeRequest, err := parseFeedroomTypeRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralBadRequestError(c, err.Error())
			return
		}

		chatroomTypeRequest := chatroom.ChatroomTypeRequest{
			ChatroomID: feedroomTypeRequest.FeedroomID,
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, chatroom.ChatroomTypeEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, chatroomTypeRequest)
	}
}

func parseFeedroomTypeRequest(c *gin.Context) (*FeedroomTypeRequest, error) {
	//POST body params
	var ctr FeedroomTypeRequest

	if err := c.ShouldBindJSON(&ctr); err != nil {
		return nil, err
	}

	if ctr.FeedroomID != nil {
		ctr.FeedroomID = utils.ParseInterfaceToString(ctr.FeedroomID)
	}

	return &ctr, nil
}
