package conversation

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/internal/cache"
	"github.com/nateshr/likeminds-authentication/internal/handlers/chatroom"
	"github.com/nateshr/likeminds-authentication/internal/handlers/community"
	"github.com/nateshr/likeminds-authentication/internal/handlers/pubsubPublish"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
	"net/http"
	"strconv"
	"time"
)

type PollObject struct {
	Text string `json:"text"`
}

type ConversationPreview struct {
	InternalLink string                    `json:"internal_link"`
	PreviewType  string                    `json:"preview_type"`
	PreviewText  string                    `json:"preview_text"`
	Title        string                    `json:"title"`
	Community    community.CommunityObject `json:"community"`
	Action       string                    `json:"action"`
	ActionRoute  string                    `json:"action_route"`
}

type ConversationAttachment struct {
	Name         string      `json:"name,omitempty"`
	Url          string      `json:"url,omitempty"`
	Type         string      `json:"type,omitempty"`
	ThumbnailUrl string      `json:"thumbnail_url,omitempty"`
	Index        int         `json:"index,omitempty"`
	Height       int         `json:"height,omitempty"`
	Width        int         `json:"width,omitempty"`
	Meta         interface{} `json:"meta,omitempty"`
	LocationName string      `json:"location_name,omitempty"`
	LocationLat  int         `json:"location_lat,omitempty"`
	LocationLong int         `json:"location_long,omitempty"`
}

type CreateConversationRequest struct {
	ChatroomID            interface{}              `json:"chatroom_id"`
	Text                  string                   `json:"text"`
	PollType              *int32                   `json:"poll_type,omitempty"`
	AllowAddOption        bool                     `json:"allow_add_option,omitempty"`
	ExpiryTime            int64                    `json:"expiry_time,omitempty"`
	Polls                 []PollObject             `json:"polls,omitempty"`
	MultilpleSelectState  *int64                   `json:"multiple_select_state,omitempty"`
	MultilpleSelectNo     int64                    `json:"multiple_select_no,omitempty"`
	AttachmentCount       int64                    `json:"attachment_count,omitempty"`
	RepliedConversationId interface{}              `json:"replied_conversation_id,omitempty"`
	RepliedChatroomID     string                   `json:"replied_chatroom_id,omitempty"`
	InternalLink          string                   `json:"internal_link,omitempty"`
	Preview               ConversationPreview      `json:"preview,omitempty"`
	IsAnonymous           bool                     `json:"is_anonymous,omitempty"`
	State                 int32                    `json:"state"`
	HasFiles              bool                     `json:"has_files,omitempty"`
	TemporaryID           string                   `json:"temporary_id,omitempty"`
	OGTags                interface{}              `json:"og_tags,omitempty"`
	ShareLink             string                   `json:"share_link,omitempty"`
	Attachments           []ConversationAttachment `json:"attachments,omitempty"`
	Metadata              interface{}              `json:"metadata,omitempty"`
	TriggerBot            bool                     `json:"trigger_bot,omitempty"`
}

type EditConversationRequest struct {
	ConversationID interface{} `json:"conversation_id" binding:"required"`
	Text           string      `json:"text" binding:"required"`
	ShareLink      string      `json:"share_link,omitempty"`
	Metadata       interface{} `json:"metadata,omitempty"`
}

type DeleteConversationRequest struct {
	ConversationIDs []interface{} `json:"conversation_ids" binding:"required"`
	TagID           int64         `json:"tag_id"`
	Reason          string        `json:"reason"`
}

// CreateConversation is used to create a new conversation in chatroom
func CreateConversation(c *gin.Context) {
	Conversation(c, utils.POSTMethod)
}

// EditConversation is used to edit a specific conversation
func EditConversation(c *gin.Context) {
	Conversation(c, utils.PUTMethod)
}

// GetConversation is used to get conversation
func GetConversation(c *gin.Context) {
	Conversation(c, utils.GETMethod)
}

