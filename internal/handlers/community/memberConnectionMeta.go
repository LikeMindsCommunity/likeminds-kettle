package community

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// GetMemberConnection is used to get a member connections
func GetMemberConnectionMeta(c *gin.Context) {
	MemberConnectionMeta(c, utils.GETMethod)
}

// MemberConnectionMeta handles all connection meta related objects
func MemberConnectionMeta(c *gin.Context, method int) {
	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}
	// user id received in path params
	paramUUID := c.Param(ParamUserId)
	memberConnectionMetaEndPoint := fmt.Sprintf(MemberConnectionMetaEndPoint, paramUUID)

	// Send request
	switch method {
	case utils.GETMethod:
		getMemberConnectionMetaInternal(c, userId, memberConnectionMetaEndPoint)
	}
}

func getMemberConnectionMetaInternal(c *gin.Context, userId string, getMemberConnectionMetaEndPoint string) {
	// Send Request
	utils.SendRequest(c, utils.CoreService, getMemberConnectionMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
