package pubsub

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nateshr/likeminds-authentication/internal/handlers/channel"
	"github.com/nateshr/likeminds-authentication/internal/handlers/chatroom"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
	"github.com/nateshr/likeminds-authentication/internal/utils/api_client"
	"log"
	"net/http"
	"strings"
	"time"
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
	topic := c.Param(ParamTopic)
	topicSplit, err := GetTopicSplit(topic)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	switch topicSplit[0] {
	case TopicTypeChatroom:
		UUID := user.GetRequestingUserId(c)
		var chatroomID string
		if len(topicSplit) > 1 {
			chatroomID = topicSplit[1]
		}
		//If UUID is missing, return error
		if UUID == "" || UUID == "null" {
			utils.GeneralBadRequestError(c, ErrorUserUUIDMissing)
			return
		}
		//If chatroomID is missing, return error
		if chatroomID == "" || chatroomID == "null" {
			utils.GeneralBadRequestError(c, ErrorChatroomIDMissing)
			return
		}
		params := map[string]string{
			channel.ParamChannelId:          chatroomID,
			channel.ParamChannelActionTypes: c.Query(channel.ParamChannelActionTypes),
		}
		//Get chatroom details to verify if user has access to any chatroom / cohort based chatroom / secret chatroom
		respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.CoreService, channel.SyncChannelDetailEndppoint, utils.GETRequest, utils.CreateHeaders(c, UUID), params, nil)
		if err != nil || statusCode != 200 {
			utils.GeneralBadRequestError(c, fmt.Sprintf(ErrorUserChatroomAccess, err))
			return
		} else {
			var chatroomDetailParentResponse chatroom.ChatroomDetailParentResponse
			if err := json.Unmarshal(respBytes, &chatroomDetailParentResponse); err != nil {
				utils.GeneralAPIError(c, fmt.Sprintf(ErrorUnmarshalErrorJson, err))
				return
			}
			chatroomDetailArray := chatroomDetailParentResponse.ChatroomDetail
			if len(chatroomDetailArray) < 1 {
				utils.GeneralAPIError(c, ErrorChatroomResponseInvalid)
				return
			}
			chatroomDetail := chatroomDetailArray[0]
			canAccessSecretChatroom := chatroomDetail.CanAccessSecretChatroom
			if canAccessSecretChatroom != nil {
				if *canAccessSecretChatroom == false {
					utils.GeneralBadRequestError(c, ErrorUserChatroomAccess)
					return
				}
			}
			cohortAccess := chatroomDetail.CohortAccess
			if cohortAccess != nil {
				if *cohortAccess != 200 {
					utils.GeneralBadRequestError(c, ErrorUserChatroomAccess)
					return
				}
			}
		}
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
	go readFromClientAndWriteToServer(conn, serverConn)
	go readFromServerAndWriteToClient(conn, serverConn)
}

// GetTopicSplit will decode topic and return split
func GetTopicSplit(topic string) ([]string, error) {
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

	go startPingMessageToServer(serverConn)
	serverConn.SetPongHandler(func(string) error {
		log.Println(PongReceivedWs)
		return nil
	})

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			logging.Error(fmt.Sprintf(ErrorReadClientWs, err))
			return
		}
		logging.Info(fmt.Sprintf(ReceivedMessageClientWs, messageType))

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

	go startPingMessageToClient(conn)
	conn.SetPongHandler(func(string) error {
		log.Println(PongReceivedClient)
		return nil
	})

	for {
		messageType, msg, err := serverConn.ReadMessage()
		if err != nil {
			logging.Error(fmt.Sprintf(ErrorReadServerWs, err))
			return
		}
		logging.Info(fmt.Sprintf(ReceivedMessageServerWs, messageType))

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
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf(fmt.Sprintf(ErrorPingSendWs, err))
			return
		}
		log.Println(PingSendWs)
	}
}