// DeleteConversation is used to delete conversation
func DeleteConversation(c *gin.Context) {
	Conversation(c, utils.DELETEMethod)
}

// Conversation method handles conversation object
func Conversation(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}
	deviceID := user.GetRequestingUserDeviceId(c)

	//Send request
	switch method {
	case utils.GETMethod:

		getConversationInternal(c, userId)

	case utils.POSTMethod:

		createConversationInternal(c, userId, deviceID)

	case utils.PUTMethod:

		editConversationInternal(c, userId)

	case utils.DELETEMethod:

		deleteConversationInternal(c, userId)
	}
}

func parseCreateConversationRequest(c *gin.Context) (*CreateConversationRequest, error) {
	//POST body params
	var ccr CreateConversationRequest

	if err := c.ShouldBindJSON(&ccr); err != nil {
		return nil, err
	}

	if ccr.ChatroomID != nil {
		ccr.ChatroomID = utils.ParseInterfaceToString(ccr.ChatroomID)
	}

	if ccr.RepliedConversationId != nil {
		ccr.RepliedConversationId = utils.ParseInterfaceToString(ccr.RepliedConversationId)
	}

	return &ccr, nil
}

func parseEditConversationRequest(c *gin.Context) (*EditConversationRequest, error) {
	//POST body params
	var ecr EditConversationRequest

	if err := c.ShouldBindJSON(&ecr); err != nil {
		return nil, err
	}

	if ecr.Metadata != nil {
		metadataString, _ := json.Marshal(ecr.Metadata)

		if metadataString != nil {
			ecr.Metadata = string(metadataString)
		}
	}

	ecr.ConversationID = utils.ParseInterfaceToString(ecr.ConversationID)

	return &ecr, nil
}

func parseDeleteConversationRequest(c *gin.Context) (*DeleteConversationRequest, error) {
	//POST body params
	var dcr DeleteConversationRequest

	if err := c.ShouldBindJSON(&dcr); err != nil {
		return nil, err
	}

	// parse conversation ids to string
	for i := 0; i < len(dcr.ConversationIDs); i++ {
		dcr.ConversationIDs[i] = utils.ParseInterfaceToString(dcr.ConversationIDs[i])
	}

	return &dcr, nil
}

func getConversationInternal(c *gin.Context, userId string) {

	//GET Request params
	meta := c.Query(ParamMeta)

	if meta == "" || meta == "false" {
		//If meta is missing, call api/conversation/fetch api internally
		params := map[string]string{
			ParamChatroomId:                 c.Query(ParamChatroomId),
			ParamConversationId:             c.Query(ParamConversationId),
			ParamPaginateBy:                 c.Query(ParamPaginateBy),
			ParamScrollDirection:            c.Query(ParamScrollDirection),
			ParamIncludeConversationId:      c.Query(ParamIncludeConversationId),
			ParamExcludedConversationStates: c.Query(ParamExcludedConversationStates),
		}

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchConversationEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, true)

	} else {
		//else, call api/conversation_meta api internally
		params := map[string]string{
			ParamChatroomId:     c.Query(ParamChatroomId),
			ParamConversationId: c.Query(ParamConversationId),
		}
		//Send Request
		utils.SendRequest(c, utils.CoreService, ConversationMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	}
}

func createConversationInternal(c *gin.Context, userId string, deviceID string) {

	//Body to be sent in the create conversation api internally
	createConversationRequest, err := parseCreateConversationRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, CreateConversationEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createConversationRequest)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR != nil {
		utils.GenerateResponse(c, apiCR.Response, true)
		if apiCR.Success == true {
			//headers with userID and chatroomID
			headers := utils.CreateHeadersFromToken(c, userId, deviceID)
			redisClient := utils.GetRedisClientFromContext(c)

			go parseAndPublishConversationOnTopicTypeChatroom(redisClient, headers, createConversationRequest.ChatroomID, apiCR.Response)
			go parseAndPublishConversationOnTopicTypeCommunity(redisClient, headers, apiCR.Response)
		}
	}
}

