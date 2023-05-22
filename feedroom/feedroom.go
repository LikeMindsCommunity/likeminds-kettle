package feedroom

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateFeedroomRequest struct {
	Title                      string        `json:"title" binding:"required"`
	Header                     string        `json:"header"`
	CohortIDs                  []int64       `json:"cohort_ids"`
	IsSecret                   bool          `json:"is_secret"`
	FeedroomParticipants       []interface{} `json:"feedroom_participants"`
	UUIDs                      []string      `json:"uuids"`
	AutoFollowDone             bool          `json:"auto_follow_done"`
	IncludeMembersLater        bool          `json:"include_members_later"`
	SecretFeedroomParticipants []interface{} `json:"secret_feedroom_participants"`
	FeedroomImageUrl           string        `json:"feedroom_image_url"`
	Tag                        string        `json:"tag"`
}

type EditFeedroomRequest struct {
	FeedroomID       interface{} `json:"feedroom_id"`
	Title            string      `json:"title"`
	Header           string      `json:"header"`
	FeedroomImageUrl string      `json:"feedroom_image_url"`
	Tag              string      `json:"tag"`
}

type DeleteFeedroomRequest struct {
	FeedroomID interface{} `json:"feedroom_id" binding:"required"`
}

// CreateFeedroom is used to create a new feedroom
func CreateFeedroom(c *gin.Context) {
	Feedroom(c, utils.POSTMethod)
}

// EditFeedroom is used to edit feedroom details
func EditFeedroom(c *gin.Context) {
	Feedroom(c, utils.PUTMethod)
}

// GetFeedroom is used to get feedrooms details
func GetFeedroom(c *gin.Context) {
	Feedroom(c, utils.GETMethod)
}

// DeleteFeedroom is used to delete an existing feedroom
func DeleteFeedroom(c *gin.Context) {
	Feedroom(c, utils.DELETEMethod)
}

// Feedroom method handles feedroom objects
func Feedroom(c *gin.Context, method int) {

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

		getFeedroomInternal(c, userId)

	case utils.POSTMethod:

		createFeedroomInternal(c, userId)

	case utils.PUTMethod:

		editFeedroomInternal(c, userId)

	case utils.DELETEMethod:

		deleteFeedroomInternal(c, userId)
	}
}

func parseCreateFeedroomRequest(c *gin.Context) (*CreateFeedroomRequest, error) {
	//POST body params
	var ccr CreateFeedroomRequest

	if err := c.ShouldBindJSON(&ccr); err != nil {
		return nil, err
	}

	return &ccr, nil
}

func parseEditFeedroomRequest(c *gin.Context) (*EditFeedroomRequest, error) {
	//POST body params
	var ecr EditFeedroomRequest

	if err := c.ShouldBindJSON(&ecr); err != nil {
		return nil, err
	}

	if ecr.FeedroomID != nil {
		ecr.FeedroomID = utils.ParseInterfaceToString(ecr.FeedroomID)
	}

	return &ecr, nil
}

func parseDeleteFeedroomRequest(c *gin.Context) (*DeleteFeedroomRequest, error) {
	//POST body params
	var dcr DeleteFeedroomRequest

	if err := c.ShouldBindJSON(&dcr); err != nil {
		return nil, err
	}

	return &dcr, nil
}

func getFeedroomInternal(c *gin.Context, userId string) {

	//GET Request params
	feedroom_id := c.Query(ParamFeedroomId)
	if feedroom_id == "" {
		//If feedroom_id is missing, call api/chatroom/fetch_all api internally

		//Params to be sent in the api/chatroom/fetch_all request
		params := map[string]string{
			chatroom.ParamPage:       c.Query(chatroom.ParamPage),
			chatroom.ParamFilterType: fmt.Sprintf("[%d]", FeedChatroomType),
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, chatroom.FetchAllChatroomEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	} else {
		//else, call api/chatroom/fetch api internally

		//Params to be sent in the api/chatroom/fetch request
		params := map[string]string{
			chatroom.ParamChatroomId: feedroom_id,
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, chatroom.FetchChatroomEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	}
}

func createFeedroomInternal(c *gin.Context, userId string) {
	//Body to be sent in the api/chatroom/create POST request
	createFeedroomRequest, err := parseCreateFeedroomRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	createChatroomRequest := chatroom.CreateChatroomRequest{
		Title:                      createFeedroomRequest.Title,
		Header:                     createFeedroomRequest.Header,
		Type:                       FeedChatroomType,
		CohortIDs:                  createFeedroomRequest.CohortIDs,
		IsSecret:                   createFeedroomRequest.IsSecret,
		ChatroomParticipants:       createFeedroomRequest.FeedroomParticipants,
		AutoFollowDone:             createFeedroomRequest.AutoFollowDone,
		IncludeMembersLater:        createFeedroomRequest.IncludeMembersLater,
		SecretChatroomParticipants: createFeedroomRequest.SecretFeedroomParticipants,
		UUIDs:                      createFeedroomRequest.UUIDs,
		ChatroomImageUrl:           createFeedroomRequest.FeedroomImageUrl,
		Tag:                        createFeedroomRequest.Tag,
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, chatroom.CreateChatroomEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createChatroomRequest)

}

func editFeedroomInternal(c *gin.Context, userId string) {
	//Body to be sent in the api/chatroom/edit POST request
	editFeedroomRequest, err := parseEditFeedroomRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	editChatroomRequest := chatroom.EditChatroomRequest{
		ChatroomID:       editFeedroomRequest.FeedroomID,
		Title:            editFeedroomRequest.Title,
		Header:           editFeedroomRequest.Header,
		ChatroomImageUrl: editFeedroomRequest.FeedroomImageUrl,
		Tag:              editFeedroomRequest.Tag,
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, chatroom.EditChatroomEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, editChatroomRequest)

}

func deleteFeedroomInternal(c *gin.Context, userId string) {
	//Body to be sent in the api/chatroom_delete POST request
	deleteFeedroomRequest, err := parseDeleteFeedroomRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	deleteChatroomRequest := chatroom.DeleteChatroomRequest{
		ChatroomID: deleteFeedroomRequest.FeedroomID,
	}

	if reflect.TypeOf(deleteChatroomRequest.ChatroomID).String() == "float64" {
		deleteChatroomRequest.ChatroomID = strconv.Itoa(int(deleteChatroomRequest.ChatroomID.(float64)))
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, chatroom.DeleteChatroomEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, deleteChatroomRequest)

}
