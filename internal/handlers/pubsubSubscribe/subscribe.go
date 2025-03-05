package pubsubSubscribe

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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
		ReadBufferSize:  ReadBufferSizeDefault,
		WriteBufferSize: WriteBufferSizeDefault,
	}
}

// upgraderHTTPToWs to upgrade the incoming HTTP request to a WebSocket connection
func upgraderHTTPToWs(c *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	return conn, err
}

// Subscribe to open WS against a topic
func Subscribe(c *gin.Context) {
	// Validate params and headers before subscribing to a topic
	topicSplit, userID, topicID, err := validateParamsAndHeaders(c)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	// Check if user has access to chatroom
	statusCode, err := hasAccessToChatroom(c, topicSplit, userID, topicID)
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
		updatedErr := fmt.Sprintf(ErrorFailedUpgrader, err)
		logging.Error(updatedErr)
		utils.GeneralAPIError(c, updatedErr)
		return
	}

	// Connect to the websocket server
	serverConn, err := dialToWs(c)
	if err != nil {
		updatedErr := fmt.Sprintf(ErrorFailedDial, err)
		logging.Error(updatedErr)
		utils.GeneralAPIError(c, updatedErr)
		return
	}

	logging.Info(WsConnectionEstablished)

	// Handle communication between the client and the websocket server
	utils.SafeGo(func() { readFromClientAndWriteToServer(conn, serverConn) })
	utils.SafeGo(func() { readFromServerAndWriteToClient(conn, serverConn) })
}

// validateParamsAndHeaders to validate params and headers sent while subscribing to topic
func validateParamsAndHeaders(c *gin.Context) ([]string, string, string, error) {
	//Validate topic
	topic := c.Param(ParamTopic)
	topicSplit, err := getTopicSplit(topic)
	if err != nil {
		return nil, "", "", err
	}
	userID := user.GetRequestingUserId(c)
	var topicID string
	if len(topicSplit) > 1 {
		topicID = topicSplit[1]
	}
	//If userID is missing, return error
	if userID == "" || userID == "null" {
		return nil, "", "", errors.New(ErrorUserUUIDMissing)
	}
	//If topicID is missing, return error
	if topicID == "" || topicID == "null" {
		return nil, "", "", errors.New(ErrorTopicIDMissing)
	}
	return topicSplit, userID, topicID, nil
}

// hasAccessToChatroom to check if userID has access to topicID while subscribing to chatroom topic
func hasAccessToChatroom(c *gin.Context, topicSplit []string, userID string, topicID string) (int, error) {
	switch topicSplit[0] {
	case TopicTypeChatroom:
		params := map[string]string{
			channel.ParamChannelId:          topicID,
			channel.ParamChannelActionTypes: c.Query(channel.ParamChannelActionTypes),
		}
		//Get chatroom details to verify if user has access to any chatroom / cohort based chatroom / secret chatroom
		respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.CoreService, channel.SyncChannelDetailEndppoint, utils.GETRequest, utils.CreateHeaders(c, userID), params, nil)
		if err != nil || statusCode != 200 {
			return statusCode, errors.New(fmt.Sprintf(ErrorUserChatroomAccess, err))
		} else {
			var chatroomDetailParentResponse chatroom.ChatroomDetailParentResponse
			if err := json.Unmarshal(respBytes, &chatroomDetailParentResponse); err != nil {
				return statusCode, errors.New(fmt.Sprintf(ErrorUnmarshalErrorJson, err))
			}
			chatroomDetailArray := chatroomDetailParentResponse.ChatroomDetail
			if len(chatroomDetailArray) < 1 {
				return statusCode, errors.New(ErrorChatroomResponseInvalid)
			}
			//If user has access to secret chatroom
			chatroomDetail := chatroomDetailArray[0]
			canAccessSecretChatroom := chatroomDetail.CanAccessSecretChatroom
			if canAccessSecretChatroom != nil {
				if *canAccessSecretChatroom == false {
					return statusCode, errors.New(ErrorUserChatroomAccess)
				}
			}
			//If user has access to cohort based chatroom
			cohortAccess := chatroomDetail.CohortAccess
			if cohortAccess != nil {
				if *cohortAccess != 200 {
					return statusCode, errors.New(ErrorUserChatroomAccess)
				}
			}
		}
	}
	return http.StatusOK, nil
}

// getTopicSplit will decode topic and return split
func getTopicSplit(topic string) ([]string, error) {
	if topic == "" || topic == "null" {
		return nil, errors.New(ErrorTopicMissing)
	}
	topicSplit := strings.Split(topic, ":")
	if len(topicSplit) <= 1 {
		return nil, errors.New(ErrorTopicInvalid)
	}
	return topicSplit, nil
}

