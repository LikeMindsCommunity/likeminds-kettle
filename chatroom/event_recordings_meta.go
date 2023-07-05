package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// Request body params for Uploading Event Recordings Meta
type uploadEventRecordingsMetaRequest struct {
	ChatroomId        string `json:"chatroom_id"`
	ConversationId    string `json:"conversation_id"`
	RecordingUrl      string `json:"recording_url"`
	RecordingUrlTitle string `json:"recording_url_title"`
	AboutRecording    string `json:"about_recording"`
	IsRecording       bool   `json:"is_recording"`
}

// Request body params for Deleting Event Recordings Meta
type deleteEventRecordingsMetaRequest struct {
	Id             string `json:"id" binding:"required"`
	ChatroomId     string `json:"chatroom_id"`
	ConversationId string `json:"conversation_id"`
}

// function to parse POST body params for uploadEventRecordingsMetaRequest
func parseUploadEventRecordingsMetaRequest(c *gin.Context) (*uploadEventRecordingsMetaRequest, error) {

	var uer uploadEventRecordingsMetaRequest

	if err := c.ShouldBindJSON(&uer); err != nil {
		return nil, err
	}

	return &uer, nil
}

// function to parse POST body params for deleteEventRecordingsMetaRequest
func parseDeleteEventRecordingsMetaRequest(c *gin.Context) (*deleteEventRecordingsMetaRequest, error) {

	var der deleteEventRecordingsMetaRequest

	if err := c.ShouldBindJSON(&der); err != nil {
		return nil, err
	}

	return &der, nil
}

// Exposed function to upload event recordings meta
func UploadEventRecordingsMeta(c *gin.Context) {
	recordingsMeta(c, utils.POSTMethod)
}

// Exposed function to delete event recordings meta
func DeleteEventRecordingsMeta(c *gin.Context) {
	recordingsMeta(c, utils.DELETEMethod)
}

// recordings meta object to upload/delete event recordings meta
func recordingsMeta(c *gin.Context, method int) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	switch method {
	case utils.POSTMethod:

		uploadEventRecordingsMetaInternal(c, userId)

	case utils.DELETEMethod:

		deleteEventRecordingsMetaInternal(c, userId)
	}
}

// Internal function to upload event recordings meta
func uploadEventRecordingsMetaInternal(c *gin.Context, userId string) {

	// Parse request body params
	uploadEventRecordingsMetaRequest, err := parseUploadEventRecordingsMetaRequest(c)

	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	// Send request to /api/chatroom/event/upload_recordings_meta
	utils.SendRequest(c, utils.CoreService, UploadEventRecordingsMetaEndPoint, utils.POSTMethod, utils.CreateHeaders(c, userId), nil, uploadEventRecordingsMetaRequest)
}

// Internal function to delete event recordings meta
func deleteEventRecordingsMetaInternal(c *gin.Context, userId string) {

	// Parse request body params
	deleteEventRecordingsMetaRequest, err := parseDeleteEventRecordingsMetaRequest(c)

	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	// Send request to /api/chatroom/event/delete_recordings_meta
	utils.SendRequest(c, utils.CoreService, DeleteEventRecordingsMetaEndPoint, utils.POSTMethod, utils.CreateHeaders(c, userId), nil, deleteEventRecordingsMetaRequest)
}