// parseAndPublishConversationOnTopicTypeCommunity to publish Conversation on TopicTypeCommunityDynamic
func parseAndPublishConversationOnTopicTypeChatroom(redisClient *redis.Client, headers map[string]interface{}, chatroomID interface{}, response map[string]interface{}) {
	if chatroomID == nil {
		logging.Error("parseAndPublishConversationOnTopicTypeChatroom: chatroom ID is missing")
		return
	}
	if response == nil {
		logging.Error("parseAndPublishConversationOnTopicTypeChatroom: response is missing")
		return
	}

	chatroomIDStr := fmt.Sprintf("%v", chatroomID)
	// Get the chatroom data
	chatroomData, _, err := getChatroomInternal(redisClient, headers, chatroomIDStr)
	if err != nil {
		logging.Error(fmt.Sprintf("Error fetching chatroom data for chatroomID %v: %v", chatroomIDStr, err))
		return
	}

	// Check if the chatroom is secret or not
	isSecret, ok := chatroomData["is_secret"].(bool)
	if !ok {
		logging.Error("parseAndPublishConversationOnTopicTypeChatroom: is_secret key is missing or is not a valid bool in chatroomData")
		isSecret = false // Default to false if not specified, or handle as needed
	}

	// Get total participants count
	totalParticipantsCount, err := getTotalParticipantsInternal(redisClient, headers, chatroomIDStr, isSecret)
	if err != nil {
		logging.Error(fmt.Sprintf("Error fetching total participants count for chatroomID %v: %v", chatroomIDStr, err))
		return
	}

	// Add total_participants_count to the response
	response["total_participants_count"] = totalParticipantsCount

	// Publish the conversation with updated response
	pubsubPublish.PublishConversationOnTopicTypeChatroom(headers, chatroomID, response)
}

// parseAndPublishConversationOnTopicTypeCommunity to publish Conversation on TopicTypeCommunityDynamic
func parseAndPublishConversationOnTopicTypeCommunity(redisClient *redis.Client, headers map[string]interface{}, response map[string]interface{}) {
	if response == nil {
		logging.Error("parseAndPublishConversationOnTopicTypeCommunity: response is missing")
		return
	}
	// Check if "conversation" exists and is a map
	conversation, ok := response["conversation"].(map[string]interface{})
	if !ok || conversation == nil {
		logging.Error("parseAndPublishConversationOnTopicTypeCommunity: conversation key is missing or is not a valid map in response")
		return // Exit if "conversation" is missing or invalid
	}

	// Check if "chatroom_id" exists within "conversation"
	chatroomIDFloat, ok := conversation["chatroom_id"].(float64)
	if !ok || chatroomIDFloat == 0 {
		logging.Error("parseAndPublishConversationOnTopicTypeCommunity: chatroom_id key is missing or is not a valid float in conversation")
		return // Exit if "chatroom_id" is missing or invalid
	}

	// Convert chatroom_id to string format
	chatroomID := fmt.Sprintf("%.0f", chatroomIDFloat)
	apiCR, _, err := getChatroomInternal(redisClient, headers, chatroomID)
	if apiCR != nil {
		isSecret, okSecret := apiCR["is_secret"].(bool)
		chatroomType, okChatroomType := apiCR["type"].(float64)

		if (okSecret && isSecret == true) || (okChatroomType && chatroomType == chatroom.DMChatroomType) {
			allParticipantIDs, err := getParticipantsInternal(redisClient, headers, chatroomID, isSecret)
			if allParticipantIDs != nil {
				response["participants"] = allParticipantIDs
				pubsubPublish.PublishConversationOnTopicTypeCommunity(headers, response)
			} else {
				logging.Error(fmt.Sprintf("Error in getting participants data before publishing: %v", err))
			}
		} else {
			pubsubPublish.PublishConversationOnTopicTypeCommunity(headers, response)
		}

	} else {
		logging.Error(fmt.Sprintf("Error in getting chatroom data before publishing: %v", err))
	}
}

