package community

import (
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetCommunityQuestionFilters(c *gin.Context) {
	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, FetchCommunityQuestionFiltersEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
}
