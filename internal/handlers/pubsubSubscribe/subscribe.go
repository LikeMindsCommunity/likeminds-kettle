package pubsubSubscribe

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/nateshr/likeminds-authentication/internal/cache"
	"github.com/nateshr/likeminds-authentication/internal/handlers/common"
	"github.com/nateshr/likeminds-authentication/internal/handlers/common/models"
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
		ReadBufferSize:  common.ReadBufferSizeDefault,
		WriteBufferSize: common.WriteBufferSizeDefault,
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
	accessChatroomStatusCode, err := hasAccessToChatroom(topicSplit, headers)
	if err != nil {
		switch accessChatroomStatusCode {
		case http.StatusBadRequest:
			utils.GeneralBadRequestError(c, err.Error())
		default:
			utils.GeneralAPIError(c, err.Error())
		}
		return
	}

	// Upgrade HTTP request
	clientConn, err := upgraderHTTPToPandemoniumWs(c)
	if err != nil {
		logging.Error(err)
		utils.GeneralAPIError(c, err.Error())
		return
	}

	// Connect to the websocket server
	serverConn, err := dialToPandemoniumWs(c.Param(common.ParamTopic), headers)
	if err != nil {
		logging.Error(err)
		utils.GeneralAPIError(c, err.Error())
		return
	}

	logging.Info(common.WsConnectionEstablished)

	// Handle communication between the client and the websocket server
	redisClient := utils.GetRedisClientFromContext(c)
	utils.SafeGo(func() { readFromClientAndWriteToServer(clientConn, serverConn, redisClient, headers, topicSplit) })
	utils.SafeGo(func() { readFromServerAndWriteToClient(clientConn, serverConn) })
}

// validateParamsAndHeaders to validate params and headers sent while subscribing to topic
func validateParamsAndHeaders(c *gin.Context) ([]string, map[string]interface{}, error) {
	//Validate topic
	topic := c.Param(common.ParamTopic)
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
		return nil, nil, errors.New(common.ErrorTopicIDMissing)
	}

	userID := user.GetRequestingUserId(c)
	//If userID is missing, return error
	if userID == "" || userID == "null" {
		return nil, nil, errors.New(common.ErrorUserUUIDMissing)
	}

	deviceID := user.GetRequestingUserDeviceId(c)

	return topicSplit, utils.CreateHeadersFromToken(c, userID, deviceID), nil
}

// getTopicSplit will decode topic and return split
func getTopicSplit(topic string) ([]string, error) {
	if topic == "" || topic == "null" {
		return nil, errors.New(common.ErrorTopicMissing)
	}
	topicSplit := strings.Split(topic, ":")

	if len(topicSplit) <= 1 {
		return nil, errors.New(common.ErrorTopicInvalid)
	}

	return topicSplit, nil
}

func isSubscribeTopicSupported(topicSplit []string) bool {
	switch topicSplit[0] {
	case common.TopicTypeChatroom:
		return true
	case common.TopicTypeCommunity:
		return true

	default:
		return false
	}
}

// hasAccessToChatroom to check if userID has access to topicID while subscribing to chatroom topic
func hasAccessToChatroom(topicSplit []string, headers map[string]interface{}) (int, error) {
	switch topicSplit[0] {
	case common.TopicTypeChatroom:
		params := map[string]string{
			channel.ParamChannelId: topicSplit[1],
		}
		//Get chatroom details to verify if user has access to any chatroom / cohort based chatroom / secret chatroom
		accessChatroomResponseBytes, accessChatroomStatusCode, err := utils.GetRequestResponseWithoutContext(utils.CoreService, channel.SyncChannelDetailEndppoint, utils.GETRequest, headers, params, nil)
		if err != nil || accessChatroomStatusCode != 200 {
			return accessChatroomStatusCode, errors.New(fmt.Sprintf(common.ErrorUserChatroomAccess, err))
		} else {
			var chatroomDetailParentResponse chatroom.ChatroomDetailParentResponse
			if err := json.Unmarshal(accessChatroomResponseBytes, &chatroomDetailParentResponse); err != nil {
				return accessChatroomStatusCode, errors.New(fmt.Sprintf(common.ErrorUnmarshalErrorJson, err))
			}
			chatroomDetailArray := chatroomDetailParentResponse.ChatroomDetail
			if len(chatroomDetailArray) < 1 {
				return accessChatroomStatusCode, errors.New(common.ErrorChatroomResponseInvalid)
			}

			//If user has access to secret chatroom
			chatroomDetail := chatroomDetailArray[0]
			canAccessSecretChatroom := chatroomDetail.CanAccessSecretChatroom
			if canAccessSecretChatroom != nil {
				if *canAccessSecretChatroom == false {
					return accessChatroomStatusCode, errors.New(common.ErrorUserChatroomAccess)
				}
			}
			//If user has access to cohort based chatroom
			cohortAccess := chatroomDetail.CohortAccess
			if cohortAccess != nil {
				if *cohortAccess != 200 {
					return accessChatroomStatusCode, errors.New(common.ErrorUserChatroomAccess)
				}
			}
		}
	}
	return http.StatusOK, nil
}

