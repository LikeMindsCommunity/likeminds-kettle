package search

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/community"
	"github.com/nateshr/likeminds-authentication/conversation"
	"github.com/nateshr/likeminds-authentication/feed"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// GeneralSearch is used to perform search on the content
func GeneralSearch(c *gin.Context) {
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

	//Params to be sent in the general post search api internally
	params := map[string]string{
		ParamSearch:                c.Query(ParamSearch),
		ParamPage:                  c.Query(ParamPage),
		ParamPageSize:              c.Query(ParamPageSize),
		ParamUserIsCm:              fmt.Sprint(response.IsCm),
		chatroom.ParamFollowStatus: c.Query(chatroom.ParamFollowStatus),
	}

	dataResponse := map[string]interface{}{}

	params[ParamSearchType] = "title"
	//Send Request to fetch the chatroom search results
	respBytes, _, _ := utils.GetRequestResponseWithoutContext(utils.CoreService, chatroom.ChatroomSearchEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	var apiCR1 api_client.APIClientResponse
	err := api_client.UnmarshalAPIClientResponse(respBytes, &apiCR1)
	if err == nil {
		dataResponse["chatrooms"] = apiCR1.Response["chatrooms"]
	}

	//Send Request to fetch the message search results
	respBytes, _, _ = utils.GetRequestResponseWithoutContext(utils.CoreService, conversation.ConversationSearchEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	var apiCR2 api_client.APIClientResponse
	err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR2)
	if err == nil {
		dataResponse["conversations"] = apiCR2.Response["conversations"]
	}

	//Send Request to fetch the post search results

	//Send Request to get excluded chatrooms list on Caravan Service
	respBytes, _, _ = utils.GetRequestResponseWithoutContext(utils.CoreService, community.CommunityExcludedChatroomsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

	var apiCR3 api_client.APIClientResponse
	excludedChatroomIds := []int{}
	err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR3)
	if err == nil {
		chatroomIds, ok := dataResponse["chatroom_ids"]
		if ok {
			for _, chatroomId := range chatroomIds.([]interface{}) {
				excludedChatroomIds = append(excludedChatroomIds, int(chatroomId.(float64)))
			}
		}

		temp_params, _ := json.Marshal(excludedChatroomIds)

		params[ParamSearchType] = "text"
		params[ParamExcludedChatroomIds] = fmt.Sprintf("%v", string(temp_params))
		params[ParamUserIsCm] = fmt.Sprint(response.IsCm)

		//Send Request
		respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.SwarmService, PostSearchEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

		//Validate response
		apiCR := utils.ValidateClientResponseWithoutContext(respBytes, statusCode, err)

		if err == nil {
			dataResponse["posts"] = apiCR["posts"]

			// Fetch users meta for the posts
			if value, ok := dataResponse["posts"]; ok {

				posts := value.([]interface{})

				user_data, err := user.GetUsersMetaFromFeedData(utils.CreateHeaders(c, userId), posts)

				if err != nil {
					utils.GenerateResponse(c, nil)
					return
				}

				//Update user data in dataResponse
				dataResponse["users"] = user_data
			}
		}

	}

	//Generate Response
	utils.GenerateResponse(c, dataResponse)
}
