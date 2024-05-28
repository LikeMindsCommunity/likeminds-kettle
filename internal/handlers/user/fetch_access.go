package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type MemberAccessResponse struct {
	Access bool `json:"access"`
	IsCm   bool `json:"is_cm"`
}

// FetchMemberAccess | fetch member access for sent action
func FetchMemberAccess(c *gin.Context, accessType string, userId string) (bool, *MemberAccessResponse) {

	//Params to be sent in the api/community_member/fetch_access request
	params := map[string]string{
		ParamAccessType: accessType,
	}

	//Params Validation
	if params[ParamAccessType] == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return false, nil
	}

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchUserAccessEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return false, nil
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	response := MemberAccessResponse{
		Access: false,
		IsCm:   false,
	}

	//Create Data
	if value, ok := dataResponse["access"]; ok {
		response.Access = value.(bool)
	}

	if value, ok := dataResponse["is_cm"]; ok {
		response.IsCm = value.(bool)
	}

	return true, &response
}
