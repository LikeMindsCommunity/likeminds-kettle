package search

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/community"
	"github.com/nateshr/likeminds-authentication/feed"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// PostSearch is used to perform search on the posts
func PostSearch(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Fetch member access to view post
	success, response := user.FetchMemberAccess(c, feed.VIEW_POST_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request to get excluded chatrooms list on Caravan Service
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, community.CommunityExcludedChatroomsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	excludedChatroomIds := []int{}

	chatroomIds, ok := dataResponse["chatroom_ids"]
	if ok {
		for _, chatroomId := range chatroomIds.([]interface{}) {
			excludedChatroomIds = append(excludedChatroomIds, int(chatroomId.(float64)))
		}
	}

	temp_params, _ := json.Marshal(excludedChatroomIds)

	//Params to be sent in the post search api internally
	params := map[string]string{
		ParamSearch:              c.Query(ParamSearch),
		ParamSearchType:          c.Query(ParamSearchType),
		ParamPage:                c.Query(ParamPage),
		ParamPageSize:            c.Query(ParamPageSize),
		ParamExcludedChatroomIds: fmt.Sprintf("%v", string(temp_params)),
		ParamUserIsCm:            fmt.Sprint(response.IsCm),
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, PostSearchEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}