// upgraderHTTPToPandemoniumWs to upgrade the incoming HTTP request to a WebSocket connection
func upgraderHTTPToPandemoniumWs(c *gin.Context) (*websocket.Conn, error) {
	clientConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	return clientConn, fmt.Errorf(common.ErrorFailedUpgrader, err)
}

// dialToPandemoniumWs to dial to websocket server
func dialToPandemoniumWs(topic string, headers map[string]interface{}) (*websocket.Conn, error) {
	psURL := api_client.GetPandemoniumServiceWsUrl()
	updatedPsURL := fmt.Sprintf("%s/subscribe/%s", psURL, topic)
	serverConn, _, err := websocket.DefaultDialer.Dial(updatedPsURL, createHeaderFromMap(headers))
	return serverConn, fmt.Errorf(common.ErrorFailedDial, err)
}

// createHeaderFromMap to createHeaderFromMap required for connecting with websocket server
func createHeaderFromMap(headersMap map[string]interface{}) http.Header {
	header := http.Header{}
	for key, value := range headersMap {
		header.Add(key, value.(string))
	}

	return header
}

func readFromClientAndWriteToServer(clientConn *websocket.Conn, serverConn *websocket.Conn, redisClient *redis.Client, headers map[string]interface{}, topicSplit []string) {
	defer func() {
		disconnect(clientConn)
		disconnect(serverConn)
	}()

	updateReadDeadline(clientConn)

	utils.SafeGo(func() { startPingMessage(serverConn) })

	serverConn.SetPongHandler(func(string) error {
		logging.Info(common.PongReceivedPandemoniumWs)
		updateReadDeadline(clientConn)
		return nil
	})

	for {
		readMessageFromClientType, readMessageFromClientBytes, err := clientConn.ReadMessage()
		if err != nil {
			logging.Error(fmt.Sprintf(common.ErrorReadClientWs, err))
			return
		}
		logging.Info(fmt.Sprintf(common.ReceivedMessageClientWs, readMessageFromClientType))

		var psRequest models.PSRequest
		err = json.Unmarshal(readMessageFromClientBytes, &psRequest)
		if err != nil {
			log.Printf(common.ErrorUnmarshalErrorJson, err)
			return
		}

		readMessageFromClientTMT := psRequest.TopicMessageType

		switch topicSplit[0] {
		case common.TopicTypeChatroom:
			switch readMessageFromClientTMT {
			case common.TopicMessageTypeCreateConversationRequest:
				updatedPSRequestRawDataMap, err := getUpdatedPSRequestRawDataMap(redisClient, headers, topicSplit, psRequest)
				updatedPSRequestRawDataBytes, err := json.Marshal(updatedPSRequestRawDataMap)
				if err != nil {
					logging.Error(fmt.Sprintf(common.ErrorMarshalErrorJson, err))
					return
				}
				// Forward the message to the WebSocket server
				updateWriteDeadline(serverConn)
				err = serverConn.WriteMessage(readMessageFromClientType, updatedPSRequestRawDataBytes)
				if err != nil {
					logging.Error(fmt.Sprintf(common.ErrorWriteServerWs, err))
					return
				}
			default:
				// Forward the message to the WebSocket server
				updateWriteDeadline(serverConn)
				err = serverConn.WriteMessage(readMessageFromClientType, readMessageFromClientBytes)
				if err != nil {
					logging.Error(fmt.Sprintf(common.ErrorWriteServerWs, err))
					return
				}
			}
		}
	}
}

