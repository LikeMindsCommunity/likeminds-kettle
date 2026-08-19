package community

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/constants"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/chatroom"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/feed"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
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

		//Get Community Feed
		getCommunityFeedInternal(c, userId)

	}
}

func getCommunityFeedInternal(c *gin.Context, userId string) {

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
			fetchAndRespondPostDependentFeed(c, userId, params)

		} else if params[ParamOrderType] == strconv.Itoa(OrderTypeRecentlyActive) || params[ParamOrderType] == strconv.Itoa(OrderTypeMostMessages) {
			fetchAndRespondPostIndependentFeed(c, userId, params)

		} else {
			utils.GeneralBadRequestError(c, "invalid order_type sent")
		}

	} else {

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, CommunityFetchFeedEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, true)
	}
}

func fetchAndRespondPostDependentFeed(c *gin.Context, userId string, params map[string]string) {
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

	chatroomIds := []int{}
	finalPostCounts := map[string]int{}

	chatrooms, ok := dataResponse["chatrooms"]
	if ok {
		for _, chatroom := range chatrooms.([]interface{}) {
			chatroom_id, ok := chatroom.(map[string]interface{})["id"]
			if ok {
				chatroomIds = append(chatroomIds, int(chatroom_id.(float64)))
				finalPostCounts[strconv.Itoa(int(chatroom_id.(float64)))] = 0
			}
		}
	}

	tempParams, _ := json.Marshal(chatroomIds)

	//Generate params for ExploreFeed Request
	exploreFeedParams := map[string]string{
		ParamOrderType:   params[ParamOrderType],
		ParamChatroomIds: fmt.Sprintf("%v", string(tempParams)),
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
	if postCounts, ok := swarmDataResponse[constants.ResponseKeyPostCounts]; ok {
		for chatroom_id, postCount := range postCounts.(map[string]interface{}) {
			finalPostCounts[chatroom_id] = int(postCount.(float64))
		}
	}

	dataResponse[constants.ResponseKeyPostCounts] = finalPostCounts

	//Send Response
	utils.GenerateResponse(c, dataResponse, true)
}

func fetchAndRespondPostIndependentFeed(c *gin.Context, userId string, params map[string]string) {

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
		for _, chatroomId := range chatroomIds.([]interface{}) {
			excludedChatroomIds = append(excludedChatroomIds, int(chatroomId.(float64)))
		}
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
	postCountMap := map[string]int{}

	chatroomIds, ok = swarmDataResponse["chatroom_ids"]
	if ok {
		for _, chatroomId := range chatroomIds.([]interface{}) {
			selectedChatroomIds = append(selectedChatroomIds, int(chatroomId.(float64)))
			postCountMap[strconv.Itoa(int(chatroomId.(float64)))] = 0
		}
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
	if postCounts, ok := swarmDataResponse[constants.ResponseKeyPostCounts]; ok {
		for chatroom_id, post_count := range postCounts.(map[string]interface{}) {
			postCountMap[chatroom_id] = int(post_count.(float64))
		}
	}

	caravanDataResponse[constants.ResponseKeyPostCounts] = postCountMap

	//Send Response
	utils.GenerateResponse(c, caravanDataResponse, true)
}
