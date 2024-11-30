package user

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

type UpdateDashboardUserRequest struct {
	Name         string `json:"name,omitempty"`
	ImageURL     string `json:"image_url,omitempty"`
	EmailID      string `json:"email_id,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	OTP          string `json:"otp,omitempty"`
}

// EditDashboardUser is used to edit user data in dashboard
func EditDashboardUser(c *gin.Context) {
	DashboardUser(c, utils.PatchMethod)
}

// DashboardUser method handles dashboard user API's
func DashboardUser(c *gin.Context, method int) {
	//Authorize User
	userId := GetRequestingUserId(c)
	if userId == "" {
		return
	}

	// Send request
	switch method {

	case utils.PatchMethod:
		editDashboardUserInternal(c, userId)
	}

}

// editDashboardUserInternal is internal method used to edit user data in dashboard
func editDashboardUserInternal(c *gin.Context, userId string) {

	//Body to be sent in the edit member api internally
	updateDashboardUserRequest, err := parseUpdateDashboardRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, DashboardUserEndpoint, utils.PATCHRequest, utils.CreateHeaders(c, userId), nil, updateDashboardUserRequest)
}

// parseUpdateDashboardRequest is used to parse the request body
func parseUpdateDashboardRequest(c *gin.Context) (*UpdateDashboardUserRequest, error) {
	var updateDashboardUserRequest UpdateDashboardUserRequest
	if err := c.ShouldBindJSON(&updateDashboardUserRequest); err != nil {
		return nil, err
	}

	return &updateDashboardUserRequest, nil
}
