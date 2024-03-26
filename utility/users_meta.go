package utility

import (
	"encoding/json"
	"fmt"

	"github.com/nateshr/likeminds-authentication/utils"
)

type UserInfo struct {
	UserID             int64  `json:"user_id"`
	UserUniqueID       string `json:"user_unique_id"`
	ClientUserUniqueID string `json:"clients_user_unique_id"`
}

type UsersInfo struct {
	Success      bool       `json:"success"`
	ErrorMessage string     `json:"error_message"`
	Users        []UserInfo `json:"users"`
}

func getUsersInfoInternal(headers map[string]interface{}, params map[string]string) ([]UserInfo, error) {

	// Internally call /api/community/users
	respBytes, _, err := utils.GetRequestResponseWithoutContext(utils.CoreService, UserMetaInfoInternalEndpoint, utils.GETRequest, headers, params, nil)
	if err != nil {
		return nil, err
	}

	// unmarshal response
	var usersInfo UsersInfo
	if err := json.Unmarshal(respBytes, &usersInfo); err != nil {
		return nil, err
	}

	// check if response is successful
	if !usersInfo.Success {
		return nil, fmt.Errorf(usersInfo.ErrorMessage)
	}

	// return users info
	return usersInfo.Users, nil
}

// External method to fetch user_unique_ids from user_ids, uuids or client_uuids
func FetchUserUniqueIdsFromAnyUserIds(headers map[string]interface{}, userIds interface{}) ([]string, error) {

	params := map[string]string{}

	// type assert userIds and set params
	switch userIds := userIds.(type) {
	case []string:
		if len(userIds) == 0 {
			return nil, fmt.Errorf("userIds cannot be empty")
		}
		params[ParamMemberIDs] = utils.ParseStringArrayToString(userIds)
	case []interface{}:
		if len(userIds) == 0 {
			return nil, fmt.Errorf("userIds cannot be empty")
		}
		params[ParamMemberIDs] = utils.ParseInterfaceListToStringList(userIds)
	case string:
		if userIds == "" {
			return nil, fmt.Errorf("userIds cannot be empty")
		}
		params[ParamMemberIDs] = userIds
	}

	usersInfo, err := getUsersInfoInternal(headers, params)
	if err != nil {
		return nil, err
	}

	user_unique_ids := []string{}
	for _, userInfo := range usersInfo {
		user_unique_ids = append(user_unique_ids, userInfo.UserUniqueID)
	}

	return user_unique_ids, nil
}

// GetUUIDInternally is used to get user_unique_id from user_id
func GetUUIDInternally(headers map[string]interface{}, user_id string) (string, error) {

	member_ids := []string{user_id}

	//Get user_unique_id by calling internal core service
	user_unique_ids, err := FetchUserUniqueIdsFromAnyUserIds(headers, member_ids)
	if err != nil {
		return "", err
	}

	if len(user_unique_ids) == 0 {
		return "", fmt.Errorf(utils.ErrorInvalidUserId)
	}

	return user_unique_ids[0], nil
}
