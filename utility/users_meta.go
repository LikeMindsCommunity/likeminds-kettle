package utility

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func GetUsersInfo(c *gin.Context, member_ids []interface{}, only_user_unique_ids bool) (interface{}, error) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return "", errors.New("Some in getting bot ID!")
	}

	// Create request param to member_profile
	params := map[string]string{
		ParamMemberIDs: utils.ParseArrayToString(member_ids),
	}

	// Internally call api/community_member/can_dm
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, UserMetaInfoEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
	if respBytes == nil {
		return "", errors.New("Some error occured!")
	}

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return "", errors.New("Some error occured!")
	}

	if apiCR.Success {
		dataResponse := apiCR.Response

		user_data, ok := dataResponse["users"]

		if ok && only_user_unique_ids {
			var user_unique_ids []interface{}

			for _, v := range user_data.([]interface{}) {
				user_unique_id, ok := v.(map[string]interface{})["user_unique_id"]

				if ok {
					user_unique_ids = append(user_unique_ids, user_unique_id.(string))
				}
			}

			return user_unique_ids, nil

		}
	}

	return &apiCR, nil

}
