package chatroom

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ChatroomMarkReadRequest struct {
	ChatroomID interface{} `json:"chatroom_id" binding:"required"`
}

func ChatroomMarkRead(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	chatroomMarkReadRequest, err := parseChatroomMarkReadRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, ChatroomMarkReadEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, chatroomMarkReadRequest)

}

func parseChatroomMarkReadRequest(c *gin.Context) (*ChatroomMarkReadRequest, error) {
	// POST body params
	var cmr ChatroomMarkReadRequest

	if err := c.ShouldBindBodyWith(&cmr, binding.JSON); err != nil {
		return nil, err
	}

	if cmr.ChatroomID != nil {
		cmr.ChatroomID = utils.ParseInterfaceToString(cmr.ChatroomID)
	}

	return &cmr, nil
}
