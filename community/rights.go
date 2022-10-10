package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CommunityRight struct {
	Id         int    `json:"id" binding:"required"`
	State      int    `json:"state" binding:"required"`
	Title      string `json:"title" binding:"required"`
	IsSelected bool   `json:"is_selected" binding:"required"`
	IsLocked   bool   `json:"is_locked"`
}

type EditCommunityRightsRequest struct {
	Rights []CommunityRight `json:"rights" binding:"required"`
}

// EditCommunityRights is used to update community rights
func EditCommunityRights(c *gin.Context) {
	CommunityRights(c, utils.PUTMethod)
}

func parseEditCommunityRightsRequest(c *gin.Context) (*EditCommunityRightsRequest, error) {
	//POST body params
	var ecrr EditCommunityRightsRequest

	if err := c.ShouldBindJSON(&ecrr); err != nil {
		return nil, err
	}

	return &ecrr, nil
}

func CommunityRights(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	switch method {
	case utils.PUTMethod:

		// Body to be sent in the edit rights api internally
		editCommunityRightsRequest, err := parseEditCommunityRightsRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, EditCommunityRightsEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, editCommunityRightsRequest)
	}

}
