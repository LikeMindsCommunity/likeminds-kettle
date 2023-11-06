package sdk

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/community"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// InitiateSDKEndPoint | togther service user initiate endpoint
const InitiateSDKEndPoint = "/api/sdk/initiate"

// InitiateSDKRequest | user initiate request schema
type InitiateSDKRequest struct {
	UserName        string                            `json:"user_name"`
	UserUniqueID    string                            `json:"user_unique_id"`
	UUID            string                            `json:"uuid,omitempty"`
	ImageURL        string                            `json:"image_url"`
	IsGuest         bool                              `json:"is_guest"`
	QuestionAnswers []community.QuestionAnswerWithInt `json:"question_answers"`
	User            user.User                         `json:"user,omitempty"`
	TokenExpiryBeta int64                             `json:"token_expiry_beta,omitempty"`
	SharedBy        string                            `json:"shared_by,omitempty"`
}

// InitiateSDK is used to initiate sdk
func InitiateSDK(c *gin.Context) {

	// Body to be sent in the initiate SDK api internally
	initiateSDKRequest, err := parseInitiateSDKRequest(c)

	if err != nil {
		// If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	verifyTokenMeta, ok := c.Get(token.ParamVTM)

	if ok {
		vtm := verifyTokenMeta.(*token.VerifyTokenMeta)

		if initiateSDKRequest.User.Name == "" && initiateSDKRequest.UserName != "" {
			initiateSDKRequest.User.Name = initiateSDKRequest.UserName
		}

		if initiateSDKRequest.User.Email == "" && vtm.EmailID != "" {
			initiateSDKRequest.User.Email = vtm.EmailID
		}

		if initiateSDKRequest.User.MobileNo == "" && vtm.MobileNo != "" {
			initiateSDKRequest.User.MobileNo = vtm.MobileNo
		}

		if initiateSDKRequest.User.CountryCode == "" && vtm.CountryCode != "" {
			initiateSDKRequest.User.CountryCode = vtm.CountryCode
		}
	}

	// Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, InitiateSDKEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, ""), nil, initiateSDKRequest)
	if respBytes == nil {
		return
	}

	// Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	// Send response with login, refresh token and api/sdk/initiate response
	dataResponse := apiCR.Response

	// If flow succeeds
	userUniqueID := apiCR.Response[user.ResponseUser].(map[string]interface{})[user.ResponseUserUniqueId].(string)

	// Create login and refresh token
	ltm, rtm, err := token.CreateLTMAndRTM(userUniqueID, c.GetHeader(utils.HeadersApiKey), initiateSDKRequest.TokenExpiryBeta)

	if err != nil {
		// If token creation fails
		utils.GeneralAPIError(c, err.Error())
		return
	}

	// Set ltm and user_unique_id in context
	ltm.UserUniqueID = userUniqueID
	c.Set(token.ParamLTM, ltm)

	dataResponse[token.ParamAccessToken] = ltm.AccessToken
	dataResponse[token.ParamRefreshToken] = rtm.RefreshToken

	// Generate response
	utils.GenerateResponse(c, dataResponse, true)
}

func parseInitiateSDKRequest(c *gin.Context) (*InitiateSDKRequest, error) {
	//POST body params
	var isr InitiateSDKRequest

	if err := c.ShouldBindJSON(&isr); err != nil {
		return nil, err
	}

	// If uuid is passed in the request, use it as user_unique_id
	if isr.UUID != "" {
		isr.UserUniqueID = isr.UUID
	}

	return &isr, nil
}
