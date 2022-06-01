package community

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type MemberRequest struct {
	UserName     string `json:"user_name"`
	UserUniqueId string `json:"user_unique_id"`
	ImageUrl     string `json:"image_url"`
}

//GetMember is used to get community members
func GetMember(c *gin.Context) {
	Member(c, utils.GETMethod)
}

//AddMember is used to add a member in community
func AddMember(c *gin.Context) {
	Member(c, utils.POSTMethod)
}

//EditMember is used to edit member in community
func EditMember(c *gin.Context) {
	Member(c, utils.PUTMethod)
}

//Member method handles members for a commuinty
func Member(c *gin.Context, method int) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Call GET api/bot to get bot
	response := user.GetBotResponse(c, utils.GETMethod)
	if response == nil {
		return
	}

	//Send request
	var respBytes []byte
	var err error
	switch method {
	case utils.GETMethod:

		respBytes = getMemberInternal(c, client, response)

	case utils.POSTMethod:

		respBytes = addMemberInternal(c, client, response)

	case utils.PUTMethod:

		respBytes = editMemberInternal(c, client, response)

	}

	if respBytes == nil {
		return
	}

	//Parse response
	var apiCR api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
	if err != nil {
		//Internal unmarshal error
		utils.GeneralAPIError(c, err.Error())
		return
	}
	if !apiCR.Success {
		//If chatroom apis returns success as false
		c.JSON(http.StatusInternalServerError, apiCR)
		return
	}
	//If flow succeeds
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    apiCR.Response,
	})
}

func parseMemberRequest(c *gin.Context) (*MemberRequest, error) {
	//POST body params
	var mr MemberRequest

	if err := c.ShouldBindJSON(&mr); err != nil {
		return nil, err
	}

	return &mr, nil
}

func getMemberInternal(c *gin.Context, client *api_client.APIClient, response *utils.Response) []byte {
	var options api_client.GetRequestOptions

	//GET Request params
	page := c.Query(ParamPage)
	if page == "" {
		//If page is missing, call api/community/fetch_members_meta api internally

		options = api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + FetchMembersMetaEndPoint,
			CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
		}

	} else {
		//else, call api/v1/all_members api internally

		//Params to be sent in the api/v1/all_members request
		params := map[string]string{
			ParamCommunityId: c.Query(ParamCommunityId),
		}

		options = api_client.GetRequestOptions{
			Url:           client.CoreServiceBaseURL + AllMembersV1EndPoint,
			CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
			Params:        params,
		}
	}

	respBytes, err := client.GetRequest(&options)
	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}

func addMemberInternal(c *gin.Context, client *api_client.APIClient, response *utils.Response) []byte {
	//Body to be sent in the api/community/member POST request
	memberRequest, err := parseMemberRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + CommunityMemberEndPoint,
		Body:          memberRequest,
		CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
	}

	respBytes, err := client.PostRequest(&options)

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}

func editMemberInternal(c *gin.Context, client *api_client.APIClient, response *utils.Response) []byte {
	//Body to be sent in the api/community/member PUT request
	memberRequest, err := parseMemberRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + CommunityMemberEndPoint,
		Body:          memberRequest,
		CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
	}

	respBytes, err := client.PutRequest(&options)

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}
