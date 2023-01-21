package community

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/feed"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

// GetFeed method will get community feed data
func GetCommunityFeed(c *gin.Context) {
	CommunityFeed(c, utils.GETMethod)
}

// CommunityFeed mthod handles community feed
func CommunityFeed(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	switch method {
	case utils.GETMethod:

		//Params to be sent in the fetch_feed request
		params := map[string]string{
			ChatroomIDParam:               c.Query(ChatroomIDParam),
			ParamPinned:                   c.Query(ParamPinned),
			chatroom.ParamScrollDirection: c.Query(chatroom.ParamScrollDirection),
			ParamOrderType:                c.Query(ParamOrderType),
			ParamPage:                     c.Query(ParamPage),
			ParamPageSize:                 c.Query(ParamPageSize),
			ParamType:                     c.Query(ParamType),
		}

		if params[ParamType] == strconv.Itoa(PostFeedType) {
			if params[ParamOrderType] == strconv.Itoa(OrderTypeNewest) || params[ParamOrderType] == strconv.Itoa(OrderTypeMostParticipants) {
				//Generate params for CommunityPostFeed Request
				communityFeedParams := map[string]string{
					ParamPinned:    params[ParamPinned],
					ParamOrderType: params[ParamOrderType],
					ParamPage:      params[ParamPage],
					ParamPageSize:  params[ParamPageSize],
				}

				//Send Request to get post community feed on Caravan Service
				respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, CommunityFetchPostFeedEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), communityFeedParams, nil)

				//Validate response
				apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
				if apiCR == nil {
					return
				}

				//If flow succeeds
				dataResponse := apiCR.Response
				chatroom_ids := []int{}

				chatrooms, ok := dataResponse["chatrooms"]
				if ok {
					for _, chatroom := range chatrooms.([]map[string]interface{}) {
						chatroom_id, ok := chatroom["id"]
						if ok {
							convertedChatroomId, _ := strconv.Atoi(chatroom_id.(string))
							chatroom_ids = append(chatroom_ids, convertedChatroomId)
						}
					}
				}

				temp_params, _ := json.Marshal(chatroom_ids)

				//Generate params for ExploreFeed Request
				exploreFeedParams := map[string]string{
					ParamOrderType:   params[ParamOrderType],
					ParamChatroomIds: fmt.Sprintf("%v", string(temp_params)),
				}

				//Send Request to get post community feed on Swarm Service
				respBytes, statusCode = utils.GetRequestResponse(c, utils.SwarmService, feed.FeedExploreEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), exploreFeedParams, nil)

				//Validate response
				apiCR = utils.ValidateClientResponse(c, respBytes, statusCode)
				if apiCR == nil {
					return
				}

				//If flow succeeds
				swarmDataResponse := apiCR.Response
				if postCounts, ok := swarmDataResponse["post_counts"]; ok {
					dataResponse["post_counts"] = postCounts
				}

				//Send Response
				utils.GenerateResponse(c, dataResponse)
			} else if params[ParamOrderType] == strconv.Itoa(OrderTypeRecentlyActive) || params[ParamOrderType] == strconv.Itoa(OrderTypeMostMessages) {
				//Send Request to get excluded chatrooms list on Caravan Service
				respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, CommunityExcludedChatroomsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

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
					excludedChatroomIds = chatroomIds.([]int)
				}

				temp_params, _ := json.Marshal(excludedChatroomIds)

				//Generate params for ExploreFeed Request
				exploreFeedParams := map[string]string{
					ParamOrderType:           params[ParamOrderType],
					ParamPage:                params[ParamPage],
					ParamPageSize:            params[ParamPageSize],
					ParamExcludedChatroomIds: fmt.Sprintf("%v", string(temp_params)),
				}

				//Send Request to get post community feed on Swarm Service
				respBytes, statusCode = utils.GetRequestResponse(c, utils.SwarmService, feed.FeedExploreEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), exploreFeedParams, nil)

				//Validate response
				apiCR = utils.ValidateClientResponse(c, respBytes, statusCode)
				if apiCR == nil {
					return
				}

				//If flow succeeds
				swarmDataResponse := apiCR.Response
				selectedChatroomIds := []int{}

				chatroomIds, ok = swarmDataResponse["chatroom_ids"]
				if ok {
					selectedChatroomIds = chatroomIds.([]int)
				}

				temp_params, _ = json.Marshal(selectedChatroomIds)

				//Generate params for CommunityPostFeed Request
				communityFeedParams := map[string]string{
					ParamOrderType:   params[ParamOrderType],
					ParamChatroomIds: fmt.Sprintf("%v", string(temp_params)),
				}

				//Send Request to get post community feed on Caravan Service
				respBytes, statusCode = utils.GetRequestResponse(c, utils.CoreService, CommunityFetchPostFeedEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), communityFeedParams, nil)

				//Validate response
				apiCR = utils.ValidateClientResponse(c, respBytes, statusCode)
				if apiCR == nil {
					return
				}

				//If flow succeeds
				caravanDataResponse := apiCR.Response
				if postCounts, ok := swarmDataResponse["post_counts"]; ok {
					caravanDataResponse["post_counts"] = postCounts
				}

				//Send Response
				utils.GenerateResponse(c, caravanDataResponse)
			} else {
				utils.GeneralBadRequestError(c, "invalid order_type sent")
				return
			}

		} else {
			//Send Request to get chat community feed
			utils.SendRequest(c, utils.CoreService, CommunityFetchFeedEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		}
	}

}
