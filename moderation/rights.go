package moderation

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
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
	UserId      interface{} `json:"user_id" binding:"required"`
	CustomTitle string      `json:"custom_title"`
	Rights      []Right     `json:"rights"`
	IsCM        bool        `json:"is_cm"`
}

type CreateActivityRequest struct {
	Action string `json:"action"`
}

// EditRights is used to edit community rights for members
func EditRights(c *gin.Context) {
	Rights(c, utils.PUTMethod)
}

// GetRights is used to get community rights for members
func GetRights(c *gin.Context) {
	Rights(c, utils.GETMethod)
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

		//Params to be sent in the fetch rights request
		params := map[string]string{
			ParamUserId: c.Query(ParamUserId),
		}

		//GET Request params
		is_cm := c.Query(ParamIsCm)

		if is_cm == "" || is_cm == "false" {
			//If is_cm is missing or false, call fetch member rights api internally

			//Send Request
			utils.SendRequest(c, utils.CoreService, FetchMemberRights, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		} else {
			//else, call fetch cm rights api internally

			//Send Request
			utils.SendRequest(c, utils.CoreService, FetchCMRights, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
		}

	case utils.PUTMethod:

		//Body to be sent in the update rights POST request
		rightsRequest, err := parseRightsRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		is_cm := rightsRequest.IsCM

		if !is_cm {
			//If is_cm is missing or false, call update member rights api internally

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

			//If flow succeeds
			if reflect.TypeOf(rightsRequest.UserId).Kind() == reflect.String {
				post_action := ""
				comment_action := ""
				for _, right := range rightsRequest.Rights {
					if right.Id == CREATE_POST_RIGHT_ID && right.IsSelected {
						post_action = CREATE_POST_PERMISSION_ADDED_ACTION
					}

					if right.Id == COMMENT_AND_REPLY_RIGHT_ID && right.IsSelected {
						comment_action = CREATE_COMMENT_PERMISSION_ADDED_ACTION
					}
				}

				if post_action == "" {
					post_action = CREATE_POST_PERMISSION_REMOVED_ACTION
				}

				if comment_action == "" {
					comment_action = CREATE_COMMENT_PERMISSION_REMOVED_ACTION
				}

				if post_action != "" {
					createActivityRequest := CreateActivityRequest{
						Action: post_action,
					}

					fmt.Println(fmt.Sprintf(CreateFeedActivityEndpoint, rightsRequest.UserId.(string)), utils.SwarmService, createActivityRequest)

					//Send Request
					respBytes, _ := utils.GetRequestResponse(c, utils.SwarmService, fmt.Sprintf(CreateFeedActivityEndpoint, rightsRequest.UserId.(string)), utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createActivityRequest)
					if respBytes == nil {
						return
					}

					//Validate response
					apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
					if apiCR == nil {
						return
					}
				}

				if comment_action != "" {
					createActivityRequest := CreateActivityRequest{
						Action: comment_action,
					}

					//Send Request
					respBytes, _ := utils.GetRequestResponse(c, utils.SwarmService, fmt.Sprintf(CreateFeedActivityEndpoint, rightsRequest.UserId.(string)), utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createActivityRequest)
					if respBytes == nil {
						return
					}
				}
			}

			//Generate response
			utils.GenerateResponse(c, apiCR.Response)

		} else {
			//else, call update cm rights api internally

			//Send Request
			utils.SendRequest(c, utils.CoreService, UpdateCMRights, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, rightsRequest)
		}
	}
}

func parseRightsRequest(c *gin.Context) (*RightsRequest, error) {
	//POST body params
	var rr RightsRequest

	if err := c.ShouldBindJSON(&rr); err != nil {
		return nil, err
	}

	return &rr, nil
}