func editConversationInternal(c *gin.Context, userId string) {

	//Body to be sent in the edit conversation api internally
	editConversationRequest, err := parseEditConversationRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Get Request response
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, EditConversationEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, editConversationRequest)
	if respBytes == nil {
		return
	}

	//Parse and generate response
	utils.ParseResponse(c, respBytes, statusCode, true)
}

func deleteConversationInternal(c *gin.Context, userId string) {

	//Body to be sent in the delete conversation api internally
	deleteConversationRequest, err := parseDeleteConversationRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, DeleteConversationEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, deleteConversationRequest)
}

func getChatroomInternal(redisClient *redis.Client, headers map[string]interface{}, chatroomID string) (map[string]interface{}, int, error) {
	// Params to be sent in the api/chatroom/fetch request
	chatroomParams := map[string]string{
		ParamChatroomId: chatroomID,
	}

	//Custom headers since this API will be called after conversation create and headers between these two APIs can have different x-api-version
	headers[utils.HeadersApiVersion] = ChatroomAPIVersion

	// Check if the chatroom is present in the cache first
	cachedChatroom, err := getChatroomFromCache(redisClient, chatroomID)
	if err == nil && cachedChatroom != nil {
		// Return the cached chatroom data
		return cachedChatroom, http.StatusOK, nil
	}

	//Get Request response
	respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.CoreService, chatroom.FetchChatroomEndPoint, utils.GETRequest, headers, chatroomParams, nil)
	//Parse and generate response
	apiCR := utils.ValidateClientResponseWithoutContext(respBytes, statusCode, err)
	if apiCR != nil {
		chatroomAPICR, ok := apiCR["chatroom"].(map[string]interface{})
		if !ok || chatroomAPICR == nil {
			err := fmt.Errorf("getChatroomInternal: chatroom key is missing or is not a valid map in apiCR")
			logging.Error(err)
			return nil, statusCode, err
		}

		// Save the fetched chatroom data in the cache
		if err := saveChatroomInCache(redisClient, chatroomID, chatroomAPICR); err != nil {
			logging.Error(fmt.Sprintf("Error saving chatroom data to cache: %v", err))
		}

		return chatroomAPICR, statusCode, err
	}
	return nil, statusCode, err
}

// SaveChatroomInCache saves the chatroom data in Redis with a specified TTL.
func saveChatroomInCache(redisClient *redis.Client, chatroomID string, chatroomData map[string]interface{}) error {
	// Serialize the chatroom data to JSON
	data, err := json.Marshal(chatroomData)
	if err != nil {
		return fmt.Errorf("failed to marshal chatroom data: %v", err)
	}

	// Cache key for the chatroom data
	cacheKey := fmt.Sprintf(cache.ChatroomKey, chatroomID)

	// Save to Redis with a TTL of 24 hours (can be adjusted as needed)
	err = cache.Set(redisClient, cacheKey, data, time.Hour*cache.ChatroomTTL)
	if err != nil {
		return fmt.Errorf("error saving chatroom data to cache: %v", err)
	}

	return nil
}

// GetChatroomFromCache fetches the chatroom data from the cache.
func getChatroomFromCache(redisClient *redis.Client, chatroomID string) (map[string]interface{}, error) {
	// Cache key for the chatroom data
	cacheKey := fmt.Sprintf(cache.ChatroomKey, chatroomID)

	// Get data from Redis
	cacheValue, _, err := cache.Get(redisClient, cacheKey)
	if err != nil {
		return nil, err
	}
	if cacheValue == "" {
		return nil, fmt.Errorf("no data found in cache for chatroom: %s", chatroomID)
	}

	// Parse the cached data into a map[string]interface{}
	var chatroomData map[string]interface{}
	err = json.Unmarshal([]byte(cacheValue), &chatroomData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chatroom data: %v", err)
	}

	return chatroomData, nil
}

