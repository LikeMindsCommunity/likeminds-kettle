package pubsubSubscribe

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/internal/cache"
	"github.com/nateshr/likeminds-authentication/internal/handlers/pubsubCommon"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nateshr/likeminds-authentication/internal/handlers/channel"
	"github.com/nateshr/likeminds-authentication/internal/handlers/chatroom"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
	"github.com/nateshr/likeminds-authentication/internal/utils/api_client"
)

var upgrader = newUpgrader()

// newUpgrader creates a new websocket Upgrader
func newUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  pubsubCommon.ReadBufferSizeDefault,
		WriteBufferSize: pubsubCommon.WriteBufferSizeDefault,
	}
}

// Subscribe to open WS against a topic
func Subscribe(c *gin.Context) {
	// Validate params and headers before subscribing to a topic
	topicSplit, headers, err := validateParamsAndHeaders(c)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// check if topic is supported or not
	if !isSubscribeTopicSupported(topicSplit) {
		return
	}

	// Check if user has access to chatroom
	statusCode, err := hasAccessToChatroom(topicSplit, headers)
	if err != nil {
		switch statusCode {
		case http.StatusBadRequest:
			utils.GeneralBadRequestError(c, err.Error())
		default:
			utils.GeneralAPIError(c, err.Error())
		}
		return
	}

	// Upgrade HTTP request
	conn, err := upgraderHTTPToWs(c)
	if err != nil {
		updatedErr := fmt.Sprintf(pubsubCommon.ErrorFailedUpgrader, err)
		logging.Error(updatedErr)
		utils.GeneralAPIError(c, updatedErr)
		return
	}

	// Connect to the websocket server
	serverConn, err := dialToWs(c.Param(pubsubCommon.ParamTopic), headers)
	if err != nil {
		updatedErr := fmt.Sprintf(pubsubCommon.ErrorFailedDial, err)
		logging.Error(updatedErr)
		utils.GeneralAPIError(c, updatedErr)
		return
	}

	logging.Info(pubsubCommon.WsConnectionEstablished)

	// Handle communication between the client and the websocket server
	redisClient := utils.GetRedisClientFromContext(c)
	utils.SafeGo(func() { readFromClientAndWriteToServer(conn, serverConn, redisClient, headers, topicSplit) })
	utils.SafeGo(func() { readFromServerAndWriteToClient(conn, serverConn) })
}

// validateParamsAndHeaders to validate params and headers sent while subscribing to topic
func validateParamsAndHeaders(c *gin.Context) ([]string, map[string]interface{}, error) {
	//Validate topic
	topic := c.Param(pubsubCommon.ParamTopic)
	topicSplit, err := getTopicSplit(topic)
	if err != nil {
		return nil, nil, err
	}
	var topicID string
	if len(topicSplit) > 1 {
		topicID = topicSplit[1]
	}
	//If topicID is missing, return error
	if topicID == "" || topicID == "null" {
		return nil, nil, errors.New(pubsubCommon.ErrorTopicIDMissing)
	}

	userID := user.GetRequestingUserId(c)
	//If userID is missing, return error
	if userID == "" || userID == "null" {
		return nil, nil, errors.New(pubsubCommon.ErrorUserUUIDMissing)
	}

	deviceID := user.GetRequestingUserDeviceId(c)

	return topicSplit, utils.CreateHeadersFromToken(c, userID, deviceID), nil
}

// getTopicSplit will decode topic and return split
func getTopicSplit(topic string) ([]string, error) {
	if topic == "" || topic == "null" {
		return nil, errors.New(pubsubCommon.ErrorTopicMissing)
	}
	topicSplit := strings.Split(topic, ":")
	if len(topicSplit) <= 1 {
		return nil, errors.New(pubsubCommon.ErrorTopicInvalid)
	}
	return topicSplit, nil
}

func isSubscribeTopicSupported(topicSplit []string) bool {
	switch topicSplit[0] {
	case pubsubCommon.TopicTypeChatroom:
		return true
	case pubsubCommon.TopicTypeCommunity:
		return true
	default:
		return false
	}
}

