package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nateshr/likeminds-authentication/logging"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// Platform | SDK platform schema
type Platform struct {
	Type        int    `json:"type"`
	Package     string `json:"package"`
	Certificate string `json:"certificate"`
}

// CommunityBasicBranding | community basic branding schema
type CommunityBasicBranding struct {
	PrimaryColour string `json:"primary_colour"`
}

// CommunityAdvancedBranding | community advanced branding schema
type CommunityAdvancedBranding struct {
	HeaderColour       string `json:"header_colour"`
	ButtonsIconsColour string `json:"buttons_icons_colour"`
	TextLinksColour    string `json:"text_links_colour"`
}

// CommunityBrandingRequest | create SDK api key platform schema
type CommunityBrandingRequest struct {
	Basic    CommunityBasicBranding    `json:"basic"`
	Advanced CommunityAdvancedBranding `json:"advanced"`
}

// ProjectRequest | create SDK api key request schema
type CreateProjectRequest struct {
	CommunityName     string                   `json:"name" binding:"required"`
	Branding          CommunityBrandingRequest `json:"branding"`
	Headline          string                   `json:"headline"`
	ImageURL          string                   `json:"image_url"`
	FirebaseServerKey string                   `json:"firebase_server_key"`
	Platform          []Platform               `json:"platform"`
	ProjectCreator    string                   `json:"project_creator"`
	IsJoinFormEnabled bool                     `json:"is_join_form_enabled"`
}

type UpdateProjectRequest struct {
	CommunityName     string                   `json:"name"`
	Branding          CommunityBrandingRequest `json:"branding"`
	Headline          string                   `json:"headline"`
	ImageURL          string                   `json:"image_url"`
	FirebaseServerKey string                   `json:"firebase_server_key"`
	Platform          []Platform               `json:"platform"`
	ProjectCreator    string                   `json:"project_creator"`
	IsJoinFormEnabled bool                     `json:"is_join_form_enabled"`
}

// CreateProject is used to create a new sdk project
func CreateProject(c *gin.Context) {
	Project(c, utils.POSTMethod)
}

// EditProject is used to edit an sdk project
func EditProject(c *gin.Context) {
	Project(c, utils.PUTMethod)
}

// GetProject is used to get an existing sdk project
func GetProject(c *gin.Context) {
	Project(c, utils.GETMethod)
}

// DeleteProject is used to delete an existing sdk project
func DeleteProject(c *gin.Context) {
	Project(c, utils.DELETEMethod)
}

// Project method handles community sdk project for each client
func Project(c *gin.Context, method int) {

	//Authorizing User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	switch method {
	case utils.GETMethod:

		//Params to be sent in the SDK fetch request internally
		params := map[string]string{
			ParamCommunityCreator: userId,
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, ProjectEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	case utils.POSTMethod:

		//Call POST api/bot to create bot
		response := user.GetBotResponse(c, utils.POSTMethod)
		if response == nil {
			return
		}

		botId := user.GetUserUniqueIDFromResponse(response)

		//Params to be sent in the create sdk project request internally
		projectRequest, err := parseCreateProjectRequest(c, userId)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		//Send Create Project request and recieve API key
		respBytes, _ := utils.GetRequestResponse(c, utils.CoreService, ProjectEndpoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, botId), nil, projectRequest)
		if respBytes == nil {
			return
		}

		//Parse response to get api key
		projectResponse := make(map[string]interface{})
		err = json.Unmarshal(respBytes, &projectResponse)
		if err != nil {
			utils.GeneralAPIError(c, err.Error())
			return
		}
		// Fetch Community ID from API Key
		redis := utils.GetRedisClientFromContext(c)
		communityId, err := utils.FetchCommunityIdFromApiKey(redis, projectResponse["api_key"].(string))
		if err != nil {
			utils.GeneralAPIError(c, err.Error())
			return
		}
		headers := utils.CreateHeaders(c, "")

		//Send Request to create billing plan for community
		communityBillingRequest := map[string]interface{}{}
		skulkRespBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.SubscriptionService, fmt.Sprintf("/%s/%d", utils.BillingPlanEnpoint, communityId), utils.POSTRequestRawBody, headers, nil, communityBillingRequest)
		if err != nil {
			logging.Error(err.Error())
			return
		}
		skulkResponse := make(map[string]interface{})
		err = json.Unmarshal(skulkRespBytes, &skulkResponse)
		if err != nil {
			logging.Error(err.Error())
			return
		}
		if statusCode == http.StatusInternalServerError || statusCode == http.StatusBadRequest {
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Response{
				Success:      false,
				ErrorMessage: skulkResponse["error_message"].(string),
			})
			return
		}
		//Send Response
		utils.ParseResponse(c, respBytes, statusCode, true)

	case utils.PUTMethod:

		botId := user.GetBotId(c)
		if botId != "" {
			userId = botId
		}

		//Params to be sent in the update sdk project request internally
		projectRequest, err := parseUpdateProjectRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, ProjectEndpoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, projectRequest)

	case utils.DELETEMethod:

		botId := user.GetBotId(c)
		if botId != "" {
			userId = botId
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, ProjectEndpoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, nil)

	}
}

func parseCreateProjectRequest(c *gin.Context, projectCreatorID string) (*CreateProjectRequest, error) {
	//POST body params
	var cpr CreateProjectRequest
	if err := c.ShouldBindBodyWith(&cpr, binding.JSON); err != nil {
		return nil, err
	}

	if len(projectCreatorID) > 0 {
		cpr.ProjectCreator = projectCreatorID
	}

	return &cpr, nil
}

func parseUpdateProjectRequest(c *gin.Context) (*UpdateProjectRequest, error) {
	//POST body params
	var upr UpdateProjectRequest
	if err := c.ShouldBindBodyWith(&upr, binding.JSON); err != nil {
		return nil, err
	}

	return &upr, nil
}