// createHeaders to createHeaders required for connecting with websocket server
func createHeaders(c *gin.Context) http.Header {
	out := http.Header{}
	userId := user.GetRequestingUserId(c)
	deviceID := user.GetRequestingUserDeviceId(c)
	headersMap := utils.CreateHeadersFromToken(c, userId, deviceID)
	for key, value := range headersMap {
		out.Add(key, value.(string))
	}
	return out
}

// dialToWs to dial to websocket server
func dialToWs(c *gin.Context) (*websocket.Conn, error) {
	topic := c.Param("topic")
	psURL := api_client.GetPandemoniumServiceWsUrl()
	updatedPsURL := fmt.Sprintf("%s/subscribe/%s", psURL, topic)
	serverConn, _, err := websocket.DefaultDialer.Dial(updatedPsURL, createHeaders(c))
	return serverConn, err
}

func readFromClientAndWriteToServer(conn *websocket.Conn, serverConn *websocket.Conn) {
	defer func() {
		disconnect(conn)
		disconnect(serverConn)
	}()

	//conn is set to read till PongWait (60 seconds)
	updateReadDeadline(conn)

	utils.SafeGo(func() {
		//write PING message to conn after every PingPeriod (54 seconds)
		startPingMessageToClient(conn)
	})

	//conn replied with PONG
	conn.SetPongHandler(func(string) error {
		log.Println(PongReceivedClient)
		//serverConn is set to read till PongWait (60 seconds)
		updateReadDeadline(conn)
		return nil
	})

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			logging.Error(fmt.Sprintf(ErrorReadClientWs, err))
			return
		}
		logging.Info(fmt.Sprintf(ReceivedMessageClientWs, messageType))

		//serverConn is set to write till WriteWait (10 seconds)
		updateWriteDeadline(serverConn)
		// Forward the message to the WebSocket server
		err = serverConn.WriteMessage(messageType, msg)
		if err != nil {
			logging.Error(fmt.Sprintf(ErrorWriteServerWs, err))
			return
		}
	}
}

func readFromServerAndWriteToClient(conn *websocket.Conn, serverConn *websocket.Conn) {
	defer func() {
		disconnect(conn)
		disconnect(serverConn)
	}()

	//serverConn is set to read till PongWait (60 seconds)
	updateReadDeadline(serverConn)

	utils.SafeGo(func() {
		//write PING message to serverConn after PingPeriod (54 seconds)
		startPingMessageToServer(serverConn)
	})

	//serverConn replied with PONG
	serverConn.SetPongHandler(func(string) error {
		log.Println(PongReceivedWs)
		//serverConn is set to read till PongWait (60 seconds)
		updateReadDeadline(serverConn)
		return nil
	})

	for {
		messageType, msg, err := serverConn.ReadMessage()
		if err != nil {
			logging.Error(fmt.Sprintf(ErrorReadServerWs, err))
			return
		}
		logging.Info(fmt.Sprintf(ReceivedMessageServerWs, messageType))

		//conn is set to write till WriteWait (10 seconds)
		updateWriteDeadline(conn)
		// Forward the message to the client
		err = conn.WriteMessage(messageType, msg)
		if err != nil {
			logging.Error(fmt.Sprintf(ErrorWriteClientWs, err))
			return
		}
	}
}

func disconnect(conn *websocket.Conn) {
	logging.Info(ConnectionClosed)
	err := conn.Close()
	if err != nil {
		log.Println(ErrorUnableToCloseWs, err)
		return
	}
}

func startPingMessageToClient(conn *websocket.Conn) {
	// Start a goroutine to send pings periodically to the client
	for {
		time.Sleep(PingPeriod) // Interval between pings
		//conn is set to write till WriteWait (10 seconds)
		updateWriteDeadline(conn)
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf(fmt.Sprintf(ErrorPingSendClient, err))
			return
		}
		log.Println(PingSendClient)
	}
}

func startPingMessageToServer(conn *websocket.Conn) {
	// Start a goroutine to send pings periodically to the client
	for {
		time.Sleep(PingPeriod) // Interval between pings
		//conn is set to write till WriteWait (10 seconds)
		updateWriteDeadline(conn)
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf(fmt.Sprintf(ErrorPingSendWs, err))
			return
		}
		log.Println(PingSendWs)
	}
}

func updateReadDeadline(conn *websocket.Conn) {
	// SetReadDeadline to time.Now() + PongWait
	err := conn.SetReadDeadline(time.Now().Add(PongWait))
	if err != nil {
		logging.Info(ErrorReadDeadlineWs, err)
		return
	}
}

func updateWriteDeadline(conn *websocket.Conn) {
	// SetWriteDeadline to time.Now() + WriteWait
	err := conn.SetWriteDeadline(time.Now().Add(WriteWait))
	if err != nil {
		logging.Info(ErrorWriteDeadlineWs, err)
		return
	}
}
