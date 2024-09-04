package sdk

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/handlers/token"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type SdkLoginRequest struct {
	Name 				string		`json:"name"`
	ImageUrl 			string		`json:"image_url"`
	OrganisationName 	string		`json:"organisation_name"`
	MobileNo    		string 		`json:"mobile_no"`
	CountryCode 		string   	`json:"country_code"`
	EmailID     		string 		`json:"email"`
	TokenExpiryBeta    	int64       `json:"token_expiry_beta,omitempty"`
	RTMTokenExpiryBeta 	int64       `json:"rtm_token_expiry_beta,omitempty"`
	DeviceID           	string      `json:"device_id,omitempty"`
}

func extractLoginDetailsFromVTM(vtm *constants.VerifyTokenMeta, slr *SdkLoginRequest) *SdkLoginRequest {

	if vtm.EmailID != "" {
		slr.EmailID = vtm.EmailID
	}

	if vtm.MobileNo != "" {
		slr.MobileNo = vtm.MobileNo
	}

	if vtm.CountryCode != "" {
		slr.CountryCode = vtm.CountryCode
	}

	return slr
}

func parseSdkLoginRequest(c *gin.Context) (*SdkLoginRequest, error) {
	// POST body params
	var slr SdkLoginRequest

	if err := c.ShouldBindJSON(&slr); err != nil {
		return nil, err
	}

	return &slr, nil
}

func SdkLogin(c *gin.Context) {
	
	sdkLoginRequest, err := parseSdkLoginRequest(c)
	if err != nil {
		// If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	verifyTokenMeta, ok := c.Get(constants.ParamVTM)
	if ok {
		sdkLoginRequest = extractLoginDetailsFromVTM(verifyTokenMeta.(*constants.VerifyTokenMeta), sdkLoginRequest)
	}
	
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, SdkLoginEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, ""), nil, sdkLoginRequest)
	if respBytes == nil {
		return
	}

	// Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	dataResponse := apiCR.Response

	// If flow succeeds, fetch User object
	userObject := apiCR.Response[user.ResponseUser].(map[string]interface{})

	// If flow succeeds
	userUniqueID := userObject[user.ResponseUserUniqueId].(string)
	userIsGuest := userObject[user.ResponseUserIsGuest].(bool)

	// Create login and refresh token
	ltm, rtm, err := token.CreateLTMAndRTM(userUniqueID, c.GetHeader(utils.HeadersApiKey), sdkLoginRequest.TokenExpiryBeta, sdkLoginRequest.RTMTokenExpiryBeta, userIsGuest, sdkLoginRequest.DeviceID)
	if err != nil {
		// If token creation fails
		utils.GeneralAPIError(c, err.Error())
		return
	}

	// Set ltm and user_unique_id in context
	ltm.UserUniqueID = userUniqueID
	c.Set(constants.ParamLTM, ltm)

	dataResponse[constants.ParamAccessToken] = ltm.AccessToken
	dataResponse[constants.ParamRefreshToken] = rtm.RefreshToken

	// Generate response
	utils.GenerateResponse(c, dataResponse, true)
}