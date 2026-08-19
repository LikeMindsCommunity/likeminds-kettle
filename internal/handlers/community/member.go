package community

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
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

type DeleteMembersRequest struct {
	UUIDs  []string `json:"uuids" binding:"required"`
	Reason string   `json:"reason,omitempty"`
	TagId  int      `json:"tag_id,omitempty"`
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

func parseDeleteMembersRequest(c *gin.Context) (*DeleteMembersRequest, error) {
	//POST body params
	var dmr DeleteMembersRequest

	if err := c.ShouldBindJSON(&dmr); err != nil {
		return nil, err
	}

	return &dmr, nil
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

// RemoveMembers is used to remove members from community
func RemoveMembers(c *gin.Context) {
	Member(c, utils.DELETEMethod)
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

	case utils.DELETEMethod:

		removeMembersInternal(c, userId)

	}

}

func getMemberInternal(c *gin.Context, userId string) {

	//GET Request params
	page := c.Query(ParamPage)
	if page == "" {
		//If page is missing, call api/community/fetch_members_meta api internally

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchMembersMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, true)

	} else {
		//else, call api/v1/all_members api internally

		//Params to be sent in the fetch all members api internally
		params := map[string]string{
			ParamPage:                   page,
			ParamMemberState:            c.Query(ParamMemberState),
			ParamQuestionAnswersVersion: c.Query(ParamQuestionAnswersVersion),
			ParamFilterMemberRoles:      c.Query(ParamFilterMemberRoles),
			ParamExcludeSelfMember:      c.Query(ParamExcludeSelfMember),
		}

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, AllMembersV1EndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, true)
	}
}

func addMemberInternal(c *gin.Context, userId string) {

	//Body to be sent in the add member api internally
	memberRequest, err := parseAddMemberRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, CommunityMemberEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, memberRequest)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode, true)

}

func editMemberInternal(c *gin.Context, userId string) {

	//Body to be sent in the edit member api internally
	memberRequest, err := parseEditMemberRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, CommunityMemberEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, memberRequest)

}

func removeMembersInternal(c *gin.Context, userId string) {

	//Body to be sent in the remove member api internally
	memberRequest, err := parseDeleteMembersRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request to internal core service
	utils.SendRequest(c, utils.CoreService, CommunityMemberEndPoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, memberRequest)

}

// Exposed method to leave community
func LeaveCommunity(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, LeaveCommunityEndPoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, nil)
}
