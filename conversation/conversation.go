package conversation

import (
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

type EditConversationRequest struct {
	ConversationID int64  `json:"conversation_id" binding:"required"`
	Text           string `json:"text" binding:"required"`
	ShareLink      string `json:"share_link,omitempty"`
}

type DeleteConversationRequest struct {
	ConversationIDs []int64 `json:"conversation_ids" binding:"required"`
	TagID           int64   `json:"tag_id"`
	Reason          string  `json:"reason" binding:"required"`
}

//CreateConversation is used to create a new conversation in chatroom
func CreateConversation(c *gin.Context) {
	Conversation(c, utils.POSTMethod)
}

//EditConversation is used to edit a specific conversation
func EditConversation(c *gin.Context) {
	Conversation(c, utils.PUTMethod)
}

//GetConversation is used to get conversation
func GetConversation(c *gin.Context) {
	Conversation(c, utils.GETMethod)
}

//DeleteConversation is used to delete conversation
func DeleteConversation(c *gin.Context) {
	Conversation(c, utils.DELETEMethod)
}

//Conversation method handles conversation object
func Conversation(c *gin.Context, method int) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Send request
	var respBytes []byte
	switch method {
	case utils.GETMethod:

		respBytes = getConversationInternal(c, client, ltm)

	case utils.POSTMethod:

		respBytes = createConversationInternal(c, client, ltm)

	case utils.PUTMethod:

		respBytes = editConversationInternal(c, client, ltm)

	case utils.DELETEMethod:

		respBytes = deleteConversationInternal(c, client, ltm)
	}

	if respBytes == nil {
		return
	}

	//Parse response
	utils.ParseResponse(c, respBytes)
}

func parseCreateConversationRequest(c *gin.Context) (*CreateConversationRequest, error) {
	//POST body params
	var ccr CreateConversationRequest

	if err := c.ShouldBindJSON(&ccr); err != nil {
		return nil, err
	}

	return &ccr, nil
}

func parseEditConversationRequest(c *gin.Context) (*EditConversationRequest, error) {
	//POST body params
	var ecr EditConversationRequest

	if err := c.ShouldBindJSON(&ecr); err != nil {
		return nil, err
	}

	return &ecr, nil
}

func parseDeleteConversationRequest(c *gin.Context) (*DeleteConversationRequest, error) {
	//POST body params
	var dcr DeleteConversationRequest

	if err := c.ShouldBindJSON(&dcr); err != nil {
		return nil, err
	}

	return &dcr, nil
}

func getConversationInternal(c *gin.Context, client *api_client.APIClient, ltm *token.LoginTokenMeta) []byte {
	var options api_client.GetRequestOptions

	//GET Request params
	chatroom_id := c.Query(ParamChatroomId)
	meta := c.Query(ParamMeta)

	if chatroom_id == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return nil
	}

	if meta == "" || meta == "false" {
		//If meta is missing, call api/conversation/fetch api internally

		//Params to be sent in the api/conversation/fetch request
		params := map[string]string{
			ParamChatroomId: chatroom_id,
		}

		options = api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + FetchConversationEndPoint,
			CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
			Params:        params,
		}

	} else {
		//else, call api/conversation_meta api internally

		//GET Request params
		conversationId := c.Query(ParamConversationId)
		if conversationId == "" {
			//If GET params are missing
			utils.GETQueryParamsMissingError(c)
			return nil
		}

		//Params to be sent in the api/conversation_meta request
		params := map[string]string{
			ParamChatroomId:     chatroom_id,
			ParamConversationId: conversationId,
		}

		options = api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + ConversationMetaEndPoint,
			CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
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

func createConversationInternal(c *gin.Context, client *api_client.APIClient, ltm *token.LoginTokenMeta) []byte {
	//Body to be sent in the api/conversation/create POST request
	createConversationRequest, err := parseCreateConversationRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + CreateConversationEndPoint,
		Body:          createConversationRequest,
		CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
	}

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeRaw)

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}

func editConversationInternal(c *gin.Context, client *api_client.APIClient, ltm *token.LoginTokenMeta) []byte {
	//Body to be sent in the api/edit_conversation POST request
	editConversationRequest, err := parseEditConversationRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + EditConversationEndPoint,
		Body:          editConversationRequest,
		CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
	}

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeFormUrlEncoded)

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}

func deleteConversationInternal(c *gin.Context, client *api_client.APIClient, ltm *token.LoginTokenMeta) []byte {
	//Body to be sent in the api/delete_conversation POST request
	deleteConversationRequest, err := parseDeleteConversationRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + DeleteConversationEndPoint,
		Body:          deleteConversationRequest,
		CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
	}

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeRaw)

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}
