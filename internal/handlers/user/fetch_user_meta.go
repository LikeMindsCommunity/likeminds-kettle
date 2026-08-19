package user

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

func UserMeta(c *gin.Context) {
	// Authorize User
	userId := GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Send Request
	utils.SendRequest(c, utils.CoreService, FetchUserMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
