package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/token"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type User struct {
	MobileNo         string `json:"mobile_no,omitempty"`
	CountryCode      string `json:"country_code,omitempty"`
	Name             string `json:"name,omitempty"`
	Email            string `json:"email,omitempty"`
	ImageUrl         string `json:"image_url,omitempty"`
	OrganisationName string `json:"organisation_name,omitempty"`
	UserUniqueId     string `json:"user_unique_id,omitempty"`
}

type LoginRequest struct {
	LoginType string `json:"type" binding:"required"`
	User      User   `json:"user"`
}

// Login used when user is signing up and generate login and refresh tokens
func Login(c *gin.Context) {

	//Body to be sent in the login api internally
	loginRequest, err := parseLoginRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, LoginEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, ""), nil, updateLoginRequest(loginRequest))
	if respBytes == nil {
		return
	}

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	userObject := apiCR.Response[ResponseUser].(map[string]interface{})
	userUniqueID := userObject[ResponseUserUniqueId].(string)
	userIsGuest := userObject[ResponseUserIsGuest].(bool)

	//Create login and refresh token
	ltm, rtm, err := token.CreateLTMAndRTM(userUniqueID, "", token.BETA_AUTH_TOKEN_EXPIRY, token.DEFAULT_TOKEN_EXPIRY, userIsGuest)
	if err != nil {
		//If token creation fails
		utils.GeneralAPIError(c, err.Error())
		return
	}

	// Set ltm and user_unique_id in context
	ltm.UserUniqueID = userUniqueID
	c.Set(constants.ParamLTM, ltm)

	//Send response with login, refresh token and api/user/login response
	dataResponse := apiCR.Response
	dataResponse[constants.ParamAccessToken] = ltm.AccessToken
	dataResponse[constants.ParamRefreshToken] = rtm.RefreshToken

	//Generate response
	utils.GenerateResponse(c, dataResponse, true)
}

func parseLoginRequest(c *gin.Context) (*LoginRequest, error) {
	//POST body params
	var lr LoginRequest

	if err := c.ShouldBindJSON(&lr); err != nil {
		return nil, err
	}

	return &lr, nil
}

func updateLoginRequest(lr *LoginRequest) interface{} {
	updatedLr := make(map[string]interface{})
	user := make(map[string]interface{})

	user[UserName] = lr.User.Name
	user[UserEmail] = lr.User.Email
	user[UserImageUrl] = lr.User.ImageUrl
	user[UserOrganisationName] = lr.User.OrganisationName
	user[ResponseUserUniqueId] = lr.User.UserUniqueId

	updatedLr[UserMobileNo] = lr.User.MobileNo
	updatedLr[UserCountryCode] = lr.User.CountryCode
	updatedLr[UserLoginType] = lr.LoginType
	updatedLr[ResponseUser] = user

	return updatedLr
}
