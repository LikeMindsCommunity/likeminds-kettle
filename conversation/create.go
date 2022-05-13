package conversation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type PollObject struct {
	Text string `json:"text"`
}

type CommunityObject struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Purpose        string `json:"purpose"`
	ImageUrl       string `json:"image_url"`
	ImageUrlRound  string `json:"image_url_round"`
	CreatedBy      string `json:"created_by"`
	PromotersCount int32  `json:"promoters_count"`
	MembersCount   int32  `json:"members_count"`
	MemberState    int32  `json:"member_state"`
}

type ConversationPreview struct {
	InternalLink string          `json:"internal_link"`
	PreviewType  string          `json:"preview_type"`
	PreviewText  string          `json:"preview_text"`
	Title        string          `json:"title"`
	Community    CommunityObject `json:"community"`
	Action       string          `json:"action"`
	ActionRoute  string          `json:"action_route"`
}

type CreateConversationRequest struct {
	ChatroomID            string              `json:"chatroom_id"`
	Text                  string              `json:"reply"`
	PollType              int32               `json:"poll_type"`
	AllowAddOption        bool                `json:"allow_add_option"`
	ExpiryTime            int64               `json:"expiry_time"`
	Polls                 []PollObject        `json:"polls"`
	AttachmentsCount      int64               `json:"attachments_count"`
	RepliedConversationId int64               `json:"replied_conversation_id"`
	InternalLink          string              `json:"internal_link"`
	Preview               ConversationPreview `json:"preview"`
	IsAnonymous           bool                `json:"is_anonymous"`
	State                 int32               `json:"state"`
}

//CreateConversation is used to create new conversation
func CreateConversation(c *gin.Context) {

	//Check if request has LTM token or not
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
	var ccr CreateConversationRequest
	if err := c.ShouldBindJSON(&ccr); err != nil {
		//If POST body params are missing
		utils.POSTBodyParamsMissingError(c)
		return
	}

	//Create internal API client
	apiClient := api_client.NewAPIClient()

	//Send request
	respBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + CreateConversationEndPoint,
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
		//If api/conversation/create returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}

	//Send response with api/conversation/create response
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    apiCR.Response,
	})
}
