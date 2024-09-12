package conversation

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/internal/cache"
	"github.com/nateshr/likeminds-authentication/internal/handlers/chatroom"
	"github.com/nateshr/likeminds-authentication/internal/handlers/community"
	"github.com/nateshr/likeminds-authentication/internal/handlers/pubsub"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
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
			go publishConversationOnTopicTypeChatroom(c, createConversationRequest, userId, deviceID, apiCR.Response)
			go publishConversationOnTopicTypeCommunity(c, userId, deviceID, apiCR.Response)
		}
	}
}

// publishConversationOnTopicTypeChatroom to publish Conversation on TopicTypeChatroom
func publishConversationOnTopicTypeChatroom(c *gin.Context, createConversationRequest *CreateConversationRequest, userId string, deviceID string, response map[string]interface{}) {
	headers := utils.CreateHeadersFromToken(c, userId, deviceID)
	topicChatroom := fmt.Sprintf(pubsub.TopicTypeChatroom, createConversationRequest.ChatroomID)
	params := map[string]string{
		pubsub.ParamTopicMessageType: pubsub.TopicMessageTypeConversation,
	}

	publishDataOnPandemonium(topicChatroom, headers, params, response)
}

// publishConversationOnTopicTypeChatroom to publish Conversation on TopicTypeCommunity
func publishConversationOnTopicTypeCommunity(c *gin.Context, userId string, deviceID string, response map[string]interface{}) {
	headers := utils.CreateHeadersFromToken(c, userId, deviceID)
	var communityID = response["conversation"].(map[string]interface{})["member"].(map[string]interface{})["sdk_client_info"].(map[string]interface{})["community"]
	topicCommunity := fmt.Sprintf(pubsub.TopicTypeCommunity, communityID)
	publishParams := map[string]string{
		pubsub.ParamTopicMessageType: pubsub.TopicMessageTypeConversation,
	}

	chatroomID := fmt.Sprintf("%.0f", response["conversation"].(map[string]interface{})["chatroom_id"].(float64))
	apiCR, _, err := getChatroom(c, userId, chatroomID)
	if apiCR != nil {
		isSecret := apiCR["chatroom"].(map[string]interface{})["is_secret"].(bool)
		chatroomType := apiCR["chatroom"].(map[string]interface{})["type"].(float64)

		if isSecret == true || chatroomType == 10 {
			allParticipantIDs, err := getParticipants(c, userId, chatroomID, isSecret)
			if allParticipantIDs != nil {
				response["participants"] = allParticipantIDs
				publishDataOnPandemonium(topicCommunity, headers, publishParams, response)
			} else {
				logging.Error(fmt.Sprintf("Error in getting participants data before publishing: %v", err))
			}
		} else {
			publishDataOnPandemonium(topicCommunity, headers, publishParams, response)
		}

	} else {
		logging.Error(fmt.Sprintf("Error in getting chatroom data before publishing: %v", err))
	}
}

func publishDataOnPandemonium(topicChatroom string, headers map[string]interface{}, params map[string]string, response map[string]interface{}) {
	respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.PandemoniumService, fmt.Sprintf(PublishEndPoint, topicChatroom), utils.POSTRequestRawBody, headers, params, response)
	if err != nil || statusCode != 200 {
		logging.Error(fmt.Sprintf("Error in publishing data on pandemonium: %v | response %v", err, string(respBytes)))
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

func getChatroom(c *gin.Context, userID string, chatroomID string) (map[string]interface{}, int, error) {
	// Params to be sent in the api/chatroom/fetch request
	chatroomParams := map[string]string{
		ParamChatroomId: chatroomID,
	}

	//Custom headers since this API will be called after conversation create and headers between these two APIs can have different x-api-version
	headers := utils.CreateHeadersFromToken(c, userID, chatroomID)
	headers[utils.HeadersApiVersion] = "1"
	//Get Request response
	respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.CoreService, chatroom.FetchChatroomEndPoint, utils.GETRequest, headers, chatroomParams, nil)
	//Parse and generate response
	apiCR := utils.ValidateClientResponseWithoutContext(respBytes, statusCode, err)
	if apiCR != nil {
		return apiCR, statusCode, err
	}
	return nil, statusCode, err
}

func getParticipants(c *gin.Context, userID string, chatroomID string, isSecret bool) ([]string, error) {
	cacheKey := fmt.Sprintf(cache.ChatroomParticipantsKey, chatroomID)
	// Check if the participants are already in the Redis cache
	redisClient := utils.GetRedisClientFromContext(c)
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
			ParamPage:       "0",
			ParamPageSize:   "100",
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
			headers := utils.CreateHeadersFromToken(c, userID, chatroomID)
			headers[utils.HeadersPlatformCode] = "an"
			headers[utils.HeadersVersionCode] = "210"
			headers[utils.HeadersApiVersion] = "1"
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
			if err != nil || respBytes == nil || statusCode != 200 {
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
