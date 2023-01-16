package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/community"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type PollObject struct {
	Text string `json:"text"`
}

type ConversationPreview struct {
	InternalLink string                    `json:"internal_link"`
	PreviewType  string                    `json:"preview_type"`
	PreviewText  string                    `json:"preview_text"`
	Title        string                    `json:"title"`
	Community    community.CommunityObject `json:"community"`
	Action       string                    `json:"action"`
	ActionRoute  string                    `json:"action_route"`
}

type CreateConversationRequest struct {
	ChatroomID            int64               `json:"chatroom_id"`
	Text                  string              `json:"text"`
	PollType              int32               `json:"poll_type"`
	AllowAddOption        bool                `json:"allow_add_option"`
	ExpiryTime            int64               `json:"expiry_time"`
	Polls                 []PollObject        `json:"polls"`
	AttachmentsCount      int64               `json:"attachments_count"`
	RepliedConversationId int64               `json:"replied_conversation_id,omitempty"`
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

// CreateConversation is used to create a new conversation in chatroom
func CreateConversation(c *gin.Context) {
	Conversation(c, utils.POSTMethod)
}

// EditConversation is used to edit a specific conversation
func EditConversation(c *gin.Context) {
	Conversation(c, utils.PUTMethod)
}

// GetConversation is used to get conversation
func GetConversation(c *gin.Context) {
	Conversation(c, utils.GETMethod)
}

// DeleteConversation is used to delete conversation
func DeleteConversation(c *gin.Context) {
	Conversation(c, utils.DELETEMethod)
}

// Conversation method handles conversation object
func Conversation(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send request
	switch method {
	case utils.GETMethod:

		getConversationInternal(c, userId)

	case utils.POSTMethod:

		createConversationInternal(c, userId)

	case utils.PUTMethod:

		editConversationInternal(c, userId)

	case utils.DELETEMethod:

		deleteConversationInternal(c, userId)
	}
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

func getConversationInternal(c *gin.Context, userId string) {

	//GET Request params
	meta := c.Query(ParamMeta)
	params := map[string]string{
		ParamChatroomId:     c.Query(ParamChatroomId),
		ParamConversationId: c.Query(ParamConversationId),
	}

	if meta == "" || meta == "false" {
		//If meta is missing, call api/conversation/fetch api internally

		//Send Request
		utils.SendRequest(c, utils.CoreService, FetchConversationEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	} else {
		//else, call api/conversation_meta api internally

		//Send Request
		utils.SendRequest(c, utils.CoreService, ConversationMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	}
}

func createConversationInternal(c *gin.Context, userId string) {

	//Body to be sent in the create conversation api internally
	createConversationRequest, err := parseCreateConversationRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, CreateConversationEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createConversationRequest)
}

func editConversationInternal(c *gin.Context, userId string) {

	//Body to be sent in the edit conversation api internally
	editConversationRequest, err := parseEditConversationRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, EditConversationEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, editConversationRequest)
}

func deleteConversationInternal(c *gin.Context, userId string) {

	//Body to be sent in the delete conversation api internally
	deleteConversationRequest, err := parseDeleteConversationRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, DeleteConversationEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, deleteConversationRequest)
}
