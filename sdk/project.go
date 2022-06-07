package sdk

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

//Platform | SDK platform schema
type Platform struct {
	Type        int    `json:"type"`
	Package     string `json:"package"`
	Certificate string `json:"certificate"`
}

//CommunityBasicBranding | community basic branding schema
type CommunityBasicBranding struct {
	PrimaryColour string `json:"primary_colour"`
}

//CommunityAdvancedBranding | community advanced branding schema
type CommunityAdvancedBranding struct {
	HeaderColour       string `json:"header_colour"`
	ButtonsIconsColour string `json:"buttons_icons_colour"`
	TextLinksColour    string `json:"text_links_colour"`
}

//CommunityBrandingRequest | create SDK api key platform schema
type CommunityBrandingRequest struct {
	Basic    CommunityBasicBranding    `json:"basic"`
	Advanced CommunityAdvancedBranding `json:"advanced"`
}

//ProjectRequest | create SDK api key request schema
type CreateProjectRequest struct {
	CommunityName  string                   `json:"name" binding:"required"`
	Branding       CommunityBrandingRequest `json:"branding"`
	Headline       string                   `json:"headline"`
	ImageURL       string                   `json:"image_url"`
	Platform       []Platform               `json:"platform"`
	ProjectCreator string                   `json:"project_creator"`
}

type UpdateProjectRequest struct {
	CommunityName  string                   `json:"name"`
	Branding       CommunityBrandingRequest `json:"branding"`
	Headline       string                   `json:"headline"`
	ImageURL       string                   `json:"image_url"`
	Platform       []Platform               `json:"platform"`
	ProjectCreator string                   `json:"project_creator"`
}

//CreateProject is used to create a new sdk project
func CreateProject(c *gin.Context) {
	Project(c, utils.POSTMethod)
}

//EditProject is used to edit an sdk project
func EditProject(c *gin.Context) {
	Project(c, utils.PUTMethod)
}

//GetProject is used to get an existing sdk project
func GetProject(c *gin.Context) {
	Project(c, utils.GETMethod)
}

//DeleteProject is used to delete an existing sdk project
func DeleteProject(c *gin.Context) {
	Project(c, utils.DELETEMethod)
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

	//Send request
	var respBytes []byte
	var err error
	switch method {
	case utils.GETMethod:
		//Params to be sent in the api/sdk/project GET request
		params := map[string]string{
			ParamCommunityCreator: ltm.UserUniqueID,
		}
		options := api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + ProjectEndpoint,
			CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
			Params:        params,
		}
		respBytes, err = client.GetRequest(&options)
	case utils.POSTMethod:
		//Call POST api/bot to create bot
		response := user.GetBotResponse(c, utils.POSTMethod)
		if response == nil {
			return
		}
		//Params to be sent in the api/sdk/project POST request
		projectRequest, err := parseCreateProjectRequest(c, ltm.UserUniqueID)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}
		options := api_client.PostRequestOptions{
			Url:           client.CoreServiceBaseURL + ProjectEndpoint,
			Body:          projectRequest,
			CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
		}
		respBytes, err = client.PostRequest(&options)
	case utils.PUTMethod:
		//Call POST api/bot to create bot
		response := user.GetBotResponse(c, utils.GETMethod)
		if response == nil {
			return
		}
		//Params to be sent in the api/sdk/project PUT request
		projectRequest, err := parseUpdateProjectRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}
		options := api_client.PostRequestOptions{
			Url:           client.CoreServiceBaseURL + ProjectEndpoint,
			Body:          projectRequest,
			CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
		}
		respBytes, err = client.PutRequest(&options)
	case utils.DELETEMethod:
		//Call GET api/bot to get bot
		response := user.GetBotResponse(c, utils.GETMethod)
		if response == nil {
			return
		}
		//api/sdk/project with delete request
		options := api_client.PostRequestOptions{
			Url:           client.CoreServiceBaseURL + ProjectEndpoint,
			CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
		}
		respBytes, err = client.DeleteRequest(&options)
	}

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
		return
	}
	if !apiCR.Success {
		//If api/sdk/project returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}
	//If flow succeeds
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    apiCR.Response,
	})
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
