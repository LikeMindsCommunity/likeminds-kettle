package moderation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type Right struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	SubTitle   string `json:"sub_title"`
	State      int32  `json:"state"`
	IsSelected bool   `json:"is_selected"`
	IsLocked   bool   `json:"is_locked"`
}

type RightsRequest struct {
	UserId      int64   `json:"user_id" binding:"required"`
	CustomTitle string  `json:"custom_title"`
	Rights      []Right `json:"rights"`
	IsCM        bool    `json:"is_cm"`
}

//EditRights is used to edit community rights for members
func EditRights(c *gin.Context) {
	Rights(c, utils.PUTMethod)
}

//GetRights is used to get community rights for members
func GetRights(c *gin.Context) {
	Rights(c, utils.GETMethod)
}

//Rigths method handles community rights for members
func Rights(c *gin.Context, method int) {

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

		//Params to be sent in the fetch rights request
		params := map[string]string{
			ParamUserId: c.Query(ParamUserId),
		}

		//GET Request params
		is_cm := c.Query(ParamIsCm)

		if is_cm == "" || is_cm == "false" {
			//If is_cm is missing or false, call fetch member rights api internally

			//Send Request
			utils.SendRequest(c, utils.CoreService, FetchMemberRights, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		} else {
			//else, call fetch cm rights api internally

			//Send Request
			utils.SendRequest(c, utils.CoreService, FetchCMRights, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		}

	case utils.PUTMethod:

		//Body to be sent in the update rights POST request
		rightsRequest, err := parseRightsRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		is_cm := rightsRequest.IsCM

		if !is_cm {
			//If is_cm is missing or false, call update member rights api internally

			//Send Request
			utils.SendRequest(c, utils.CoreService, UpdateMemberRights, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, rightsRequest)
		} else {
			//else, call update cm rights api internally

			//Send Request
			utils.SendRequest(c, utils.CoreService, UpdateCMRights, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, rightsRequest)
		}
	}
}

func parseRightsRequest(c *gin.Context) (*RightsRequest, error) {
	//POST body params
	var rr RightsRequest

	if err := c.ShouldBindJSON(&rr); err != nil {
		return nil, err
	}

	return &rr, nil
}
