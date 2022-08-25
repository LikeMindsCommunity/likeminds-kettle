package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type User struct {
	MobileNo         string `json:"mobile_no"`
	CountryCode      string `json:"country_code"`
	Name             string `json:"name" binding:"required"`
	Email            string `json:"email"`
	ImageUrl         string `json:"image_url"`
	OrganisationName string `json:"organisation_name"`
	UserUniqueId     string `json:"user_unique_id"`
}

type LoginRequest struct {
	LoginType string `json:"type" binding:"required"`
	User      User   `json:"user"`
}

//Login used when user is signing up and generate login and refresh tokens
func Login(c *gin.Context) {

	//Body to be sent in the login api internally
	loginRequest, err := parseLoginRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
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
	userUniqueID := apiCR.Response[ResponseUser].(map[string]interface{})[ResponseUserUniqueId].(string)
	//Create login and refresh token

	ltm, rtm, err := token.CreateLTMAndRTM(userUniqueID)
	if err != nil {
		//If token creation fails
		utils.GeneralAPIError(c, err.Error())
		return
	}
	//Send response with login, refresh token and api/user/login response
	dataResponse := apiCR.Response
	dataResponse[token.ParamAccessToken] = ltm.AccessToken
	dataResponse[token.ParamRefreshToken] = rtm.RefreshToken

	//Generate response
	utils.GenerateResponse(c, dataResponse)
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
