package sdk

import (
	"fmt"
	"net/http"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/constants"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/logging"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-redis/redis/v7"
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
	CommunityName           string                   `json:"name"`
	Branding                CommunityBrandingRequest `json:"branding"`
	Headline                string                   `json:"headline"`
	ImageURL                string                   `json:"image_url"`
	FirebaseServerKey       string                   `json:"firebase_server_key"`
	Platform                []Platform               `json:"platform"`
	ProjectCreator          string                   `json:"project_creator"`
	IsJoinFormEnabled       bool                     `json:"is_join_form_enabled"`
	GCPServiceAccountFile   map[string]interface{}   `json:"gcp_service_account_file"`
	PSGCPServiceAccountFile map[string]interface{}   `json:"ps_gcp_service_account_file"`
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
		getProjectInternal(c, userId)

	case utils.POSTMethod:
		createProjectInternal(c, userId)

	case utils.PUTMethod:
		updateProjectInternal(c, userId)

	case utils.DELETEMethod:

		botId := user.GetBotId(c)
		if botId != "" {
			userId = botId
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, ProjectEndpoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, nil)
	}
}

func getProjectInternal(c *gin.Context, userId string) {

	headers := utils.CreateHeaders(c, userId)

	//Params to be sent in the SDK fetch request internally
	params := map[string]string{
		ParamCommunityCreator: userId,
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, ProjectEndpoint, utils.GETRequest, headers, params, nil)
}

func createProjectInternal(c *gin.Context, userId string) {

	//Call POST api/bot to create bot
	response := user.GetBotResponse(c, utils.POSTMethod)
	if response == nil {
		return
	}

	botId := user.GetUserUniqueIDFromResponse(response)

	headers := utils.CreateHeaders(c, botId)

	//Params to be sent in the create sdk project request internally
	projectRequest, err := parseCreateProjectRequest(c, userId)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Create Project request and recieve API key
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, ProjectEndpoint, utils.POSTRequestRawBody, headers, nil, projectRequest)
	if respBytes == nil {
		return
	}

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	projectApiResp := apiCR.Response

	if projectApiResp[constants.ResponseKeyApiKey] == nil {
		utils.GeneralAPIError(c, "Community not created")
		return
	}

	// Create billing plan for community
	err = createCommunityBillingPlan(utils.GetRedisClientFromContext(c), projectApiResp, headers)
	if err != nil {
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Response
	utils.ParseResponse(c, respBytes, statusCode, false)
}

func createCommunityBillingPlan(redis *redis.Client, projectApiResp map[string]interface{}, headers map[string]interface{},
) error {

	// Fetch Community ID from API Key
	communityId, err := utils.FetchCommunityIdFromApiKey(redis, projectApiResp[constants.ResponseKeyApiKey].(string))
	if err != nil {
		return err
	}

	//Send Request to create billing plan for community
	communityBillingRequest := map[string]interface{}{}
	billingEndpoint := utils.BillingPlanEnpoint + "/" + fmt.Sprint(communityId)
	billingPlanRespBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.SubscriptionService, billingEndpoint, utils.POSTRequestRawBody, headers, nil, communityBillingRequest)
	if err != nil || statusCode != http.StatusOK {

		if err == nil {
			err = fmt.Errorf("error creating billing plan")
		}

		logging.Error(fmt.Sprintf("Error creating billing plan, response: %s err: %s ", string(billingPlanRespBytes), err.Error()))

		// Delete the created community
		headers[utils.HeadersApiKey] = projectApiResp[constants.ResponseKeyApiKey].(string)
		resp, statusCode, err := utils.GetRequestResponseWithoutContext(utils.CoreService, ProjectEndpoint, utils.DELETERequest, headers, nil, nil)
		if err != nil || statusCode != http.StatusOK {
			if err == nil {
				err = fmt.Errorf("error deleting community")
			}
			logging.Error(fmt.Sprintf("Error deleting community, response: %s err: %s ", string(resp), err.Error()))
		}

		return fmt.Errorf("Error creating billing plan")
	}

	return nil
}

func updateProjectInternal(c *gin.Context, userId string) {

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	headers := utils.CreateHeaders(c, userId)

	//Params to be sent in the update sdk project request internally
	projectRequest, err := parseUpdateProjectRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, ProjectEndpoint, utils.PUTRequest, headers, nil, projectRequest)
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
