package chatroom

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

type EditCohortAccessRequest struct {
	ChatroomID   interface{} `json:"chatroom_id" binding:"required"`
	CohortID     int         `json:"cohort_id"`
	CohortAccess int         `json:"cohort_access"`
}

// Get cohort access in chatroom
func GetCohortAccess(c *gin.Context) {
	CohortAccess(c, utils.GETMethod)
}

// Edit cohort access in chatroom
func EditCohortAccess(c *gin.Context) {
	CohortAccess(c, utils.PUTMethod)
}

func CohortAccess(c *gin.Context, method int) {

	// Authorize User
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

		// Params to be sent in fetch cohort access api internally
		params := map[string]string{
			ParamChatroomId: c.Query(ParamChatroomId),
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, GetCohortAccessEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	case utils.PUTMethod:

		// Body to be sent in the edit cohort access api internally
		editCohortAccessRequest, err := parseEditCohortAccessRequest(c)
		if err != nil {
			// If POST body params are missing
			utils.GeneralBadRequestError(c, err.Error())
			return
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, EditCohortAccessEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, editCohortAccessRequest)
	}
}

func parseEditCohortAccessRequest(c *gin.Context) (*EditCohortAccessRequest, error) {

	// POST body params
	var ecar EditCohortAccessRequest

	if err := c.ShouldBindJSON(&ecar); err != nil {
		return nil, err
	}

	if ecar.ChatroomID != nil {
		ecar.ChatroomID = utils.ParseInterfaceToString(ecar.ChatroomID)
	}

	return &ecar, nil
}
