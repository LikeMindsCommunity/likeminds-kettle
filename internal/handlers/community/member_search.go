package community

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

func MemberSearch(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	// Params to be sent in the api/search/member_directory request
	requestParams := map[string]string{
		SearchParam:                 c.Query(SearchParam),
		SearchTypeParam:             c.Query(SearchTypeParam),
		ParamPage:                   c.Query(ParamPage),
		ParamPageSize:               c.Query(ParamPageSize),
		ParamOrderType:              c.Query(ParamOrderType),
		ParamMemberStates:           c.Query(ParamMemberStates),
		ParamQuestionAnswersVersion: c.Query(ParamQuestionAnswersVersion),
		ParamExcludeSelfMember:      c.Query(ParamExcludeSelfMember),
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, UserSearchEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode, true)

}
