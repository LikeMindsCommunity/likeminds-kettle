package utility

import (
	"errors"

	"github.com/nateshr/likeminds-authentication/utils"
)

type UserInfo struct {
	UserID             int64  `json:"user_id"`
	UserUniqueID       string `json:"user_unique_id"`
	ClientUserUniqueID string `json:"clients_user_unique_id"`
}

type UsersInfo struct {
	Users []UserInfo `json:"users"`
}

func GetUsersInfo(headers map[string]interface{}, member_ids []interface{}, only_user_unique_ids bool) (interface{}, error) {

	var response []interface{}

	if len(member_ids) == 0 {
		return response, nil
	}

	// Create request param to member_profile
	params := map[string]string{
		ParamMemberIDs: utils.ParseArrayToString(member_ids),
	}

	// Internally call /api/community/users
	respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.CoreService, UserMetaInfoEndpoint, utils.GETRequest, headers, params, nil)

	if err != nil {
		return nil, err
	}

	dataResponse := utils.ValidateClientResponseWithoutContext(respBytes, statusCode, err)

	// Parse response
	if dataResponse != nil {
		user_data, ok := dataResponse["users"]

		if !ok {
			return nil, errors.New("No users found!")
		}

		if only_user_unique_ids {
			var user_unique_ids []interface{}

			for _, v := range user_data.([]interface{}) {
				user_unique_id, ok := v.(map[string]interface{})["user_unique_id"]

				if ok {
					user_unique_ids = append(user_unique_ids, user_unique_id.(string))
				}
			}

			return user_unique_ids, nil
		} else {
			return user_data, nil
		}
	}

	return dataResponse, nil

}