func readFromServerAndWriteToClient(clientConn *websocket.Conn, serverConn *websocket.Conn) {
	defer func() {
		disconnect(clientConn)
		disconnect(serverConn)
	}()

	updateReadDeadline(serverConn)

	utils.SafeGo(func() { startPingMessage(clientConn) })

	clientConn.SetPongHandler(func(string) error {
		logging.Info(common.PongReceivedClient)
		updateReadDeadline(serverConn)
		return nil
	})

	for {
		messageType, msg, err := serverConn.ReadMessage()
		if err != nil {
			logging.Error(fmt.Sprintf(common.ErrorReadServerWs, err))
			return
		}
		logging.Info(fmt.Sprintf(common.ReceivedMessageServerWs, messageType))

		// Forward the message to the client
		updateWriteDeadline(clientConn)
		err = clientConn.WriteMessage(messageType, msg)
		if err != nil {
			logging.Error(fmt.Sprintf(common.ErrorWriteClientWs, err))
			return
		}
	}
}

func disconnect(conn *websocket.Conn) {
	logging.Info(common.ConnectionClosed)
	err := conn.Close()
	if err != nil {
		logging.Error(fmt.Sprintf(common.ErrorUnableToCloseWs, err))
		return
	}
}

func startPingMessage(conn *websocket.Conn) {
	// Start a goroutine to send pings periodically to the client
	for {
		time.Sleep(common.PingPeriod) // Interval between pings
		updateWriteDeadline(conn)
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			logging.Error(fmt.Sprintf(common.ErrorPingSendWs, err))
			return
		}
		logging.Info(common.PingSendWs)
	}
}

// getUpdatedPSRequestRawDataMap to get updated payload to be sent in case of message.create.request
func getUpdatedPSRequestRawDataMap(redisClient *redis.Client, headers map[string]interface{}, topicSplit []string, psRequest models.PSRequest) (map[string]interface{}, error) {
	chatroomID := topicSplit[0]

	chatroomInternalMap, _, err := getChatroomInternal(redisClient, headers, chatroomID)
	if chatroomInternalMap != nil {
		isSecret, _ := chatroomInternalMap[common.ParamIsSecret].(bool)
		allParticipantIDs, err := getParticipantsInternal(redisClient, headers, chatroomID, isSecret)
		if allParticipantIDs != nil {
			var updatedPSRequestRawDataMap map[string]interface{}
			err := json.Unmarshal(psRequest.RawData, &updatedPSRequestRawDataMap)
			if err != nil {
				return nil, err
			}
			updatedPSRequestRawDataMap[common.ParamParticipants] = allParticipantIDs
			updatedPSRequestRawDataMap[common.ParamTotalParticipantsCount] = len(allParticipantIDs)
			return updatedPSRequestRawDataMap, nil
		} else {
			return nil, err
		}

	} else {
		return nil, err
	}
}

