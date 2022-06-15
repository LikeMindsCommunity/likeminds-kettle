package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type AddMemberRequest struct {
	UserName     string `json:"user_name" binding:"required"`
	UserUniqueId string `json:"user_unique_id"`
	ImageUrl     string `json:"image_url"`
}

type EditMemberRequest struct {
	UserName     string `json:"user_name"`
	UserUniqueId string `json:"user_unique_id" binding:"required"`
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

	//Parse Response
	utils.ParseResponse(c, respBytes)
}

func parseAddMemberRequest(c *gin.Context) (*AddMemberRequest, error) {
	//POST body params
	var amr AddMemberRequest

	if err := c.ShouldBindJSON(&amr); err != nil {
		return nil, err
	}

	return &amr, nil
}

func parseEditMemberRequest(c *gin.Context) (*EditMemberRequest, error) {
	//POST body params
	var emr EditMemberRequest

	if err := c.ShouldBindJSON(&emr); err != nil {
		return nil, err
	}

	return &emr, nil
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
	memberRequest, err := parseAddMemberRequest(c)

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

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeRaw)

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}

func editMemberInternal(c *gin.Context, client *api_client.APIClient, response *utils.Response) []byte {
	//Body to be sent in the api/community/member PUT request
	memberRequest, err := parseEditMemberRequest(c)

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
