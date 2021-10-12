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
		c.JSON(http.StatusBadRequest, utils.Response{
			Success:      false,
			ErrorMessage: "Query params missing!",
		})
		return
	}

	//Params to be sent in the request
	params := map[string]string{
		"country_code": countryCode,
		"mobile_no":    mobileNo,
	}
	//http client and request options
	client := api_client.NewAPIClient()
	options := api_client.GetRequestOptions{
		Url:           client.CoreServiceBaseURL + "/api/generate_otp",
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

	//Check api/generate_otp success and response
	if !apiCR.Success {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success:      false,
			ErrorMessage: apiCR.ErrorMessage,
		})
		return
	}
	var apiCR api_client.APIClientResponse
	err = json.Unmarshal(respBytes, &apiCR)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api_client.APIClientResponse{
			Success:      false,
			ErrorMessage: "Something went wrong! Please try after sometime",
		})
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: apiCR.Success,
	})
}
