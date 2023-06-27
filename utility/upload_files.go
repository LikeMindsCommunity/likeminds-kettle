package utility

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type UploadFilesRequest struct {
	ChatroomID     interface{} `json:"chatroom_id"`
	ConversationID interface{} `json:"conversation_id"`
	PollID         interface{} `json:"poll_id,omitempty"`
	Name           string      `json:"name,omitempty"`
	Url            string      `json:"url,omitempty"`
	Type           string      `json:"type,omitempty"`
	FilesCount     int64       `json:"files_count,omitempty"`
	ThumbnailUrl   string      `json:"thumbnail_url,omitempty"`
	Index          int64       `json:"index,omitempty"`
	Height         int64       `json:"height,omitempty"`
	Width          int64       `json:"width,omitempty"`
	Meta           interface{} `json:"meta,omitempty"`
}

func UploadFiles(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	uploadFileRequest, err := parseUploadFilesRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, UploadFilesEndpoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, uploadFileRequest)
}

func parseUploadFilesRequest(c *gin.Context) (*UploadFilesRequest, error) {
	// POST body params
	var ufr UploadFilesRequest

	if err := c.ShouldBindBodyWith(&ufr, binding.JSON); err != nil {
		return nil, err
	}

	if ufr.ChatroomID != nil {
		ufr.ChatroomID = utils.ParseInterfaceToString(ufr.ChatroomID)
	}

	if ufr.ConversationID != nil {
		ufr.ConversationID = utils.ParseInterfaceToString(ufr.ConversationID)
	}

	if ufr.PollID != nil {
		ufr.PollID = utils.ParseInterfaceToString(ufr.PollID)
	}

	return &ufr, nil
}
