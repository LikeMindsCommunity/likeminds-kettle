package otp

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/core_client"
	"github.com/nateshr/likeminds-authentication/utils"
)

type GenerateOTPResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message"`
}

// GenerateOTP is used to generate otp
func GenerateOTP(c *gin.Context) {
	//GET Request params
	mobileNo := c.Query("mobile_no")
	countryCode := c.Query("country_code")
	if mobileNo == "" || countryCode == "" {
		c.JSON(http.StatusBadRequest, utils.AuthenticationResponse{
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
	client := core_client.NewClient()
	options := core_client.GetRequestOptions{
		Url:     client.BaseURL + "/api/generate_otp",
		Params:  params,
		Headers: nil,
	}
	res, _ := client.GetRequest(&options)

	// marshaling and unmarshaling of response
	var resp GenerateOTPResponse
	response, _ := json.Marshal(res)
	json.Unmarshal(response, &resp)

	//Check api/generate_otp success and response
	if !resp.Success {
		c.JSON(http.StatusInternalServerError, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: resp.ErrorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, utils.AuthenticationResponse{
		Success: true,
	})

}
