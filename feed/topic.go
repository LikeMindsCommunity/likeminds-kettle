package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateTopicRequest struct {
	Name string `json:"name"`
}

type EditTopicRequest struct {
	Name      string `json:"name"`
	IsEnabled bool   `json:"is_enabled"`
	UserIsCm  bool   `json:"user_is_cm"`
}

func parseCreateTopicRequest(c *gin.Context) (*CreateTopicRequest, error) {
	//POST body params
	var ctr CreateTopicRequest

	if err := c.ShouldBindJSON(&ctr); err != nil {
		return nil, err
	}

	return &ctr, nil
}
func parseEditTopicRequest(c *gin.Context) (*EditTopicRequest, error) {
	//POST body params
	var etr EditTopicRequest

	if err := c.ShouldBindJSON(&etr); err != nil {
		return nil, err
	}

	return &etr, nil
}

// CreateTopic is used to create a new topic
func CreateTopic(c *gin.Context) {
	Topic(c, utils.POSTMethod)
}

// GetTopic is used to get topics of a specific community
func GetTopic(c *gin.Context) {
	Topic(c, utils.GETMethod)
}

// EditTopic is used to edit an existing topic
func EditTopic(c *gin.Context) {
	Topic(c, utils.PUTMethod)
}

// Topic method handles topic objects
func Topic(c *gin.Context, method int) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Send request
	switch method {
	case utils.GETMethod:
		GetTopicInternal(c, userId)

	case utils.POSTMethod:
		createTopicInternal(c, userId)

	case utils.PUTMethod:
		editTopicInternal(c, userId)
	}
}

func GetTopicInternal(c *gin.Context, userId string) {
	//Params to be sent in the /topic GET request
	params := map[string]string{
		ParamPage:       c.Query(ParamPage),
		ParamPageSize:   c.Query(ParamPageSize),
		ParamIsEnabled:  c.Query(ParamIsEnabled),
		ParamSearchType: c.Query(ParamSearchType),
		ParamSearch:     c.Query(ParamSearch),
	}

	//Fetch member access to view topics
	success, response := user.FetchMemberAccess(c, IS_MEMBER, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, TopicEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}

func createTopicInternal(c *gin.Context, userId string) {
	//Body to be sent in the /topic POST request
	createTopicRequest, err := parseCreateTopicRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Fetch member access to create topic
	success, response := user.FetchMemberAccess(c, CREATE_TOPIC_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, TopicEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createTopicRequest)
}

func editTopicInternal(c *gin.Context, userId string) {
	topicId := c.Param("topic_id")
	EditTopicEndPoint := fmt.Sprintf(SingleTopicEndPoint, topicId)

	//Body to be sent in the /topic PUT request
	editTopicRequest, err := parseEditTopicRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Fetch member access to create topic
	success, response := user.FetchMemberAccess(c, EDIT_TOPIC_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, EditTopicEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, editTopicRequest)
}
