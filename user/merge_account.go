package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

//MergeAccount used when user wants to merge account and generate login and refresh tokens
func MergeAccount(c *gin.Context) {

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Create headers from login token
	headers := make(map[string]interface{})
	headers[utils.HeadersMemberId] = ltm.UserUniqueID
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
	utils.ParseResponse(c, respBytes)
}
