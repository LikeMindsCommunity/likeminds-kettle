package community

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type RemoveCommunityManagerRequest struct {
	UserId interface{} `json:"user_id"`
	UUID   string      `json:"uuid"`
}

// RemoveCommunityManager is used to remove a CM from community
func RemoveCommunityManager(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Body to be sent in the remove member request
	removeCMRequest, err := parseRemoveCMRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	if removeCMRequest.UserId != nil && reflect.TypeOf(removeCMRequest.UserId).String() == "float64" {
		removeCMRequest.UserId = strconv.Itoa(int(removeCMRequest.UserId.(float64)))
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, RemoveCMEndPoint, utils.POSTRequestFormUrlEncodedBody, utils.CreateHeaders(c, userId), nil, removeCMRequest)

	// delete cached user feed access rights
	user.DeleteAccessDataAgainstUserIdAndAccessTypeFromCache(utils.GetRedisClientFromContext(c), removeCMRequest.UserId.(string))

}

func parseRemoveCMRequest(c *gin.Context) (*RemoveCommunityManagerRequest, error) {
	//POST body params
	var rcmr RemoveCommunityManagerRequest

	if err := c.ShouldBindJSON(&rcmr); err != nil {
		return nil, err
	}

	return &rcmr, nil
}