// hasAccessToChatroom to check if userID has access to topicID while subscribing to chatroom topic
func hasAccessToChatroom(topicSplit []string, headers map[string]interface{}) (int, error) {
	switch topicSplit[0] {
	case pubsubCommon.TopicTypeChatroom:
		params := map[string]string{
			channel.ParamChannelId: topicSplit[1],
		}
		//Get chatroom details to verify if user has access to any chatroom / cohort based chatroom / secret chatroom
		respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.CoreService, channel.SyncChannelDetailEndppoint, utils.GETRequest, headers, params, nil)
		if err != nil || statusCode != 200 {
			return statusCode, errors.New(fmt.Sprintf(pubsubCommon.ErrorUserChatroomAccess, err))
		} else {
			var chatroomDetailParentResponse chatroom.ChatroomDetailParentResponse
			if err := json.Unmarshal(respBytes, &chatroomDetailParentResponse); err != nil {
				return statusCode, errors.New(fmt.Sprintf(pubsubCommon.ErrorUnmarshalErrorJson, err))
			}
			chatroomDetailArray := chatroomDetailParentResponse.ChatroomDetail
			if len(chatroomDetailArray) < 1 {
				return statusCode, errors.New(pubsubCommon.ErrorChatroomResponseInvalid)
			}
			//If user has access to secret chatroom
			chatroomDetail := chatroomDetailArray[0]
			canAccessSecretChatroom := chatroomDetail.CanAccessSecretChatroom
			if canAccessSecretChatroom != nil {
				if *canAccessSecretChatroom == false {
					return statusCode, errors.New(pubsubCommon.ErrorUserChatroomAccess)
				}
			}
			//If user has access to cohort based chatroom
			cohortAccess := chatroomDetail.CohortAccess
			if cohortAccess != nil {
				if *cohortAccess != 200 {
					return statusCode, errors.New(pubsubCommon.ErrorUserChatroomAccess)
				}
			}
		}
	}
	return http.StatusOK, nil
}

// upgraderHTTPToWs to upgrade the incoming HTTP request to a WebSocket connection
func upgraderHTTPToWs(c *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	return conn, err
}

// dialToWs to dial to websocket server
func dialToWs(topic string, headers map[string]interface{}) (*websocket.Conn, error) {
	psURL := api_client.GetPandemoniumServiceWsUrl()
	updatedPsURL := fmt.Sprintf("%s/subscribe/%s", psURL, topic)
	serverConn, _, err := websocket.DefaultDialer.Dial(updatedPsURL, createHeaders(headers))
	return serverConn, err
}

// createHeaders to createHeaders required for connecting with websocket server
func createHeaders(headersMap map[string]interface{}) http.Header {
	out := http.Header{}

	for key, value := range headersMap {
		out.Add(key, value.(string))
	}
	return out
}

func readFromClientAndWriteToServer(conn *websocket.Conn, serverConn *websocket.Conn, redisClient *redis.Client, headers map[string]interface{}, topicSplit []string) {
	defer func() {
		disconnect(conn)
		disconnect(serverConn)
	}()

	utils.SafeGo(func() { startPingMessageToServer(serverConn) })

	serverConn.SetPongHandler(func(string) error {
		log.Println(pubsubCommon.PongReceivedWs)
		return nil
	})

	for {
		readMessageType, readMessagePayload, err := conn.ReadMessage()
		if err != nil {
			logging.Error(fmt.Sprintf(pubsubCommon.ErrorReadClientWs, err))
			return
		}
		logging.Info(fmt.Sprintf(pubsubCommon.ReceivedMessageClientWs, readMessageType))

		var readMessageJsonMap map[string]interface{}
		err = json.Unmarshal(readMessagePayload, &readMessageJsonMap)
		if err != nil {
			logging.Error(pubsubCommon.ErrorInvalidJSONFormat, err)
			return
		}

		topicMessageType := readMessageJsonMap[pubsubCommon.ParamTopicMessageType]

		switch topicSplit[0] {
		case pubsubCommon.TopicTypeChatroom:
			switch topicMessageType {
			case pubsubCommon.TopicMessageTypeCreateConversationRequest:
				updatedMessageRequest, err := getUpdatedMessageRequest(redisClient, headers, topicSplit, readMessageJsonMap)
				updatedMessageRequestBytes, err := json.Marshal(updatedMessageRequest)
				if err != nil {
					fmt.Errorf("failed to marshal chatroom data: %v", err)
					return
				}
				// Forward the message to the WebSocket server
				err = serverConn.WriteMessage(readMessageType, updatedMessageRequestBytes)
				if err != nil {
					logging.Error(fmt.Sprintf(pubsubCommon.ErrorWriteServerWs, err))
					return
				}
			default:
				// Forward the message to the WebSocket server
				err = serverConn.WriteMessage(readMessageType, readMessagePayload)
				if err != nil {
					logging.Error(fmt.Sprintf(pubsubCommon.ErrorWriteServerWs, err))
					return
				}
			}
		}
	}
}

