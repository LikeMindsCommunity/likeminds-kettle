package otp

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

// GenerateOTP is used to generate otp
func GenerateOTP(c *gin.Context) {
	//GET Request params
	mobileNo := c.Query("mobile_no")
	countryCode := c.Query("country_code")
	if mobileNo == "" || countryCode == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//Params to be sent in the api/generate_otp request
	params := map[string]string{
		"country_code": countryCode,
		"mobile_no":    mobileNo,
	}
	//Create internal API client
	client := api_client.NewAPIClient()
	options := api_client.GetRequestOptions{
		Url:           client.CoreServiceBaseURL + "/api/generate_otp",
		Params:        params,
		CustomHeaders: nil,
	}
	//Send request
	respBytes, err := client.GetRequest(&options)
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
		utils.SomethingWentWrongError(c)
		return
	}

	if !apiCR.Success {
		//If api/generate_otp returns success as false
		utils.APIClientError(c, apiCR)
		return
	}
	//If flow succeeds
	c.JSON(http.StatusOK, utils.Response{
		Success: apiCR.Success,
	})
}
