package chatroom

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

type Attachment struct {
	URL    string `json:"url"`
	Index  int32  `json:"index"`
	Type   string `json:"type"`
	Height int32  `json:"height"`
	Width  int32  `json:"width"`
}

type UpdateFilesRequest struct {
	ChatroomID  interface{}  `json:"chatroom_id"`
	Attachments []Attachment `json:"attachments"`
}

// UpdateFiles is used to update files in chatroom
func UpdateFiles(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Body to be sent in the api/chatroom/update_files POST request
	updateChatroomFilesRequest, err := parseUpdateChatroomFilesRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, UpdateFilesEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, updateChatroomFilesRequest)
}

func parseUpdateChatroomFilesRequest(c *gin.Context) (*UpdateFilesRequest, error) {
	//POST body params
	var ufr UpdateFilesRequest

	if err := c.ShouldBindJSON(&ufr); err != nil {
		return nil, err
	}

	if ufr.ChatroomID != nil {
		ufr.ChatroomID = utils.ParseInterfaceToString(ufr.ChatroomID)
	}

	return &ufr, nil
}
