package webhook

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateWebhookRequest struct {
	WebhookType string `json:"webhook_type" binding:"required"`
	URL         string `json:"url" binding:"required"`
	IsActive    *bool  `json:"is_active" binding:"required"`
}

func parseCreateWebhookRequest(c *gin.Context) (*CreateWebhookRequest, error) {
	//POST body params
	var cwr CreateWebhookRequest

	if err := c.ShouldBindJSON(&cwr); err != nil {
		return nil, err
	}

	return &cwr, nil
}

type EditWebhookRequest struct {
	URL      string `json:"url,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

func parseEditWebhookRequest(c *gin.Context) (*EditWebhookRequest, error) {
	//POST body params
	var ewr EditWebhookRequest

	if err := c.ShouldBindJSON(&ewr); err != nil {
		return nil, err
	}

	return &ewr, nil
}

func GetWebhooks(c *gin.Context) {
	Webhook(c, utils.GETMethod)
}

func CreateWebhook(c *gin.Context) {
	Webhook(c, utils.POSTMethod)
}

func EditWebhook(c *gin.Context) {
	Webhook(c, utils.PatchMethod)
}

func DeleteWebhook(c *gin.Context) {
	Webhook(c, utils.DELETEMethod)
}

func Webhook(c *gin.Context, method int) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Get bot id if request from dashboard
	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	switch method {

	// Get all webhooks
	case utils.GETMethod:

		// send request to core service
		utils.SendRequest(c, utils.CoreService, WebhooksEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

	// Create a new webhook
	case utils.POSTMethod:

		createWebhookInternal(c, userId)

	// Edit an existing webhook
	case utils.PatchMethod:

		editWebhookInternal(c, userId)

	// Delete an existing webhook
	case utils.DELETEMethod:

		deleteWebhookInternal(c, userId)

	}

}

// Exposed function to get a webhook
func GetWebhook(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Get bot id if request from dashboard
	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	webhookId := c.Param(utils.ParamWebhookId)

	webhookEndpoint := fmt.Sprintf(WebhookEndpoint, webhookId)

	// Send request to core service
	utils.SendRequest(c, utils.CoreService, webhookEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}

func createWebhookInternal(c *gin.Context, userId string) {

	// Parse request
	cwr, err := parseCreateWebhookRequest(c)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Send request to core service
	utils.SendRequest(c, utils.CoreService, WebhooksEndpoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, cwr)
}

func editWebhookInternal(c *gin.Context, userId string) {

	webhookId := c.Param("webhook_id")

	webhookEndpoint := fmt.Sprintf(WebhookEndpoint, webhookId)

	// Parse request
	ewr, err := parseEditWebhookRequest(c)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Send request to core service
	utils.SendRequest(c, utils.CoreService, webhookEndpoint, utils.PATCHRequest, utils.CreateHeaders(c, userId), nil, ewr)
}

func deleteWebhookInternal(c *gin.Context, userId string) {

	webhookId := c.Param("webhook_id")

	webhookEndpoint := fmt.Sprintf(WebhookEndpoint, webhookId)

	// Send request to core service
	utils.SendRequest(c, utils.CoreService, webhookEndpoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, nil)
}
