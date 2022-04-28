package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
)

const MergeAccountEndPoint = "/api/merge_account"
const ParamMobileNo = "mobile_no"
const ParamCountryCode = "country_code"

//MergeAccount used when user wants to merge account and generate login and refresh tokens
func MergeAccount(c *gin.Context) {
	//Check if request has valid login token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.SomethingWentWrongError(c)
		return
	}

	//Create headers from login token
	headers := make(map[string]interface{})
	headers[utils.HeadersMemberId] = ltm.UserID
	//Params to be sent in the api/merge_account request
	params := map[string]string{
		//TODO - get mobile number and country code
	}
	//Create internal API client
	apiClient := api_client.NewAPIClient()
	//Send request
	respBytes, err := apiClient.PostRequest(&api_client.PostRequestOptions{
		Url:           apiClient.CoreServiceBaseURL + MergeAccountEndPoint,
		CustomHeaders: headers,
		Body:          params,
	})
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
	}

	if !apiCR.Success {
		//If api/merge_account returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}
	//Send response with api/merge_account response
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    apiCR.Response,
	})
}
