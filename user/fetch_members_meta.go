package user

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/utils"
)

type MemberMeta struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	ImageUrl     string `json:"image_url"`
	UserUniqueId string `json:"user_unique_id"`
	IsGuest      bool   `json:"is_guest"`
	IsDeleted    bool   `json:"is_deleted"`
	CustomTitle  string `json:"custom_title"`
}

type MemberMetaResponse struct {
	Success bool         `json:"success"`
	Members []MemberMeta `json:"members"`
}

// FetchMemberMeta | fetch member meta for sent ids
func FetchMemberMeta(c *gin.Context, member_ids []string, userId string) (bool, map[string]MemberMeta) {

	response := map[string]MemberMeta{}

	if len(member_ids) == 0 {
		return true, response
	}

	temp_params, _ := json.Marshal(member_ids)

	//Params to be sent in the api/community_member/fetch_access request
	params := map[string]string{
		ParamMemberIds: fmt.Sprintf("%v", string(temp_params)),
	}

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchMembersMetaEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return false, nil
	}

	//Parse response
	var membersMetaResponse MemberMetaResponse
	if err := json.Unmarshal(respBytes, &membersMetaResponse); err != nil {
		//Internal unmarshal error
		utils.GeneralAPIError(c, err.Error())
		return false, nil
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

	return true, response
}
