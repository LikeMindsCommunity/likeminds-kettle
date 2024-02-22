package feed

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func GetUserFeedMeta(c *gin.Context) {
	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Access query params and url generation
	userID := c.Param("user_id")

	if userID == "" {
		utils.GeneralBadRequestError(c, utils.ErrorInvalidUserId)
	}

	endpoint := fmt.Sprintf(FetchUserFeedMetaEndPoint, userID)

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, endpoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)

	// Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	// If flow succeeds
	dataResponse := apiCR.Response

	//Fetch user data for given user_unique_ids
	userIds := []string{userID}

	user_data, err := user.FetchMemberMeta(utils.CreateHeaders(c, userId), userIds)
	if err != nil {
		utils.GeneralBadRequestError(c, utils.ErrorFetchingUserData)
		return
	}

	dataResponse["users"] = user_data

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}
