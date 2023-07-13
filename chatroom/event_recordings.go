package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// Request body params to upload event recordings
type uploadEventRecordingsRequest struct {
	ChatroomId     string `json:"chatroom_id"`
	ConversationId string `json:"conversation_id"`
	Url            string `json:"url"`
	Type           string `json:"type"`
	Index          int64  `json:"index"`
	Width          int64  `json:"width"`
	Height         int64  `json:"height"`
	ThumbnailUrl   string `json:"thumbnail_url"`
	Name           string `json:"name"`
	Meta           string `json:"meta"`
	About          string `json:"about"`
	IsRecording    bool   `json:"is_recording"`
}

// Request body params to delete event recordings
type deleteEventRecordingsRequest struct {
	Id string `json:"id" binding:"required"`
}

// function to parse POST body params for uploadEventRecordingsRequest
func parseUploadEventRecordingsRequest(c *gin.Context) (*uploadEventRecordingsRequest, error) {

	var uer uploadEventRecordingsRequest

	if err := c.ShouldBindJSON(&uer); err != nil {
		return nil, err
	}

	return &uer, nil
}

// function to parse POST body params for deleteEventRecordingsRequest
func parseDeleteEventRecordingsRequest(c *gin.Context) (*deleteEventRecordingsRequest, error) {

	var der deleteEventRecordingsRequest

	if err := c.ShouldBindJSON(&der); err != nil {
		return nil, err
	}

	return &der, nil
}

// Exposed function to upload event recordings
func UploadEventRecordings(c *gin.Context) {
	recordings(c, utils.POSTMethod)
}

// Exposed function to delete event recordings
func DeleteEventRecordings(c *gin.Context) {
	recordings(c, utils.DELETEMethod)
}

// Recordings object to handle both upload and delete event recordings
func recordings(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send request according to method
	switch method {
	case utils.POSTMethod:

		uploadEventRecordingsInternal(c, userId)

	case utils.DELETEMethod:

		deleteEventRecordingsInternal(c, userId)
	}
}

// Internal function to upload event recordings
func uploadEventRecordingsInternal(c *gin.Context, userId string) {

	// Parse request body params
	uploadEventRecordingsRequest, err := parseUploadEventRecordingsRequest(c)

	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	// Send request to /api/chatroom/event/upload_recordings
	utils.SendRequest(c, utils.CoreService, UploadEventRecordingsEndPoint, utils.POSTMethod, utils.CreateHeaders(c, userId), nil, uploadEventRecordingsRequest)
}

// Internal function to delete event recordings
func deleteEventRecordingsInternal(c *gin.Context, userId string) {

	// Parse request body params
	deleteEventRecordingsRequest, err := parseDeleteEventRecordingsRequest(c)

	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	// Send request to /api/chatroom/event/delete_recordings
	utils.SendRequest(c, utils.CoreService, DeleteEventRecordingsEndPoint, utils.POSTMethod, utils.CreateHeaders(c, userId), nil, deleteEventRecordingsRequest)
}
