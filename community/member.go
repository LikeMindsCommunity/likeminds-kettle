package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type AddMemberRequest struct {
	UserName     string `json:"user_name" binding:"required"`
	UserUniqueId string `json:"user_unique_id"`
	UUID         string `json:"uuid"`
	ImageUrl     string `json:"image_url"`
}

type EditMemberRequest struct {
	UserName     string `json:"user_name"`
	UserUniqueId string `json:"user_unique_id"`
	UUID         string `json:"uuid"`
	ImageUrl     string `json:"image_url"`
}

// GetMember is used to get community members
func GetMember(c *gin.Context) {
	Member(c, utils.GETMethod)
}

// AddMember is used to add a member in community
func AddMember(c *gin.Context) {
	Member(c, utils.POSTMethod)
}

// EditMember is used to edit member in community
func EditMember(c *gin.Context) {
	Member(c, utils.PUTMethod)
}

// Member method handles members for a commuinty
func Member(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Send request
	switch method {
	case utils.GETMethod:

		getMemberInternal(c, userId)

	case utils.POSTMethod:

		addMemberInternal(c, userId)

	case utils.PUTMethod:

		editMemberInternal(c, userId)

	}
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

func getMemberInternal(c *gin.Context, userId string) {

	//GET Request params
	page := c.Query(ParamPage)
	if page == "" {
		//If page is missing, call api/community/fetch_members_meta api internally

		//Send Request
		utils.SendRequest(c, utils.CoreService, FetchMembersMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
	} else {
		//else, call api/v1/all_members api internally

		//Params to be sent in the fetch all members api internally
		params := map[string]string{
			ParamPage:        page,
			ParamMemberState: c.Query(ParamMemberState),
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, AllMembersV1EndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	}
}

func addMemberInternal(c *gin.Context, userId string) {

	//Body to be sent in the add member api internally
	memberRequest, err := parseAddMemberRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, CommunityMemberEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, memberRequest)
}

func editMemberInternal(c *gin.Context, userId string) {

	//Body to be sent in the edit member api internally
	memberRequest, err := parseEditMemberRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, CommunityMemberEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, memberRequest)
}
