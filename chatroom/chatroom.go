package chatroom

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateChatroomRequest struct {
	Title                      string        `json:"title" binding:"required"`
	Header                     string        `json:"header"`
	ShareLink                  string        `json:"share_link"`
	AttachmentCount            int64         `json:"attachment_count"`
	PdfCount                   int64         `json:"pdf_count"`
	ImageCount                 int64         `json:"image_count"`
	VideoCount                 int64         `json:"video_count"`
	AudioCount                 int64         `json:"audio_count"`
	Type                       int32         `json:"type"`
	DateTime                   int64         `json:"date_time"`
	EndDate                    int64         `json:"end_date"`
	Duration                   int64         `json:"duration"`
	Location                   string        `json:"location"`
	LocationLat                float64       `json:"location_lat"`
	LocationLong               float64       `json:"location_long"`
	About                      string        `json:"about"`
	DraftID                    int64         `json:"draft_id"`
	InternalLink               string        `json:"internal_link"`
	Preview                    interface{}   `json:"preview"`
	CoHosts                    []int64       `json:"co_hosts"`
	CohortIDs                  []int64       `json:"cohort_ids"`
	OnlineLink                 string        `json:"online_link"`
	IsSecret                   bool          `json:"is_secret"`
	ChatroomParticipants       []interface{} `json:"chatroom_participants"`
	AutoFollowDone             bool          `json:"auto_follow_done"`
	IncludeMembersLater        bool          `json:"include_members_later"`
	SecretChatroomParticipants []interface{} `json:"secret_chatroom_participants"`
	ThirdPartyUniqueID         string        `json:"third_party_unique_id"`
	ScheduleTime               int64         `json:"schedule_time"`
	ScheduleTimeBefore         int64         `json:"schedule_time_before"`
	EndTime                    int64         `json:"end_time"`
	EndTimeAfter               int64         `json:"end_time_after"`
	ChatroomImageUrl           string        `json:"chatroom_image_url"`
}

type EditChatroomRequest struct {
	ChatroomID       int64  `json:"chatroom_id"`
	Title            string `json:"title"`
	Header           string `json:"header"`
	ChatroomImageUrl string `json:"chatroom_image_url"`
}

type DeleteChatroomRequest struct {
	ChatroomID interface{} `json:"chatroom_id" binding:"required"`
	TagID      int32       `json:"tag_id"`
	Reason     string      `json:"reason"`
}

// CreateChatroom is used to create a new chatroom
func CreateChatroom(c *gin.Context) {
	Chatroom(c, utils.POSTMethod)
}

// EditChatroom is used to edit chatroom details
func EditChatroom(c *gin.Context) {
	Chatroom(c, utils.PUTMethod)
}

// GetChatroom is used to get chatrooms details
func GetChatroom(c *gin.Context) {
	Chatroom(c, utils.GETMethod)
}

// DeleteChatroom is used to delete an existing chatroom
func DeleteChatroom(c *gin.Context) {
	Chatroom(c, utils.DELETEMethod)
}

// Chatroom method handles chatroom objects
func Chatroom(c *gin.Context, method int) {

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

		getChatroomInternal(c, userId)

	case utils.POSTMethod:

		createChatroomInternal(c, userId)

	case utils.PUTMethod:

		editChatroomInternal(c, userId)

	case utils.DELETEMethod:

		deleteChatroomInternal(c, userId)
	}
}

func parseCreateChatroomRequest(c *gin.Context) (*CreateChatroomRequest, error) {
	//POST body params
	var ccr CreateChatroomRequest

	if err := c.ShouldBindJSON(&ccr); err != nil {
		return nil, err
	}

	return &ccr, nil
}

func parseEditChatroomRequest(c *gin.Context) (*EditChatroomRequest, error) {
	//POST body params
	var ecr EditChatroomRequest

	if err := c.ShouldBindJSON(&ecr); err != nil {
		return nil, err
	}

	return &ecr, nil
}

func parseDeleteChatroomRequest(c *gin.Context) (*DeleteChatroomRequest, error) {
	//POST body params
	var dcr DeleteChatroomRequest

	if err := c.ShouldBindJSON(&dcr); err != nil {
		return nil, err
	}

	return &dcr, nil
}

func getChatroomInternal(c *gin.Context, userId string) {

	//GET Request params
	chatroom_id := c.Query(ParamChatroomId)
	if chatroom_id == "" {
		//If chatroom_id is missing, call api/chatroom/fetch_all api internally

		//Params to be sent in the api/chatroom/fetch_all request
		params := map[string]string{
			ParamPage:         c.Query(ParamPage),
			ParamExcludedType: fmt.Sprintf("[%d]", FeedChatroomType),
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, FetchAllChatroomEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	} else {
		//else, call api/chatroom/fetch api internally

		version := c.GetHeader(utils.HeadersVersionCode)

		if version == "v2" {
			//Params to be sent in the /api/v2/fetch_chatroom request
			params := map[string]string{
				ParamChatroomId: chatroom_id,
				ParamAJ:         c.Query(ParamAJ),
				ParamSourceId:   c.Query(ParamSourceId),
				ParamApiType:    strconv.Itoa(SdkApiType),
			}

			//Send Request
			utils.SendRequest(c, utils.CoreService, FetchChatroomV2EndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		} else {
			//Params to be sent in the api/chatroom/fetch request
			params := map[string]string{
				ParamChatroomId: c.Query(ParamChatroomId),
			}

			//Send Request
			utils.SendRequest(c, utils.CoreService, FetchChatroomEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		}

	}
}

func createChatroomInternal(c *gin.Context, userId string) {
	//Body to be sent in the api/chatroom/create POST request
	createChatroomRequest, err := parseCreateChatroomRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, CreateChatroomEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createChatroomRequest)

}

func editChatroomInternal(c *gin.Context, userId string) {
	//Body to be sent in the api/chatroom/edit POST request
	editChatroomRequest, err := parseEditChatroomRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, EditChatroomEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, editChatroomRequest)

}

func deleteChatroomInternal(c *gin.Context, userId string) {
	//Body to be sent in the api/chatroom_delete POST request
	deleteChatroomRequest, err := parseDeleteChatroomRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	if reflect.TypeOf(deleteChatroomRequest.ChatroomID).String() == "float64" {
		deleteChatroomRequest.ChatroomID = strconv.Itoa(int(deleteChatroomRequest.ChatroomID.(float64)))
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, DeleteChatroomEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, deleteChatroomRequest)

}
