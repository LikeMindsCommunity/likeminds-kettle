package chatroom

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// CollabcardSeen is used to mark a chatroom as seen
func CollabcardSeen(c *gin.Context) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Params to be sent in the collabcard seen api internally
	params := map[string]string{
		ParamCollabcardId:   c.Query(ParamCollabcardId),
		ParamMemberId:       c.Query(ParamMemberId),
		ParamUUID:           c.Query(ParamUUID),
		ParamCollabcardType: c.Query(ParamCollabcardType),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, CollabcardSeenEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), params, nil)
}
