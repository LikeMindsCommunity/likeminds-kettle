package community

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
	"github.com/nateshr/likeminds-authentication/widget"
)

type QuestionAnswer struct {
	QuestionId interface{} `json:"question_id"`
	Answer     string      `json:"answer"`
}

type MemberProfileRequest struct {
	QuestionAnswers []QuestionAnswer       `json:"question_answers"`
	ImageUrl        string                 `json:"image_url"`
	Name            *string                `json:"name"`
	WidgetId        *string                `json:"-"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// GetProfile is used to get member profile
func GetMemberProfile(c *gin.Context) {
	Profile(c, utils.GETMethod)
}

// EditProfile is used to update a member profile in community
func EditMemberProfile(c *gin.Context) {
	Profile(c, utils.PUTMethod)
}

// Profile method handles community member profile
func Profile(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Send request
	switch method {
	case utils.GETMethod:

		// Params to be sent in the api/community_member/fetch_profile request
		requestParams := map[string]string{
			ParamUserId: c.Query(ParamUserId),
			ParamUUID:   c.Query(ParamUUID),
		}

		//Get Request response
		respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, FetchMemberProfileEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)
		if respBytes == nil {
			return
		}

		//Parse and generate response
		utils.ParseResponse(c, respBytes, statusCode, true)

	case utils.PUTMethod:

		EditMemberProfileInternal(c, userId)
	}
}

func parseMemberProfileRequest(c *gin.Context) (*MemberProfileRequest, error) {
	// POST body params
	var mpr MemberProfileRequest

	if err := c.ShouldBindJSON(&mpr); err != nil {
		return nil, err
	}

	for i, QuestionAnswer := range mpr.QuestionAnswers {
		if QuestionAnswer.QuestionId != nil {
			mpr.QuestionAnswers[i].QuestionId = utils.ParseInterfaceToString(QuestionAnswer.QuestionId)
		}
	}

	return &mpr, nil
}

func EditMemberProfileInternal(c *gin.Context, userId string) {

	//Body to be sent in the edit profile api internally
	mpr, err := parseMemberProfileRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	widgetId, created := "", false

	// If metadata is present, create or edit widget
	if mpr.Metadata != nil {

		widgetId, created = createOrUpdateWidgetForMemberProfile(c, userId, mpr.Metadata)
		if widgetId != "" {
			mpr.WidgetId = &widgetId
		}
	}

	//Send Request to edit member profile
	respBytes, statusCode := utils.GetRequestResponse(c, utils.CoreService, EditMemberProfileEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, mpr)

	apiCr := utils.ValidateClientResponse(c, respBytes, statusCode)

	// if a widget was created and edit_member API failed, delete the widget
	if created && apiCr == nil && widgetId != "" {

		deleteWidgetEndPoint := fmt.Sprintf(widget.SingleWidgetEndPoint, widgetId)

		// Send request to delete widget
		respBytes, _ := utils.GetRequestResponse(c, utils.SwarmService, deleteWidgetEndPoint, utils.DELETERequest, utils.CreateHeaders(c, userId), nil, nil)

		apiCr := utils.ValidateClientResponse(c, respBytes, statusCode)
		if apiCr == nil {
			return
		}

		return
	}

	utils.GenerateResponse(c, apiCr.Response, false)

}

func createOrUpdateWidgetForMemberProfile(c *gin.Context, userId string, metaData map[string]interface{}) (string, bool) {

	widgetId, created := "", false

	// send request and check if widget exists
	fetchWidgetParams := map[string]string{
		widget.ParamParentEntityId:   userId,
		widget.ParamParentEntityType: widget.ParentEntityTypeUser,
	}

	//Send Request to /widget GET
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, widget.WidgetEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), fetchWidgetParams, nil)
	apiCr := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCr == nil {
		return widgetId, created
	}

	// Get widget id from response
	dataResponse := apiCr.Response
	if widgets, ok := dataResponse["widgets"].([]interface{}); ok {
		if len(widgets) > 0 {
			if id, ok := widgets[0].(map[string]interface{})["_id"].(string); ok {
				widgetId = id
			}
		}
	}

	if widgetId != "" { // If widget exists, edit widget
		EditWidgetEndPoint := fmt.Sprintf(widget.SingleWidgetEndPoint, widgetId)

		ewr := widget.EditWidgetRequest{
			MetaData: metaData,
		}

		// Send request to edit widget
		respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, EditWidgetEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, ewr)
		apiCr := utils.ValidateClientResponse(c, respBytes, statusCode)
		if apiCr == nil {
			return widgetId, created
		}

	} else { // If widget does not exist, create widget

		cwr := widget.CreateWidgetRequest{
			ParentEntityID:   userId,
			ParentEntityType: widget.ParentEntityTypeUser,
			MetaData:         metaData,
		}

		// Send request to create widget
		respByte, statusCode := utils.GetRequestResponse(c, utils.SwarmService, widget.WidgetEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, cwr)
		apiCr := utils.ValidateClientResponse(c, respByte, statusCode)
		if apiCr == nil {
			return widgetId, created
		}

		// Get widget id from response
		dataResponse := apiCr.Response
		if widget, ok := dataResponse["widget"].(map[string]interface{}); ok {
			widgetId = widget["_id"].(string)
			created = true
		}
	}

	return widgetId, created
}
