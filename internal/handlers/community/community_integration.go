package community

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

type UpdateCommunityIntegrationStatusRequest struct {
	StatusType string `json:"status_type"`
	Status     bool   `json:"status"`
}

// GetCommunityIntegration handles the get community integration requests
func GetCommumityIntegrations(c *gin.Context) {
	CommunityIntegration(c, utils.GETMethod)
}

// UpdateCommunityIntegrations handles the update community integration requests
func UpdateCommunityIntegrations(c *gin.Context) {
	CommunityIntegration(c, utils.PUTMethod)
}

// CommunityIntegration handles the community integration requests
func CommunityIntegration(c *gin.Context, method int) {
	// Validate user id
	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	switch method {
	case utils.GETMethod:
		getCommunityIntegrationsInternal(c, userId)
	case utils.PUTMethod:
		updateCommunityIntegrationsInternal(c, userId)
	}
}

// getCommunityIntegrations handles the get community integrations request internally
func getCommunityIntegrationsInternal(c *gin.Context, userId string) {
	// Send Request
	utils.SendRequest(c, utils.CoreService, FetchCommunityIntegrationsEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}

// updateCommunityIntegrations handles the update community integrations request internally
func updateCommunityIntegrationsInternal(c *gin.Context, userId string) {
	// Body to be sent in the edit community integrations api internally
	editCommunityIntegrationsRequest, err := parseEditCommunityIntegrationsRequest(c)
	if err != nil {
		// If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, UpdateCommunityIntegrationsEndpoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, editCommunityIntegrationsRequest)
}

// parseEditCommunityIntegrationsRequest parses the edit community integrations request
func parseEditCommunityIntegrationsRequest(c *gin.Context) (*UpdateCommunityIntegrationStatusRequest, error) {

	// POST body params
	var ucir UpdateCommunityIntegrationStatusRequest

	if err := c.ShouldBindJSON(&ucir); err != nil {
		return nil, err
	}

	return &ucir, nil
}
