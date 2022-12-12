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
}

type MemberMetaResponse struct {
	Success bool         `json:"success"`
	Members []MemberMeta `json:"members"`
}

// FetchMemberMeta | fetch member meta for sent ids
func FetchMemberMeta(c *gin.Context, member_ids []string) (bool, *MemberMetaResponse) {

	//Authorize User
	userId := GetRequestingUserId(c)
	if userId == "" {
		return false, nil
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

	var membersMetaResponse MemberMetaResponse
	if err := json.Unmarshal(respBytes, &membersMetaResponse); err != nil {
		//Internal unmarshal error
		utils.GeneralAPIError(c, err.Error())
		return false, nil
	}

	return true, &membersMetaResponse
}
