package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
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
	//POST body params
	var lr LoginRequest
	if err := c.ShouldBindJSON(&lr); err != nil {
		//If POST body params are missing
		utils.POSTBodyParamsMissingError(c)
		return
	}

	//Create internal API client
	client := api_client.NewAPIClient()
	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + LoginEndPoint,
		Body:          updateLoginRequest(lr),
		CustomHeaders: nil,
	}
	//Send request
	respBytes, err := client.PostRequest(&options)
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
	}

	if !apiCR.Success {
		//If api/user/login returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
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
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    dataResponse,
	})
}

func updateLoginRequest(lr LoginRequest) map[string]interface{} {
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