func readFromServerAndWriteToClient(conn *websocket.Conn, serverConn *websocket.Conn) {
	defer func() {
		disconnect(conn)
		disconnect(serverConn)
	}()

	utils.SafeGo(func() { startPingMessageToClient(conn) })

	conn.SetPongHandler(func(string) error {
		log.Println(pubsubCommon.PongReceivedClient)
		return nil
	})

	for {
		messageType, msg, err := serverConn.ReadMessage()
		if err != nil {
			logging.Error(fmt.Sprintf(pubsubCommon.ErrorReadServerWs, err))
			return
		}
		logging.Info(fmt.Sprintf(pubsubCommon.ReceivedMessageServerWs, messageType))

		// Forward the message to the client
		err = conn.WriteMessage(messageType, msg)
		if err != nil {
			logging.Error(fmt.Sprintf(pubsubCommon.ErrorWriteClientWs, err))
			return
		}
	}
}

func disconnect(conn *websocket.Conn) {
	logging.Info(pubsubCommon.ConnectionClosed)
	err := conn.Close()
	if err != nil {
		log.Println(pubsubCommon.ErrorUnableToCloseWs, err)
		return
	}
}

func startPingMessageToClient(conn *websocket.Conn) {
	// Start a goroutine to send pings periodically to the client
	for {
		time.Sleep(pubsubCommon.PingPeriod) // Interval between pings
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf(fmt.Sprintf(pubsubCommon.ErrorPingSendClient, err))
			return
		}
		log.Println(pubsubCommon.PingSendClient)
	}
}

func startPingMessageToServer(conn *websocket.Conn) {
	// Start a goroutine to send pings periodically to the client
	for {
		time.Sleep(pubsubCommon.PingPeriod) // Interval between pings
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf(fmt.Sprintf(pubsubCommon.ErrorPingSendWs, err))
			return
		}
		log.Println(pubsubCommon.PingSendWs)
	}
}

// getUpdatedMessageRequest to get updated payload to be sent in case of message.create.request
func getUpdatedMessageRequest(redisClient *redis.Client, headers map[string]interface{}, topicSplit []string, readMessageJsonMap map[string]interface{}) (map[string]interface{}, error) {
	chatroomID := topicSplit[0]

	apiCR, _, err := getChatroomInternal(redisClient, headers, chatroomID)
	if apiCR != nil {
		isSecret, _ := apiCR["is_secret"].(bool)
		allParticipantIDs, err := getParticipantsInternal(redisClient, headers, chatroomID, isSecret)
		if allParticipantIDs != nil {
			readMessageJsonMap["participants"] = allParticipantIDs
			return readMessageJsonMap, nil
		} else {
			return nil, fmt.Errorf("error in getting participants data before publishing: %v", err)
		}

	} else {
		return nil, fmt.Errorf("error in getting chatroom data before publishing: %v", err)
	}
}

