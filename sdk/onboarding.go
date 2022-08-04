package sdk

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

//CreateScreenRequest | create Onboarding Screen Schema
type CreateScreenRequest struct {
	Index     int64  `json:"index" binding:"required"`
	Image     string `json:"image" binding:"required"`
	Heading   string `json:"heading"`
	Text      string `json:"text"`
	CtaColour string `json:"cta_colour"`
	CtaText   string `json:"cta_text"`
}

type UpdateScreenRequest struct {
	Id        int64  `json:"id" binding:"required"`
	Index     int64  `json:"index"`
	Image     string `json:"image"`
	Heading   string `json:"heading"`
	Text      string `json:"text"`
	CtaColour string `json:"cta_colour"`
	CtaText   string `json:"cta_text"`
}

type DeleteScreenRequest struct {
	Id int64 `json:"id" binding:"required"`
}

//CreateScreen is used to create a new onboarding screen
func CreateScreen(c *gin.Context) {
	OnboardingScreen(c, utils.POSTMethod)
}

//EditScreen is used to edit an existing onboarding screen
func EditScreen(c *gin.Context) {
	OnboardingScreen(c, utils.PUTMethod)
}

//GetScreen is used to get an existing onboarding screen
func GetScreen(c *gin.Context) {
	OnboardingScreen(c, utils.GETMethod)
}

//DeleteScreen is used to delete an existing onboarding screen
func DeleteScreen(c *gin.Context) {
	OnboardingScreen(c, utils.DELETEMethod)
}

//OnboardingScreen method handles onboarding screens for each client project
func OnboardingScreen(c *gin.Context, method int) {

	//Authorizing User
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

		//Params to be sent in the get onboarding screens api internally
		params := map[string]string{
			ParamScreenId: c.Query(ParamScreenId),
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, OnboardingEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	case utils.POSTMethod:

		//Params to be sent in the create onboarding screen api internally
		screenRequest, err := parseCreateScreenRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, OnboardingEndpoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, screenRequest)

	case utils.PUTMethod:

		//Params to be sent in the update onboarding screen api internally
		screenRequest, err := parseUpdateScreenRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, OnboardingEndpoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, screenRequest)

	case utils.DELETEMethod:

		//Params to be sent in the delete onboarding screen api internally
		screenRequest, err := parseDeleteScreenRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, OnboardingEndpoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, screenRequest)

	}
}

func parseCreateScreenRequest(c *gin.Context) (*CreateScreenRequest, error) {
	//POST body params
	var csr CreateScreenRequest
	if err := c.ShouldBindJSON(&csr); err != nil {
		return nil, err
	}

	return &csr, nil
}

func parseUpdateScreenRequest(c *gin.Context) (*UpdateScreenRequest, error) {
	//POST body params
	var usr UpdateScreenRequest
	if err := c.ShouldBindJSON(&usr); err != nil {
		return nil, err
	}

	return &usr, nil
}

func parseDeleteScreenRequest(c *gin.Context) (*DeleteScreenRequest, error) {
	//POST body params
	var dsr DeleteScreenRequest
	if err := c.ShouldBindJSON(&dsr); err != nil {
		return nil, err
	}

	return &dsr, nil
}
