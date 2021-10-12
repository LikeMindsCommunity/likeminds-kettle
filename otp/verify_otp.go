package otp

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

// VerifyOTP is used to verify otp and generate VTM Token
func VerifyOTP(c *gin.Context) {
	//GET Request params
	otp := c.Query("otp")
	mobileNo := c.Query("mobile_no")
	countryCode := c.Query("country_code")
	if otp == "" || mobileNo == "" || countryCode == "" {
		c.JSON(http.StatusBadRequest, api_client.APIClientResponse{
			Success:      false,
			ErrorMessage: "Query params missing!",
		})
		return
	}

	//Params to be sent in the request
	params := map[string]string{
		"country_code": countryCode,
		"mobile_no":    mobileNo,
		"otp":          otp,
	}
	//http client and request options
	client := api_client.NewAPIClient()
	options := api_client.GetRequestOptions{
		Url:           client.CoreServiceBaseURL + "/api/verify_otp",
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

	//Check api/verify_otp success and response
	if !apiCR.Success {
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}
	profileExists := apiCR.Response["profile_exists"].(bool)
	userID := apiCR.Response["user"].(map[string]interface{})["id"].(float64)
	if profileExists {
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
	} else {
		//Create verify tokenResponse meta from the response received in api/verify_otp
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
	}
	return
}
