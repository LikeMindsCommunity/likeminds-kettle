package utility

import (
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
			return response, nil
		}

		if only_user_unique_ids {

			for _, v := range user_data.([]interface{}) {
				user_unique_id, ok := v.(map[string]interface{})["user_unique_id"]

				if ok {
					response = append(response, user_unique_id.(string))
				}
			}

		} else {
			response = user_data.([]interface{})
		}
	}

	return response, nil

}