func getParticipantsInternal(redisClient *redis.Client, headers map[string]interface{}, chatroomID string, isSecret bool) ([]string, error) {

	// Check if the participants are already in the Redis cache
	cacheKey := fmt.Sprintf(cache.ChatroomParticipantsKey, chatroomID)
	cachedParticipantIDs, err := getParticipantsFromCache(redisClient, cacheKey)
	if err != nil {
		return nil, err
	}
	if cachedParticipantIDs != nil && len(cachedParticipantIDs) > 0 {
		// If cache exists and is not nil, return the participants from cache
		logging.Info("Returning participants from Redis cache")
		return cachedParticipantIDs, nil
	} else {
		// Initialize parameters for pagination and collection of participants
		params := map[string]string{
			ParamChatroomId: chatroomID,
			ParamPage:       ChatroomParticipantsPage,
			ParamPageSize:   ChatroomParticipantsPageSize,
		}
		var allParticipantIDs []string

		// Select the correct API endpoint based on whether the chatroom is secret
		var endpoint string
		if isSecret {
			endpoint = chatroom.FetchSecretParticipantsMetaEndPoint
		} else {
			endpoint = chatroom.FetchParticipantsMetaEndPoint
		}
		// Loop to fetch participants until the response is empty
		for {
			//Custom headers since this API will be called after conversation create and headers between these two APIs can have different x-api-version
			headers[utils.HeadersPlatformCode] = ChatroomPlatformCode
			headers[utils.HeadersVersionCode] = ChatroomVersionCode
			headers[utils.HeadersApiVersion] = ChatroomParticipantsAPIVersion
			// Make the API call to fetch participants
			respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(
				utils.CoreService,
				endpoint,
				utils.GETRequest,
				headers,
				params,
				nil,
			)
			// Check if the response is empty or if there was an error
			if err != nil || respBytes == nil || statusCode != http.StatusOK {
				break
			}

			// Parse the response to extract participant IDs
			var response struct {
				TotalParticipantsCount int `json:"total_participants_count"`
				Participants           []struct {
					ID string `json:"uuid"`
				} `json:"participants"`
			}
			if err := json.Unmarshal(respBytes, &response); err != nil {
				break
			}
			// If no participants were found, exit the loop
			if len(response.Participants) == 0 {
				break
			}

			// Collect participant IDs
			for _, participant := range response.Participants {
				allParticipantIDs = append(allParticipantIDs, participant.ID)
			}
			if len(allParticipantIDs) == response.TotalParticipantsCount {
				break
			}

			// Increment page for the next request
			currentPage, _ := strconv.Atoi(params[ParamPage])
			params[ParamPage] = strconv.Itoa(currentPage + 1)
		}

		// Update Redis cache with the collected participant IDs if any participants were found
		if len(allParticipantIDs) > 0 {
			err = setParticipantsInCache(redisClient, cacheKey, allParticipantIDs)
			if err != nil {
				return nil, err
			} else {
				return allParticipantIDs, nil
			}
		}
	}
	return nil, nil
}

// getParticipantsFromCache fetches data from Redis by key and unmarshals it into a []string
func getParticipantsFromCache(redisClient *redis.Client, cacheKey string) ([]string, error) {
	// Get data from Redis
	cacheValue, _, err := cache.Get(redisClient, cacheKey)
	if err != nil {
		return nil, err
	}
	if cacheValue == "" {
		return []string{}, nil
	}

	// Parse the cached data into a slice of strings (assuming JSON array format)
	var participantIDs []string
	err = json.Unmarshal([]byte(cacheValue), &participantIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached data: %v", err)
	}

	return participantIDs, nil
}

