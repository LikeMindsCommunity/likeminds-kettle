package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/chatroom"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

func DMStatus(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Params to be sent in the api/community_member/can_dm request
	requestParams := map[string]string{
		RequestFromParam: c.Query(RequestFromParam),
		ParamMemberId:    c.Query(ParamMemberId),
		ParamUUID:        c.Query(ParamUUID),
		ChatroomIDParam:  c.Query(ChatroomIDParam),
	}

	if requestParams[RequestFromParam] == UserChannelReqFrom {

		// set req_from to member_profile
		requestParams[RequestFromParam] = MemberProfileReqFrom

		// Internally call api/community_member/can_dm
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, UserCanDMEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)
		if respBytes == nil {
			return
		}

		//Validate response
		apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
		if apiCR == nil {
			return
		}

		dataResponse := apiCR.Response

		// if show_dm is true
		if dataResponse[showDmResponse] != nil && dataResponse[showDmResponse] == true {

			//Body to be sent in the /chatroom/create_dm POST request
			createDMbody := map[string]string{
				"member_id": c.Query(ParamMemberId),
				"uuid":      c.Query(ParamUUID),
			}

			// Send request to api/chatroom/create_dm with body
			utils.SendRequest(c, utils.CoreService, chatroom.CreateDMEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createDMbody)

		} else {

			utils.GeneralBadRequestError(c, utils.ErrorUserCannotDm)
		}

	} else {

		// send request to api/community_member/can_dm
		utils.SendRequest(c, utils.CoreService, UserCanDMEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)

	}

}