func getChatroomInternal(redisClient *redis.Client, headers map[string]interface{}, chatroomID string) (map[string]interface{}, int, error) {
	// Params to be sent in the api/chatroom/fetch request
	chatroomAPIParams := map[string]string{
		common.ParamChatroomId: chatroomID,
	}

	//Custom headers since this API will be called after conversation create and headers between these two APIs can have different x-api-version
	headers[utils.HeadersApiVersion] = common.ChatroomAPIVersion

	// Check if the chatroom is present in the cache first
	chatroomFromCache, err := getChatroomFromCache(redisClient, chatroomID)
	if err == nil && chatroomFromCache != nil {
		// Return the cached chatroom data
		return chatroomFromCache, http.StatusOK, nil
	}

	//Get Request response
	chatroomAPIResponseBytes, chatroomAPIStatusCode, err := utils.GetRequestResponseWithoutContext(utils.CoreService, chatroom.FetchChatroomEndPoint, utils.GETRequest, headers, chatroomAPIParams, nil)
	//Parse and generate response
	chatroomAPIResponseMap := utils.ValidateClientResponseWithoutContext(chatroomAPIResponseBytes, chatroomAPIStatusCode, err)
	if chatroomAPIResponseMap != nil {
		chatroomResponseMap, ok := chatroomAPIResponseMap["chatroom"].(map[string]interface{})
		if !ok || chatroomResponseMap == nil {
			return nil, chatroomAPIStatusCode, err
		}

		// Save the fetched chatroom data in the cache
		if err := saveChatroomInCache(redisClient, chatroomID, chatroomResponseMap); err != nil {
			return nil, chatroomAPIStatusCode, err
		}

		return chatroomResponseMap, chatroomAPIStatusCode, err
	}
	return nil, chatroomAPIStatusCode, err
}

// GetChatroomFromCache fetches the chatroom data from the cache.
func getChatroomFromCache(redisClient *redis.Client, chatroomID string) (map[string]interface{}, error) {
	// Cache key for the chatroom data
	chatroomResponseKey := fmt.Sprintf(cache.ChatroomKey, chatroomID)

	// Get data from Redis
	chatroomResponseValue, _, err := cache.Get(redisClient, chatroomResponseKey)
	if err != nil {
		return nil, err
	}
	if chatroomResponseValue == "" {
		return nil, err
	}

	// Parse the cached data into a map[string]interface{}
	var chatroomResponseMap map[string]interface{}
	err = json.Unmarshal([]byte(chatroomResponseValue), &chatroomResponseMap)
	if err != nil {
		return nil, err
	}

	return chatroomResponseMap, nil
}

// SaveChatroomInCache saves the chatroom data in Redis with a specified TTL.
func saveChatroomInCache(redisClient *redis.Client, chatroomID string, chatroomResponseMap map[string]interface{}) error {
	// Serialize the chatroom chatroomResponseBytes to JSON
	chatroomResponseBytes, err := json.Marshal(chatroomResponseMap)
	if err != nil {
		return err
	}

	// Cache key for the chatroom chatroomResponseBytes
	chatroomResponseKey := fmt.Sprintf(cache.ChatroomKey, chatroomID)

	// Save to Redis with a TTL of 24 hours (can be adjusted as needed)
	err = cache.Set(redisClient, chatroomResponseKey, chatroomResponseBytes, time.Hour*cache.ChatroomTTL)
	if err != nil {
		return err
	}

	return nil
}

