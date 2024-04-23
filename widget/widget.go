package widget

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/feed"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreateWidgetRequest struct {
	ParentEntityID   string                 `json:"parent_entity_id" binding:"required"`
	ParentEntityType string                 `json:"parent_entity_type" binding:"required"`
	MetaData         map[string]interface{} `json:"metadata"`
}

type EditWidgetRequest struct {
	MetaData map[string]interface{} `json:"metadata"`
}

func parseCreateWidgetRequest(c *gin.Context) (*CreateWidgetRequest, error) {
	//POST body params
	var cwr CreateWidgetRequest

	if err := c.ShouldBindJSON(&cwr); err != nil {
		return nil, err
	}

	return &cwr, nil
}
func parseEditWidgetRequest(c *gin.Context) (*EditWidgetRequest, error) {
	//POST body params
	var ewr EditWidgetRequest

	if err := c.ShouldBindJSON(&ewr); err != nil {
		return nil, err
	}

	return &ewr, nil
}

// CreateWidget is used to create a new widget
func CreateWidget(c *gin.Context) {
	Widget(c, utils.POSTMethod)
}

// GetWidget is used to get widgets of a specific community
func GetWidget(c *gin.Context) {
	Widget(c, utils.GETMethod)
}

// EditWidget is used to edit an existing widget
func EditWidget(c *gin.Context) {
	Widget(c, utils.PUTMethod)
}

// Widget method handles widget objects
func Widget(c *gin.Context, method int) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send request
	switch method {
	case utils.GETMethod:
		GetWidgetInternal(c, userId)

	case utils.POSTMethod:
		createWidgetInternal(c, userId)

	case utils.PUTMethod:
		editWidgetInternal(c, userId)
	}
}

func GetWidgetInternal(c *gin.Context, userId string) {
	//Params to be sent in the /widget GET request
	params := map[string]string{
		ParamPage:             c.Query(ParamPage),
		ParamPageSize:         c.Query(ParamPageSize),
		ParamSearchKey:        c.Query(ParamSearchKey),
		ParamSearchValue:      c.Query(ParamSearchValue),
		ParamWidgetIds:        c.Query(ParamWidgetIds),
		ParamParentEntityId:   c.Query(ParamParentEntityId),
		ParamParentEntityType: c.Query(ParamParentEntityType),
	}

	//Fetch member access to view widgets
	success, response := user.FetchMemberAccess(c, feed.IS_MEMBER, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//add CM role in headers if user is cm
	if response.IsCm {
		headers := map[string]string{
			utils.HeaderMemberRole: utils.CMRole,
		}

		utils.AddHeaders(c, headers)
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, WidgetEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)
}

func createWidgetInternal(c *gin.Context, userId string) {
	//Body to be sent in the /widget POST request
	createWidgetRequest, err := parseCreateWidgetRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Fetch member access to create widget
	success, response := user.FetchMemberAccess(c, feed.IS_MEMBER, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//add CM role in headers if user is cm
	if response.IsCm {
		headers := map[string]string{
			utils.HeaderMemberRole: utils.CMRole,
		}

		utils.AddHeaders(c, headers)
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, WidgetEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createWidgetRequest)
}

func editWidgetInternal(c *gin.Context, userId string) {
	widgetId := c.Param("widget_id")
	EditWidgetEndPoint := fmt.Sprintf(SingleWidgetEndPoint, widgetId)

	//Body to be sent in the /widget PUT request
	editWidgetRequest, err := parseEditWidgetRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Fetch member access to edit widget
	success, response := user.FetchMemberAccess(c, feed.IS_MEMBER, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//add CM role in headers if user is cm
	if response.IsCm {
		headers := map[string]string{
			utils.HeaderMemberRole: utils.CMRole,
		}

		utils.AddHeaders(c, headers)
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, EditWidgetEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, editWidgetRequest)
}
