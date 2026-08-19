package community

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/chatroom"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// GetTaggingList is used to fetch the tag members list for a specific feedroom
func GetTaggingList(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	chatroomId := c.Query(ChatroomIDParam)
	if c.Query(FeedroomIDParam) != "" {
		chatroomId = c.Query(FeedroomIDParam)
	}

	//Params to be sent with pagination and search support in APIs internally
	params := map[string]string{
		ChatroomIDParam: chatroomId,
		ParamPage:       c.Query(ParamPage),
		ParamPageSize:   c.Query(ParamPageSize),
		SearchName:      c.Query(SearchName),
	}

	//Params Validation
	if params[chatroom.ParamChatroomId] == "" {

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchMembersMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, true)

	} else {

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, chatroom.GetTaggingListEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, true)

	}
}