func getParticipantsInternal(redisClient *redis.Client, headers map[string]interface{}, chatroomID string, isSecret bool) ([]string, error) {
	// Check if the participants are already in the Redis cache
	chatroomParticipantsKey := fmt.Sprintf(cache.ChatroomParticipantsKey, chatroomID)
	chatroomParticipantsValue, err := getParticipantsFromCache(redisClient, chatroomParticipantsKey)
	if err != nil {
		return nil, err
	}
	if chatroomParticipantsValue != nil && len(chatroomParticipantsValue) > 0 {
		// If cache exists and is not nil, return the participants from cache
		return chatroomParticipantsValue, nil
	} else {
		// Initialize parameters for pagination and collection of participants
		chatroomParticipantsAPIParams := map[string]string{
			common.ParamChatroomId: chatroomID,
			common.ParamPage:       common.ChatroomParticipantsAPIPage,
			common.ParamPageSize:   common.ChatroomParticipantsAPIPageSize,
		}
		var allParticipantIDs []string

		// Select the correct API chatroomParticipantsAPIEndpoint based on whether the chatroom is secret
		var chatroomParticipantsAPIEndpoint string
		if isSecret {
			chatroomParticipantsAPIEndpoint = chatroom.FetchSecretParticipantsMetaEndPoint
		} else {
			chatroomParticipantsAPIEndpoint = chatroom.FetchParticipantsMetaEndPoint
		}
		// Loop to fetch participants until the response is empty
		for {
			//Custom headers since this API will be called after conversation create and headers between these two APIs can have different x-api-version
			headers[utils.HeadersPlatformCode] = common.ChatroomParticipantsAPIPlatformCode
			headers[utils.HeadersVersionCode] = common.ChatroomParticipantsAPIVersionCode
			headers[utils.HeadersApiVersion] = common.ChatroomParticipantsAPIVersion
			// Make the API call to fetch participants
			chatroomParticipantsAPIResponseBytes, statusCode, err := utils.GetRequestResponseWithoutContext(
				utils.CoreService,
				chatroomParticipantsAPIEndpoint,
				utils.GETRequest,
				headers,
				chatroomParticipantsAPIParams,
				nil,
			)
			// Check if the chatroomParticipantAPIResponse is empty or if there was an error
			if err != nil || chatroomParticipantsAPIResponseBytes == nil || statusCode != http.StatusOK {
				break
			}

			// Parse the chatroomParticipantAPIResponse to extract participant IDs
			var chatroomParticipantAPIResponse struct {
				TotalParticipantsCount int `json:"total_participants_count"`
				Participants           []struct {
					ID string `json:"uuid"`
				} `json:"participants"`
			}
			if err := json.Unmarshal(chatroomParticipantsAPIResponseBytes, &chatroomParticipantAPIResponse); err != nil {
				break
			}
			// If no participants were found, exit the loop
			if len(chatroomParticipantAPIResponse.Participants) == 0 {
				break
			}

			// Collect participant IDs
			for _, participant := range chatroomParticipantAPIResponse.Participants {
				allParticipantIDs = append(allParticipantIDs, participant.ID)
			}
			if len(allParticipantIDs) == chatroomParticipantAPIResponse.TotalParticipantsCount {
				break
			}

			// Increment page for the next request
			currentPage, _ := strconv.Atoi(chatroomParticipantsAPIParams[common.ParamPage])
			chatroomParticipantsAPIParams[common.ParamPage] = strconv.Itoa(currentPage + 1)
		}

		// Update Redis cache with the collected participant IDs if any participants were found
		if len(allParticipantIDs) > 0 {
			err = saveParticipantsInCache(redisClient, chatroomParticipantsKey, allParticipantIDs)
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
func getParticipantsFromCache(redisClient *redis.Client, chatroomParticipantsKey string) ([]string, error) {
	// Get data from Redis
	chatroomParticipantsValue, _, err := cache.Get(redisClient, chatroomParticipantsKey)
	if err != nil {
		return nil, err
	}
	if chatroomParticipantsValue == "" {
		return []string{}, nil
	}

	// Parse the cached data into a slice of strings (assuming JSON array format)
	var allParticipantIDs []string
	err = json.Unmarshal([]byte(chatroomParticipantsValue), &allParticipantIDs)
	if err != nil {
		return nil, err
	}

	return allParticipantIDs, nil
}

// setParticipantsInCache to marshal data and set in cache
func saveParticipantsInCache(redisClient *redis.Client, chatroomParticipantsKey string, allParticipantIDs []string) error {
	// Serialize the value to JSON
	chatroomParticipantsBytes, err := json.Marshal(allParticipantIDs)
	if err != nil {
		return err
	}

	// Store the chatroomParticipantsBytes in Redis with an expiration time (e.g., 24 hour)
	err = cache.Set(redisClient, chatroomParticipantsKey, chatroomParticipantsBytes, time.Hour*cache.ChatroomParticipantsTTL)
	if err != nil {
		return err
	}
	return nil
}

func updateReadDeadline(conn *websocket.Conn) {
	//client.conn.SetReadLimit(maxMessageSize)
	// SetReadDeadline to time.Now() + PongWait
	err := conn.SetReadDeadline(time.Now().Add(common.PongWait))
	if err != nil {
		logging.Info(common.ErrorReadDeadlineWs, err)
		return
	}
}

func updateWriteDeadline(conn *websocket.Conn) {
	// SetWriteDeadline to time.Now() + WriteWait
	err := conn.SetWriteDeadline(time.Now().Add(common.WriteWait))
	if err != nil {
		logging.Info(common.ErrorWriteDeadlineWs, err)
		return
	}
}
