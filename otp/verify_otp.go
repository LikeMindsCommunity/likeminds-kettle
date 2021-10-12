package otp

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/core_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type User struct {
	Id string `json:"id"`
}

type VerifyOTPResponse struct {
	Success       bool   `json:"success"`
	ErrorMessage  string `json:"error_message"`
	ProfileExists bool   `json:"profile_exists"`
	User          User   `json:"user"`
}

// VerifyOTP is used to verify otp and generate VTM Token
func VerifyOTP(c *gin.Context) {
	//GET Request params
	otp := c.Query("otp")
	mobileNo := c.Query("mobile_no")
	countryCode := c.Query("country_code")
	if otp == "" || mobileNo == "" || countryCode == "" {
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
		"otp":          otp,
	}
	//http client and request options
	client := core_client.NewClient()
	options := core_client.GetRequestOptions{
		Url:     client.BaseURL + "/api/verify_otp",
		Params:  params,
		Headers: nil,
	}
	res, _ := client.GetRequest(&options)

	// marshaling and unmarshaling of response
	var resp VerifyOTPResponse
	response, _ := json.Marshal(res)
	json.Unmarshal(response, &resp)

	//Check api/verify_otp success and response
	if !resp.Success {
		c.JSON(http.StatusInternalServerError, utils.AuthenticationResponse{
			Success:      false,
			ErrorMessage: resp.ErrorMessage,
		})
		return
	}
	if resp.ProfileExists {
		//Create login and refresh token meta from the response received in api/verify_otp
		ltm, rtm, err := token.CreateLTMAndRTM(mobileNo, countryCode, resp.User.Id)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, utils.AuthenticationResponse{
				Success:      false,
				ErrorMessage: err.Error(),
			})
			return
		}
		token := map[string]string{
			"access_token":  ltm.AccessToken,
			"refresh_token": rtm.RefreshToken,
		}
		c.JSON(http.StatusOK, utils.AuthenticationResponse{
			Success: true,
			Data:    token,
		})
	} else {
		//Create verify token meta from the response received in api/verify_otp
		vtm, err := token.CreateVTM(mobileNo, countryCode)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		token := map[string]string{
			"access_token": vtm.AccessToken,
		}
		c.JSON(http.StatusOK, utils.AuthenticationResponse{
			Success: true,
			Data:    token,
		})
	}
}
