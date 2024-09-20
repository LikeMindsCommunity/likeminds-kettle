package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type BlockUnblockUserRequest struct {
	ShouldBlock bool `json:"should_block"`
}

// BlockUnblockUser is used to block unblock in community
func BlockUnblockUser(c *gin.Context) {
	BlockUser(c, utils.PUTMethod)
}

// GetBlockUser is used to get block user data in community
func GetBlockUser(c *gin.Context) {
	BlockUser(c, utils.GETMethod)
}

// BlockUser method handles block user API's
func BlockUser(c *gin.Context, method int) {
	//Authorize User
	userId := GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Send request
	switch method {

	case utils.PUTMethod:
		blockUnblockUserInternal(c, userId)

	case utils.GETMethod:
		getBlockUserDataInternal(c, userId)
	}
}

func blockUnblockUserInternal(c *gin.Context, userId string) {

	//Body to be sent in the edit member api internally
	blockUnblockUserRequest, err := parseBlockUnblockRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	secondUserId := c.Param(ParamUserUUID)
	apiEndpoint := fmt.Sprintf(BlockUserEndpoint, secondUserId)

	//Send Request
	utils.SendRequest(c, utils.CoreService, apiEndpoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, blockUnblockUserRequest)
}

func getBlockUserDataInternal(c *gin.Context, userId string) {
	secondUserId := c.Param(ParamUserUUID)
	apiEndpoint := fmt.Sprintf(BlockUserEndpoint, secondUserId)

	// Params
	params := map[string]string{
		ParamBlockUserType: fmt.Sprintf("[%s]", BlockingUserType),
		ParamPage:          c.Query(ParamPage),
		ParamPageSize:      c.Query(ParamPageSize),
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, apiEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}

func parseBlockUnblockRequest(c *gin.Context) (*BlockUnblockUserRequest, error) {
	//POST body params
	var buur BlockUnblockUserRequest

	if err := c.ShouldBindJSON(&buur); err != nil {
		return nil, err
	}

	return &buur, nil
}