func getChatroomInternal(redisClient *redis.Client, headers map[string]interface{}, chatroomID string) (map[string]interface{}, int, error) {
	// Params to be sent in the api/chatroom/fetch request
	chatroomParams := map[string]string{
		pubsubCommon.ParamChatroomId: chatroomID,
	}

	//Custom headers since this API will be called after conversation create and headers between these two APIs can have different x-api-version
	headers[utils.HeadersApiVersion] = pubsubCommon.ChatroomAPIVersion

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

func getParticipantsInternal(redisClient *redis.Client, headers map[string]interface{}, chatroomID string, isSecret bool) ([]string, error) {

	// Check if the participants are already in the Redis cache
	chatroomCacheKey := fmt.Sprintf(cache.ChatroomParticipantsKey, chatroomID)
	participantsCacheValue, err := getParticipantsFromCache(redisClient, chatroomCacheKey)
	if err != nil {
		return nil, err
	}
	if participantsCacheValue != nil && len(participantsCacheValue) > 0 {
		// If cache exists and is not nil, return the participants from cache
		logging.Info("Returning participants from Redis cache")
		return participantsCacheValue, nil
	} else {
		// Initialize parameters for pagination and collection of participants
		participantsAPIParams := map[string]string{
			pubsubCommon.ParamChatroomId: chatroomID,
			pubsubCommon.ParamPage:       pubsubCommon.ChatroomParticipantsAPIPage,
			pubsubCommon.ParamPageSize:   pubsubCommon.ChatroomParticipantsAPIPageSize,
		}
		var allParticipantIDs []string

		// Select the correct API participantsAPIEndpoint based on whether the chatroom is secret
		var participantsAPIEndpoint string
		if isSecret {
			participantsAPIEndpoint = chatroom.FetchSecretParticipantsMetaEndPoint
		} else {
			participantsAPIEndpoint = chatroom.FetchParticipantsMetaEndPoint
		}
		// Loop to fetch participants until the response is empty
		for {
			//Custom headers since this API will be called after conversation create and headers between these two APIs can have different x-api-version
			headers[utils.HeadersPlatformCode] = pubsubCommon.ChatroomParticipantsAPIPlatformCode
			headers[utils.HeadersVersionCode] = pubsubCommon.ChatroomParticipantsAPIVersion
			headers[utils.HeadersApiVersion] = pubsubCommon.ChatroomParticipantsAPIVersion
			// Make the API call to fetch participants
			participantsAPIResponseBytes, statusCode, err := utils.GetRequestResponseWithoutContext(
				utils.CoreService,
				participantsAPIEndpoint,
				utils.GETRequest,
				headers,
				participantsAPIParams,
				nil,
			)
			// Check if the participantsAPIResponse is empty or if there was an error
			if err != nil || participantsAPIResponseBytes == nil || statusCode != http.StatusOK {
				break
			}

			// Parse the participantsAPIResponse to extract participant IDs
			var participantsAPIResponse struct {
				TotalParticipantsCount int `json:"total_participants_count"`
				Participants           []struct {
					ID string `json:"uuid"`
				} `json:"participants"`
			}
			if err := json.Unmarshal(participantsAPIResponseBytes, &participantsAPIResponse); err != nil {
				break
			}
			// If no participants were found, exit the loop
			if len(participantsAPIResponse.Participants) == 0 {
				break
			}

			// Collect participant IDs
			for _, participant := range participantsAPIResponse.Participants {
				allParticipantIDs = append(allParticipantIDs, participant.ID)
			}
			if len(allParticipantIDs) == participantsAPIResponse.TotalParticipantsCount {
				break
			}

			// Increment page for the next request
			currentPage, _ := strconv.Atoi(participantsAPIParams[pubsubCommon.ParamPage])
			participantsAPIParams[pubsubCommon.ParamPage] = strconv.Itoa(currentPage + 1)
		}

		// Update Redis cache with the collected participant IDs if any participants were found
		if len(allParticipantIDs) > 0 {
			err = saveParticipantsInCache(redisClient, chatroomCacheKey, allParticipantIDs)
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
func saveParticipantsInCache(redisClient *redis.Client, cacheKey string, allParticipantIDs []string) error {
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
