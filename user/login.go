package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

type LoginRequest struct {
	UserAcquisitionURL string `json:"user_acquisition_url"`
	LoginJSON          string `json:"login_json"`
	GoogleIDToken      string `json:"google_id_token"`
	LoginType          string `json:"type" binding:"required"`
}

//Login used when user is signing up and generate LTM and RTM tokens
func Login(c *gin.Context) {
	vtm, ok := c.MustGet("vtm").(*token.VerifyTokenMeta)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success:      false,
			ErrorMessage: "Something went wrong! Please try after sometime",
		})
		return
	}
	var lr LoginRequest
	if err := c.ShouldBindJSON(&lr); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.Response{
			Success:      false,
			ErrorMessage: "Body params missing!",
		})
		return
	}

	//Params to be sent in the request
	params, err := utils.RequestParamsToMap(lr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success:      false,
			ErrorMessage: "Something went wrong! Please try after sometime",
		})
	}
	//http client and request options
	client := api_client.NewAPIClient()
	options := api_client.GetRequestOptions{
		Url:           client.CoreServiceBaseURL + "/api/user/login",
		Params:        params,
		CustomHeaders: nil,
	}
	//Unmarshaling of response
	respBytes, err := client.GetRequest(&options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success:      false,
			ErrorMessage: err.Error(),
		})
		return
	}
	var apiCR api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success:      false,
			ErrorMessage: "Something went wrong! Please try after sometime",
		})
	}

	//Check api/user/login success and response
	if !apiCR.Success {
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}
	emailExists := apiCR.Response["email_exists"].(bool)
	mobileNo := vtm.VerifiedMobileNo
	countryCode := vtm.CountryCode
	if emailExists {
		//Create verify tokenResponse meta from the response received in VTM
		vtm, err := token.CreateVTM(mobileNo, countryCode)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		dataResponse := apiCR.Response
		dataResponse["access_token"] = vtm.AccessToken
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Data:    dataResponse,
		})
	} else {
		userID := apiCR.Response["user"].(map[string]interface{})["id"].(float64)
		//Create login and refresh dataResponse meta from the response received in api/verify_otp
		ltm, rtm, err := token.CreateLTMAndRTM(mobileNo, countryCode, userID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, utils.Response{
				Success:      false,
				ErrorMessage: err.Error(),
			})
			return
		}
		dataResponse := apiCR.Response
		dataResponse["access_token"] = ltm.AccessToken
		dataResponse["refresh_token"] = rtm.RefreshToken
		c.JSON(http.StatusOK, utils.Response{
			Success: true,
			Data:    dataResponse,
		})
	}
	return
}
