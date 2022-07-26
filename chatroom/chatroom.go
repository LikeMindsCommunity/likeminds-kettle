package chatroom

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateChatroomRequest struct {
	Title                      string      `json:"title" binding:"required"`
	Header                     string      `json:"header"`
	ShareLink                  string      `json:"share_link"`
	AttachmentCount            int64       `json:"attachment_count"`
	PdfCount                   int64       `json:"pdf_count"`
	ImageCount                 int64       `json:"image_count"`
	VideoCount                 int64       `json:"video_count"`
	AudioCount                 int64       `json:"audio_count"`
	Type                       int32       `json:"type"`
	DateTime                   int64       `json:"date_time"`
	EndDate                    int64       `json:"end_date"`
	Duration                   int64       `json:"duration"`
	Location                   string      `json:"location"`
	LocationLat                float64     `json:"location_lat"`
	LocationLong               float64     `json:"location_long"`
	About                      string      `json:"about"`
	DraftID                    int64       `json:"draft_id"`
	InternalLink               string      `json:"internal_link"`
	Preview                    interface{} `json:"preview"`
	CoHosts                    []int64     `json:"co_hosts"`
	CohortIDs                  []int64     `json:"cohort_ids"`
	OnlineLink                 string      `json:"online_link"`
	IsSecret                   bool        `json:"is_secret"`
	ChatroomParticipants       []int64     `json:"chatroom_participants"`
	AutoFollowDone             bool        `json:"auto_follow_done"`
	IncludeMembersLater        bool        `json:"include_members_later"`
	SecretChatroomParticipants []int64     `json:"secret_chatroom_participants"`
	ThirdPartyUniqueID         string      `json:"third_party_unique_id"`
	ScheduleTime               int64       `json:"schedule_time"`
	ScheduleTimeBefore         int64       `json:"schedule_time_before"`
	EndTime                    int64       `json:"end_time"`
	EndTimeAfter               int64       `json:"end_time_after"`
}

type EditChatroomRequest struct {
	ChatroomID       int64  `json:"chatroom_id"`
	Title            string `json:"title"`
	Header           string `json:"header"`
	ChatroomImageUrl string `json:"chatroom_image_url"`
}

//CreateChatroom is used to create a new chatroom
func CreateChatroom(c *gin.Context) {
	Chatroom(c, utils.POSTMethod)
}

//EditChatroom is used to edit community details
func EditChatroom(c *gin.Context) {
	Chatroom(c, utils.PUTMethod)
}

//GetChatroom is used to get chatrooms details
func GetChatroom(c *gin.Context) {
	Chatroom(c, utils.GETMethod)
}

//Project method handles community sdk project for each client
func Chatroom(c *gin.Context, method int) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Call GET api/bot to get bot
	response := user.GetBotResponse(c, utils.GETMethod)
	if response == nil {
		return
	}

	//Send request
	var respBytes []byte
	switch method {
	case utils.GETMethod:

		respBytes = getChatroomInternal(c, client, response)

	case utils.POSTMethod:

		respBytes = createChatroomInternal(c, client, response)

	case utils.PUTMethod:

		respBytes = editChatroomInternal(c, client, response)
	}

	if respBytes == nil {
		return
	}

	//Parse response
	utils.ParseResponse(c, respBytes)
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

func getChatroomInternal(c *gin.Context, client *api_client.APIClient, response *utils.Response) []byte {
	var options api_client.GetRequestOptions

	//GET Request params
	chatroom_id := c.Query(ParamChatroomId)
	if chatroom_id == "" {
		//If chatroom_id is missing, call api/chatroom/fetch_all api internally

		//Params to be sent in the api/chatroom/fetch_all request
		params := map[string]string{
			ParamPage: c.Query(ParamPage),
		}

		options = api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + FetchAllChatroomEndPoint,
			CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
			Params:        params,
		}

	} else {
		//else, call api/chatroom/fetch api internally

		//Params to be sent in the api/chatroom/fetch request
		params := map[string]string{
			ParamChatroomId: c.Query(ParamChatroomId),
		}

		options = api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + FetchChatroomEndPoint,
			CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
			Params:        params,
		}
	}

	respBytes, err := client.GetRequest(&options)
	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}

func createChatroomInternal(c *gin.Context, client *api_client.APIClient, response *utils.Response) []byte {
	//Body to be sent in the api/chatroom/create POST request
	createChatroomRequest, err := parseCreateChatroomRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + CreateChatroomEndPoint,
		Body:          createChatroomRequest,
		CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
	}

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeRaw)

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}

func editChatroomInternal(c *gin.Context, client *api_client.APIClient, response *utils.Response) []byte {
	//Body to be sent in the api/chatroom/edit POST request
	editChatroomRequest, err := parseEditChatroomRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + EditChatroomEndPoint,
		Body:          editChatroomRequest,
		CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
	}

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeRaw)

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}
