package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CommunityRight struct {
	Id         int    `json:"id" binding:"required"`
	State      int    `json:"state"`
	Title      string `json:"title" binding:"required"`
	IsSelected bool   `json:"is_selected" binding:"required"`
	IsLocked   bool   `json:"is_locked"`
}

type EditCommunityRightsRequest struct {
	Rights []CommunityRight `json:"rights" binding:"required"`
}

func parseEditCommunityRightsRequest(c *gin.Context) (*EditCommunityRightsRequest, error) {
	//POST body params
	var ecrr EditCommunityRightsRequest

	if err := c.ShouldBindJSON(&ecrr); err != nil {
		return nil, err
	}

	return &ecrr, nil
}

// GetCommunityRights is used to get community rights
func GetCommunityRights(c *gin.Context) {
	CommunityRights(c, utils.GETMethod)
}

// EditCommunityRights is used to update community rights
func EditCommunityRights(c *gin.Context) {
	CommunityRights(c, utils.PUTMethod)
}

// UpdateCommunityRights is used to update community rights only sent in the request
func UpdateCommunityRights(c *gin.Context) {
	CommunityRights(c, utils.PatchMethod)
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
	case utils.GETMethod:

		//Send Request
		utils.SendRequest(c, utils.CoreService, FetchCommunityRightsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

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

	case utils.PatchMethod:
		UpdateCommunityRightsInternal(c, userId)

	}

}

func UpdateCommunityRightsInternal(c *gin.Context, userId string) {

	// Body to be sent in the edit rights api internally
	editCommunityRightsRequest, err := parseEditCommunityRightsRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request to api/update_community_rights with PATCH method
	utils.SendRequest(c, utils.CoreService, EditCommunityRightsEndPoint, utils.PATCHRequest, utils.CreateHeaders(c, userId), nil, editCommunityRightsRequest)

}
