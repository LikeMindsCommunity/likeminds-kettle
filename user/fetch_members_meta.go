package user

import (
	"encoding/json"
	"fmt"

	"github.com/nateshr/likeminds-authentication/utils"
)

type SdkClientInfo struct {
	CommunityId  int    `json:"community"`
	UserId       int    `json:"user"`
	UserUniqueId string `json:"user_unique_id"`
	UUID         string `json:"uuid"`
}

type MemberMeta struct {
	Id            int            `json:"id"`
	Name          string         `json:"name"`
	ImageUrl      string         `json:"image_url"`
	UserUniqueId  string         `json:"user_unique_id"`
	SdkClientInfo *SdkClientInfo `json:"sdk_client_info"`
	UUID          string         `json:"uuid"`
	IsGuest       bool           `json:"is_guest"`
	IsDeleted     bool           `json:"is_deleted"`
	CustomTitle   string         `json:"custom_title"`
}

type MemberMetaResponse struct {
	Success bool         `json:"success"`
	Members []MemberMeta `json:"members"`
}

// FetchMemberMeta without context | fetch member meta for sent ids
func FetchMemberMeta(headers map[string]interface{}, member_ids []string) (map[string]MemberMeta, error) {

	response := map[string]MemberMeta{}

	if len(member_ids) == 0 {
		return response, nil
	}

	temp_params, _ := json.Marshal(member_ids)

	//Params to be sent in the api/community_member/fetch_access request
	params := map[string]string{
		ParamMemberIds: fmt.Sprintf("%v", string(temp_params)),
	}

	//Send Request
	respBytes, _, err := utils.GetRequestResponseWithoutContext(utils.CoreService, FetchMembersMetaEndPoint, utils.GETRequest, headers, params, nil)

	if err != nil {
		return nil, err
	}

	//Parse response
	var membersMetaResponse MemberMetaResponse
	if err := json.Unmarshal(respBytes, &membersMetaResponse); err != nil {
		return nil, err
	}

	//Generate user data for received data
	for _, memberData := range membersMetaResponse.Members {
		response[memberData.UserUniqueId] = memberData
	}

	//Generate user data for remaining member_ids
	for _, memberId := range member_ids {
		if _, ok := response[memberId]; !ok {
			response[memberId] = MemberMeta{
				IsDeleted: true,
			}
		}
	}

	return response, err
}

// This function is used to fetch members meta from user_ids of feed entity data
func GetUsersMetaFromFeedData(headers map[string]interface{}, feedDataArray []interface{}) (map[string]MemberMeta, error) {

	user_unique_ids := []string{}

	// Fetch user ids from array
	for _, data := range feedDataArray {
		if user_unique_id, ok := data.(map[string]interface{})["user_id"]; ok {
			user_unique_ids = append(user_unique_ids, user_unique_id.(string))
		}
	}

	// Fetch user data for given user_unique_ids
	user_data, err := FetchMemberMeta(headers, user_unique_ids)
	if err != nil {
		return nil, err
	}

	return user_data, nil
}
