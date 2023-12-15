package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateTopicsRequest struct {
	Name  string   `json:"name,omitempty"`
	Names []string `json:"names,omitempty"`
}

type DeleteTopicsRequest struct {
	TopicIds []string `json:"topic_ids"`
}

type EditTopicRequest struct {
	Name      string `json:"name"`
	IsEnabled bool   `json:"is_enabled"`
	UserIsCm  bool   `json:"user_is_cm"`
}

func parseCreateTopicsRequest(c *gin.Context) (*CreateTopicsRequest, error) {
	//POST body params
	var ctr CreateTopicsRequest

	if err := c.ShouldBindJSON(&ctr); err != nil {
		return nil, err
	}

	return &ctr, nil
}

func parseDeleteTopicsRequest(c *gin.Context) (*DeleteTopicsRequest, error) {
	//DELETE body params
	var dtr DeleteTopicsRequest

	if err := c.ShouldBindJSON(&dtr); err != nil {
		return nil, err
	}

	return &dtr, nil
}

func parseEditTopicRequest(c *gin.Context) (*EditTopicRequest, error) {
	//PUT body params
	var etr EditTopicRequest

	if err := c.ShouldBindJSON(&etr); err != nil {
		return nil, err
	}

	return &etr, nil
}

// CreateTopics is used to create a new topic
func CreateTopics(c *gin.Context) {
	Topic(c, utils.POSTMethod)
}

// GetTopic is used to get topics of a specific community
func GetTopic(c *gin.Context) {
	Topic(c, utils.GETMethod)
}

// DeleteTopics is used to delete topics of a specific community
func DeleteTopics(c *gin.Context) {
	Topic(c, utils.DELETEMethod)
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
		createTopicsInternal(c, userId)

	case utils.DELETEMethod:
		deleteTopicsInternal(c, userId)

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
		ParamMinPosts:   c.Query(ParamMinPosts),
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

func createTopicsInternal(c *gin.Context, userId string) {
	//Body to be sent in the /topic POST request
	createTopicsRequest, err := parseCreateTopicsRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Fetch member access to create topics
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
	utils.SendRequest(c, utils.SwarmService, TopicEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createTopicsRequest)
}

func deleteTopicsInternal(c *gin.Context, userId string) {
	//Body to be sent in the /topic DELETE request
	deleteTopicsRequest, err := parseDeleteTopicsRequest(c)
	if err != nil {
		//If DELETE body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Fetch member access to delete topics
	success, response := user.FetchMemberAccess(c, DELETE_TOPIC_ACTION, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, TopicEndPoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, deleteTopicsRequest)
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