// setParticipantsInCache to marshal data and set in cache
func setParticipantsInCache(redisClient *redis.Client, cacheKey string, allParticipantIDs []string) error {
	// Serialize the value to JSON
	data, err := json.Marshal(allParticipantIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	// Store the data in Redis with an expiration time (e.g., 24 hour)
	err = cache.Set(redisClient, cacheKey, data, time.Hour*cache.ChatroomParticipantsTTL)
	if err != nil {
		return fmt.Errorf("error updating Redis cache: %v", err)
	}
	return nil
}

// getTotalParticipantsInternal returns the total participants count, fetching the first page if not present in cache.
func getTotalParticipantsInternal(redisClient *redis.Client, headers map[string]interface{}, chatroomID string, isSecret bool) (int, error) {

	// Try to get total participants count from the cache
	totalCount, err := getTotalParticipantsCountFromCache(redisClient, chatroomID)
	if err == nil {
		return totalCount, nil
	}

	// Initialize parameters for pagination to fetch the first page
	params := map[string]string{
		ParamChatroomId: chatroomID,
		ParamPage:       ChatroomParticipantsPage,
		ParamPageSize:   ChatroomParticipantsPageSize,
	}

	// Select the correct API endpoint based on whether the chatroom is secret
	var endpoint string
	if isSecret {
		endpoint = chatroom.FetchSecretParticipantsMetaEndPoint
	} else {
		endpoint = chatroom.FetchParticipantsMetaEndPoint
	}

	// Custom headers for API call
	headers[utils.HeadersPlatformCode] = ChatroomPlatformCode
	headers[utils.HeadersVersionCode] = ChatroomVersionCode
	headers[utils.HeadersApiVersion] = ChatroomParticipantsAPIVersion

	// Make the API call to fetch the first page of participants
	respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(
		utils.CoreService,
		endpoint,
		utils.GETRequest,
		headers,
		params,
		nil,
	)
	if err != nil || respBytes == nil || statusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch participants data: %v", err)
	}

	// Parse the response to get the total participants count
	var response struct {
		TotalParticipantsCount int `json:"total_participants_count"`
		Participants           []struct {
			ID string `json:"uuid"`
		} `json:"participants"`
	}
	if err := json.Unmarshal(respBytes, &response); err != nil {
		return 0, fmt.Errorf("failed to unmarshal participants data: %v", err)
	}

	// Save total participants count in cache
	if err := saveTotalParticipantsCountInCache(redisClient, chatroomID, response.TotalParticipantsCount); err != nil {
		return 0, fmt.Errorf("error saving total participants count to cache: %v", err)
	}

	// Return the total participants count
	return response.TotalParticipantsCount, nil
}

// getTotalParticipantsCountFromCache fetches the total participants count from Redis cache
func getTotalParticipantsCountFromCache(redisClient *redis.Client, chatroomID string) (int, error) {
	// Cache key for storing total participants count
	totalParticipantsCacheKey := fmt.Sprintf(cache.ChatroomTotalParticipantsKey, chatroomID)

	// Try to get the total count from the cache
	cachedTotalCount, _, err := cache.Get(redisClient, totalParticipantsCacheKey)
	if err != nil || cachedTotalCount == "" {
		return 0, fmt.Errorf("total participants count not found in cache")
	}

	// Convert the cached value to an integer and return it
	totalCount, err := strconv.Atoi(cachedTotalCount)
	if err != nil {
		return 0, fmt.Errorf("failed to convert cached total participants count: %v", err)
	}
	return totalCount, nil
}

// saveTotalParticipantsCountInCache saves the total participants count to Redis cache
func saveTotalParticipantsCountInCache(redisClient *redis.Client, chatroomID string, totalCount int) error {
	// Cache key for storing total participants count
	totalParticipantsCacheKey := fmt.Sprintf(cache.ChatroomTotalParticipantsKey, chatroomID)

	// Save the total count to Redis with an expiration time (e.g., 24 hours)
	err := cache.Set(redisClient, totalParticipantsCacheKey, strconv.Itoa(totalCount), time.Hour*cache.ChatroomParticipantsTTL)
	if err != nil {
		return fmt.Errorf("error saving total participants count to cache: %v", err)
	}
	return nil
}
