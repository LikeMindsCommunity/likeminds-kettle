package moderation

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/handlers/utility"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type Right struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	SubTitle   string `json:"sub_title"`
	State      int32  `json:"state"`
	IsSelected bool   `json:"is_selected"`
	IsLocked   bool   `json:"is_locked"`
}

type RightsRequest struct {
	UserId      interface{} `json:"user_id"`
	UUID        string      `json:"uuid"`
	CustomTitle *string     `json:"custom_title,omitempty"`
	Rights      []Right     `json:"rights,omitempty"`
	IsCM        bool        `json:"is_cm"`
}

type CreateActivityRequest struct {
	Action string `json:"action"`
}

func parseRightsRequest(c *gin.Context) (*RightsRequest, error) {
	//POST body params
	var rr RightsRequest

	if err := c.ShouldBindJSON(&rr); err != nil {
		return nil, err
	}

	return &rr, nil
}

// EditRights is used to edit community rights for members
func EditRights(c *gin.Context) {
	Rights(c, utils.PUTMethod)
}

// GetRights is used to get community rights for members
func GetRights(c *gin.Context) {
	Rights(c, utils.GETMethod)
}

// UpdateRights is used to update only rights sent in the request
func UpdateRights(c *gin.Context) {
	Rights(c, utils.PatchMethod)
}

// Rigths method handles community rights for members
func Rights(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Send request
	switch method {
	case utils.GETMethod:

		getRightsInternal(c, userId)

	case utils.PUTMethod:

		editRightsInternal(c, userId)

	case utils.PatchMethod:

		updateRightsInternal(c, userId)

	}
}

func getRightsInternal(c *gin.Context, userId string) {

	//Params to be sent in the fetch rights request
	params := map[string]string{
		ParamUserId: c.Query(ParamUserId),
		ParamUUID:   c.Query(ParamUUID),
	}

	//GET Request params
	isCm := c.Query(ParamIsCm)

	if isCm == "" || isCm == "false" {
		//If is_cm is missing or false, call fetch member rights api internally

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchMemberRights, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, true)

	} else {
		//else, call fetch cm rights api internally

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchCMRights, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, true)

	}
}

func editRightsInternal(c *gin.Context, userId string) {

	//Body to be sent in the update rights POST request
	rightsRequest, err := parseRightsRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	isCm := rightsRequest.IsCM

	if !isCm {
		//If isCm is missing or false, call update member rights api internally

		//Send Request
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, UpdateMemberRights, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, rightsRequest)
		if respBytes == nil {
			return
		}

		//Validate response
		apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
		if apiCR == nil {
			return
		}

		// Get user_unique_id from user_id internally and update user_id
		if rightsRequest.UUID != "" {
			uuid, _ := utility.GetUUIDInternally(utils.CreateHeaders(c, userId), rightsRequest.UUID)
			rightsRequest.UserId = uuid
		}

		//If rights change for FEED, create feed rights activity
		if reflect.TypeOf(rightsRequest.UserId).Kind() == reflect.String {
			err := createFeedRightsAcitivity(c, userId, rightsRequest.Rights, rightsRequest.UserId.(string))
			if err {
				return
			}
		}

		//Generate response
		utils.GenerateResponse(c, apiCR.Response, false)

	} else {
		//else, call update cm rights api internally

		//Send Request
		utils.SendRequest(c, utils.CoreService, UpdateCMRights, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, rightsRequest)
	}
}

func updateRightsInternal(c *gin.Context, userId string) {

	//Body to be sent in the update rights PATCH request
	rightsRequest, err := parseRightsRequest(c)
	if err != nil {
		//If body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	isCm := rightsRequest.IsCM

	if !isCm {
		//If isCm is missing or false, call update member rights api internally

		//Send Request
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, UpdateMemberRights, utils.PATCHRequest, utils.CreateHeaders(c, userId), nil, rightsRequest)
		if respBytes == nil {
			return
		}

		//Validate response
		apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
		if apiCR == nil {
			return
		}

		// Get UUID from internal core service and update UUID
		if rightsRequest.UUID != "" {
			rightsRequest.UUID, _ = utility.GetUUIDInternally(utils.CreateHeaders(c, userId), rightsRequest.UUID)
		}

		err := createFeedRightsAcitivity(c, userId, rightsRequest.Rights, rightsRequest.UUID)
		if err {
			return
		}

		//Generate response
		utils.GenerateResponse(c, apiCR.Response, false)

	} else {

		//Send Request to api/update_community_manager_rights
		utils.SendRequest(c, utils.CoreService, UpdateCMRights, utils.PATCHRequest, utils.CreateHeaders(c, userId), nil, rightsRequest)
	}
}

func extractActionFromRights(rights []Right) (string, string) {

	postAction, commentAction := "", ""
	for _, right := range rights {
		if right.Id == CREATE_POST_RIGHT_ID {

			if right.IsSelected {
				postAction = CREATE_POST_PERMISSION_ADDED_ACTION
			} else {
				postAction = CREATE_POST_PERMISSION_REMOVED_ACTION
			}
		}

		if right.Id == COMMENT_AND_REPLY_RIGHT_ID {

			if right.IsSelected {
				commentAction = CREATE_COMMENT_PERMISSION_ADDED_ACTION
			} else {
				commentAction = CREATE_COMMENT_PERMISSION_REMOVED_ACTION
			}
		}
	}

	return postAction, commentAction
}

// createFeedRightsAcitivity is used to create feed rights activity for members
func createFeedRightsAcitivity(c *gin.Context, userId string, rights []Right, uuid string) bool {

	postAction, commentAction := extractActionFromRights(rights)

	if postAction != "" {
		createActivityRequest := CreateActivityRequest{
			Action: postAction,
		}

		//Send Request
		respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, fmt.Sprintf(CreateFeedActivityEndpoint, uuid), utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createActivityRequest)
		if respBytes == nil {
			return true
		}

		//Validate response
		apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
		if apiCR == nil {
			return true
		}
	}

	if commentAction != "" {
		createActivityRequest := CreateActivityRequest{
			Action: commentAction,
		}

		//Send Request
		respBytes, _ := utils.GetRequestResponse(c, utils.SwarmService, fmt.Sprintf(CreateFeedActivityEndpoint, uuid), utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createActivityRequest)
		if respBytes == nil {
			return true
		}
	}

	return false
}
