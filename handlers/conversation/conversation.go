package conversation

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/community"
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

type ConversationAttachment struct {
	Name         string      `json:"name,omitempty"`
	Url          string      `json:"url,omitempty"`
	Type         string      `json:"type,omitempty"`
	ThumbnailUrl string      `json:"thumbnail_url,omitempty"`
	Index        int         `json:"index,omitempty"`
	Height       int         `json:"height,omitempty"`
	Width        int         `json:"width,omitempty"`
	Meta         interface{} `json:"meta,omitempty"`
	LocationName string      `json:"location_name,omitempty"`
	LocationLat  int         `json:"location_lat,omitempty"`
	LocationLong int         `json:"location_long,omitempty"`
}

type CreateConversationRequest struct {
	ChatroomID            interface{}              `json:"chatroom_id"`
	Text                  string                   `json:"text"`
	PollType              *int32                   `json:"poll_type,omitempty"`
	AllowAddOption        bool                     `json:"allow_add_option,omitempty"`
	ExpiryTime            int64                    `json:"expiry_time,omitempty"`
	Polls                 []PollObject             `json:"polls,omitempty"`
	MultilpleSelectState  *int64                   `json:"multiple_select_state,omitempty"`
	MultilpleSelectNo     int64                    `json:"multiple_select_no,omitempty"`
	AttachmentCount       int64                    `json:"attachment_count,omitempty"`
	RepliedConversationId interface{}              `json:"replied_conversation_id,omitempty"`
	RepliedChatroomID     string                   `json:"replied_chatroom_id,omitempty"`
	InternalLink          string                   `json:"internal_link,omitempty"`
	Preview               ConversationPreview      `json:"preview,omitempty"`
	IsAnonymous           bool                     `json:"is_anonymous,omitempty"`
	State                 int32                    `json:"state"`
	HasFiles              bool                     `json:"has_files,omitempty"`
	TemporaryID           string                   `json:"temporary_id,omitempty"`
	OGTags                interface{}              `json:"og_tags,omitempty"`
	ShareLink             string                   `json:"share_link,omitempty"`
	Attachments           []ConversationAttachment `json:"attachments,omitempty"`
	Metadata              interface{}              `json:"metadata,omitempty"`
}

type EditConversationRequest struct {
	ConversationID interface{} `json:"conversation_id" binding:"required"`
	Text           string      `json:"text" binding:"required"`
	ShareLink      string      `json:"share_link,omitempty"`
	Metadata       interface{} `json:"metadata,omitempty"`
}

type DeleteConversationRequest struct {
	ConversationIDs []interface{} `json:"conversation_ids" binding:"required"`
	TagID           int64         `json:"tag_id"`
	Reason          string        `json:"reason"`
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

	if ccr.ChatroomID != nil {
		ccr.ChatroomID = utils.ParseInterfaceToString(ccr.ChatroomID)
	}

	if ccr.RepliedConversationId != nil {
		ccr.RepliedConversationId = utils.ParseInterfaceToString(ccr.RepliedConversationId)
	}

	return &ccr, nil
}

func parseEditConversationRequest(c *gin.Context) (*EditConversationRequest, error) {
	//POST body params
	var ecr EditConversationRequest

	if err := c.ShouldBindJSON(&ecr); err != nil {
		return nil, err
	}

	if ecr.Metadata != nil {
		metadataString, _ := json.Marshal(ecr.Metadata)

		if metadataString != nil {
			ecr.Metadata = string(metadataString)
		}
	}

	ecr.ConversationID = utils.ParseInterfaceToString(ecr.ConversationID)

	return &ecr, nil
}

func parseDeleteConversationRequest(c *gin.Context) (*DeleteConversationRequest, error) {
	//POST body params
	var dcr DeleteConversationRequest

	if err := c.ShouldBindJSON(&dcr); err != nil {
		return nil, err
	}

	// parse conversation ids to string
	for i := 0; i < len(dcr.ConversationIDs); i++ {
		dcr.ConversationIDs[i] = utils.ParseInterfaceToString(dcr.ConversationIDs[i])
	}

	return &dcr, nil
}

func getConversationInternal(c *gin.Context, userId string) {

	//GET Request params
	meta := c.Query(ParamMeta)

	if meta == "" || meta == "false" {
		//If meta is missing, call api/conversation/fetch api internally
		params := map[string]string{
			ParamChatroomId:                 c.Query(ParamChatroomId),
			ParamConversationId:             c.Query(ParamConversationId),
			ParamPaginateBy:                 c.Query(ParamPaginateBy),
			ParamScrollDirection:            c.Query(ParamScrollDirection),
			ParamIncludeConversationId:      c.Query(ParamIncludeConversationId),
			ParamExcludedConversationStates: c.Query(ParamExcludedConversationStates),
		}

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchConversationEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, true)

	} else {
		//else, call api/conversation_meta api internally
		params := map[string]string{
			ParamChatroomId:     c.Query(ParamChatroomId),
			ParamConversationId: c.Query(ParamConversationId),
		}
		//Send Request
		utils.SendRequest(c, utils.CoreService, ConversationMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	}
}

func createConversationInternal(c *gin.Context, userId string) {

	//Body to be sent in the create conversation api internally
	createConversationRequest, err := parseCreateConversationRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, CreateConversationEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createConversationRequest)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode, true)
}

func editConversationInternal(c *gin.Context, userId string) {

	//Body to be sent in the edit conversation api internally
	editConversationRequest, err := parseEditConversationRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, EditConversationEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, editConversationRequest)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode, true)
}

func deleteConversationInternal(c *gin.Context, userId string) {

	//Body to be sent in the delete conversation api internally
	deleteConversationRequest, err := parseDeleteConversationRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, DeleteConversationEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, deleteConversationRequest)
}
