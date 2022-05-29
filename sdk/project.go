package sdk

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// SDKPlatform | SDK platform schema
type SDKPlatform struct {
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

// SDKRequest | create SDK api key request schema
type SDKProjectRequest struct {
	CommunityName  string                   `json:"name" binding:"required"`
	Branding       CommunityBrandingRequest `json:"branding"`
	Headline       string                   `json:"headline"`
	ImageURL       string                   `json:"image_url"`
	Platform       []SDKPlatform            `json:"platform"`
	ProjectCreator string                   `json:"project_creator"`
}

//CreateProject is used to create a new sdk project
func CreateProject(c *gin.Context) {
	Project(c, utils.POSTMethod)
}

//GetProject is used to get an existing sdk project
func GetProject(c *gin.Context) {
	Project(c, utils.GETMethod)
}

//Project method handles community sdk project for each client
func Project(c *gin.Context, method int) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	data := user.FetchBot(utils.CreateHeaders(c))
	var userUniqueId = data.Data.(map[string]interface{})["user"].(map[string]interface{})["user_unique_id"]

	headers := utils.CreateHeaders(c)
	headers[utils.HeadersMemberId] = userUniqueId

	//Send request
	var respBytes []byte
	var err error
	switch method {
	case utils.GETMethod:

		//Params to be sent in the api/sdk/fetch request
		params := map[string]string{
			ParamCommunityCreator: ltm.UserID,
		}

		options := api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + ProjectEndpoint,
			CustomHeaders: headers,
			Params:        params,
		}
		respBytes, err = client.GetRequest(&options)
		if err != nil {
			//If API fails or any other error
			utils.GeneralAPIError(c, err.Error())
			return
		}
	case utils.POSTMethod:

		spr, err := parseProjectRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		spr.ProjectCreator = ltm.UserID

		options := api_client.PostRequestOptions{
			Url:           client.CoreServiceBaseURL + ProjectEndpoint,
			Body:          spr,
			CustomHeaders: headers,
		}

		respBytes, err = client.PostRequest(&options)
		if err != nil {
			//If API fails or any other error
			utils.GeneralAPIError(c, err.Error())
			return
		}

	}

	//Parse response
	var apiCR api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
	if err != nil {
		//Internal unmarshal error
		utils.GeneralAPIError(c, err.Error())
	}
	if !apiCR.Success {
		//If api/sdk/project returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}
	//If flow succeeds
	dataResponse := apiCR.Response
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    dataResponse,
	})
}

func parseProjectRequest(c *gin.Context) (*SDKProjectRequest, error) {
	//POST body params
	var spr SDKProjectRequest
	if err := c.ShouldBindJSON(&spr); err != nil {
		return nil, err
	}
	return &spr, nil
}
