package otp

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/utils"
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
		Url:     client.CoreServiceBaseURL + "/api/generate_otp",
		Params:  params,
		CustomHeaders: nil,
	}
	res, _ := client.GetRequest(&options)

	// marshaling and unmarshaling of response
	var resp utils.Response
	response, _ := json.Marshal(res)
	json.Unmarshal(response, &resp)

	//Check api/generate_otp success and response
	if !resp.Success {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success:      false,
			ErrorMessage: resp.ErrorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
	})

}
