package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type CreateProviderMetaRequest struct {
	AssistantId         string `json:"assistant_id" binding:"required"`
	ThreadContext       int    `json:"thread_context"`
	MaxPromptTokens     int    `json:"max_prompt_tokens"`
	MaxCompletionTokens int    `json:"max_completion_tokens"`
	VisionModel         string `json:"vision_model"`
}

type CreateChatbotMetaRequest struct {
	DefaultText  string                    `json:"default_text"`
	Provider     string                    `json:"provider" binding:"required"`
	ProviderMeta CreateProviderMetaRequest `json:"provider_meta" binding:"required"`
}

type CreateChatbotRequest struct {
	UUID        string                   `json:"uuid"`
	ImageUrl    string                   `json:"image_url"`
	Name        string                   `json:"name" binding:"required"`
	ChatbotMeta CreateChatbotMetaRequest `json:"chatbot_meta" binding:"required"`
}

type UpdateChatbotRequest struct {
	ChatbotUUID string                   `json:"chatbot_uuid"`
	ImageUrl    string                   `json:"image_url"`
	Name        string                   `json:"name"`
	ChatbotMeta UpdateChatbotMetaRequest `json:"chatbot_meta"`
}

type UpdateProviderMeta struct {
	AssistantId         string `json:"assistant_id"`
	ThreadContext       int    `json:"thread_context"`
	MaxPromptTokens     int    `json:"max_prompt_tokens"`
	MaxCompletionTokens int    `json:"max_completion_tokens"`
	VisionModel         string `json:"vision_model"`
}

type UpdateChatbotMetaRequest struct {
	DefaultText  string             `json:"default_text"`
	Provider     string             `json:"provider"`
	ProviderMeta UpdateProviderMeta `json:"provider_meta"`
}

func parseChatbotRequest(c *gin.Context) (*CreateChatbotRequest, error) {

	var ccr CreateChatbotRequest
	if err := c.ShouldBindJSON(&ccr); err != nil {
		return nil, err
	}

	return &ccr, nil
}

func parseUpdateChatbotRequest(c *gin.Context) (*UpdateChatbotRequest, error) {

	var ucr UpdateChatbotRequest
	if err := c.ShouldBindJSON(&ucr); err != nil {
		return nil, err
	}

	return &ucr, nil
}

// Expose method to create chatbot in a community
func CreateChatbot(c *gin.Context) {
	chatbot(c, utils.POSTMethod)
}

// External method to fetch chatbots in a community
func FetchChatbots(c *gin.Context) {
	chatbot(c, utils.GETMethod)
}

// External method to update chatbot in a community
func UpdateChatbot(c *gin.Context) {
	chatbot(c, utils.PatchMethod)
}

// Internal method for chatbot routes
func chatbot(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Get bot id if call is from dashboard
	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	switch method {
	case utils.POSTMethod:
		createChatbot(c, userId)
	case utils.GETMethod:
		fetchChatbots(c, userId)
	case utils.PatchMethod:
		updateChatbot(c, userId)
	}
}

// Internal method to create chatbot
func createChatbot(c *gin.Context, userId string) {

	headers := utils.CreateHeaders(c, userId)

	// Parse request
	chatbotReqBody, err := parseChatbotRequest(c)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Send request to api/community/chatbot
	utils.SendRequest(c, utils.CoreService, CommunityChatbotEndpoint, utils.POSTRequestRawBody, headers, nil, chatbotReqBody)
}

// Internal method to fetch chatbots
func fetchChatbots(c *gin.Context, userId string) {

	headers := utils.CreateHeaders(c, userId)

	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
		ParamUUIDs:    c.Query(ParamUUIDs),
	}

	// Send request to api/community/chatbot
	utils.SendRequest(c, utils.CoreService, CommunityChatbotEndpoint, utils.GETRequest, headers, params, nil)
}

// Internal method to update chatbot
func updateChatbot(c *gin.Context, userId string) {

	headers := utils.CreateHeaders(c, userId)
	chatbotUUID := c.Param(ParamChatbotUUID)

	// Parse request
	reqBody, err := parseUpdateChatbotRequest(c)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	reqBody.ChatbotUUID = chatbotUUID

	// Send request to api/community/chatbot
	utils.SendRequest(c, utils.CoreService, CommunityChatbotEndpoint, utils.PATCHRequest, headers, nil, reqBody)
}
