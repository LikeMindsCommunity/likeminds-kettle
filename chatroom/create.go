package chatroom

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateChatroomRequest struct {
	CommunityID                string  `json:"community_id"`
	Title                      string  `json:"title"`
	AttachmentCount            int64   `json:"attachment_count"`
	ImageCount                 int64   `json:"image_count"`
	PdfCount                   int64   `json:"pdf_count"`
	Type                       int64   `json:"type"`
	DateTime                   int64   `json:"date_time"`
	EndDate                    int64   `json:"end_date"`
	Duration                   int64   `json:"duration"`
	Location                   string  `json:"location"`
	LocationLat                float64 `json:"location_lat"`
	LocationLong               float64 `json:"location_long"`
	About                      string  `json:"about"`
	CohortIDs                  []int64 `json:"cohort_ids"`
	OnlineLink                 string  `json:"online_link"`
	IsSecret                   bool    `json:"is_secret"`
	SecretChatroomParticipants []int64 `json:"secret_chatroom_participants"`
	ThirdPartyUniqueID         string  `json:"third_party_unique_id"`
	ScheduleTime               int64   `json:"schedule_time"`
	ScheduleTimeBefore         int64   `json:"schedule_time_before"`
	EndTime                    int64   `json:"end_time"`
	EndTimeAfter               int64   `json:"end_time_after"`
}

//CreateChatroom is used to create a new chatroom
func CreateChatroom(c *gin.Context) {

	//Check if request has valid login token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Create headers from login token
	headers := make(map[string]interface{})
	headers[utils.HeadersMemberId] = ltm.UserID

	//POST body params
	var ccr CreateChatroomRequest
	if err := c.ShouldBindJSON(&ccr); err != nil {
		//If POST body params are missing
		utils.POSTBodyParamsMissingError(c)
		return
	}

	//Create internal API client
	apiClient := api_client.NewAPIClient()

	//Send request
	respBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + CreateChatroomEndPoint,
		CustomHeaders: headers,
		Body:          ccr,
	})

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Parse response
	var apiCR api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
	if err != nil {
		//Internal unmarshal error
		utils.GeneralAPIError(c, err.Error())
	}

	if !apiCR.Success {
		//If api/chatroom/create returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}

	//Send response with api/chatroom/create response
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    apiCR.Response,
	})
}
