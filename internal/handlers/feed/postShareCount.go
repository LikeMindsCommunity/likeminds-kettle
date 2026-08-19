package feed

import (
	"fmt"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/handlers/user"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/gin-gonic/gin"
)

// PostShareCountRequest struct
type PostShareCountRequest struct {
	CountNumberType string `json:"count_number_type"`
	ShareNumber     int    `json:"share_number"`
}

// UpdatePostShareCountEndPoint is the endpoint for updating the post share count
func UpdatePostShareCount(c *gin.Context) {
	PostShareCount(c, utils.PUTMethod)
}

// PostShareCount is used to call API's based on the method
func PostShareCount(c *gin.Context, method int) {
	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	switch method {
	case utils.PUTMethod:
		postId := c.Param(ParamPostId)
		updatePostShareCountInteral(c, userId, postId)
	}
}

// Internal method for updating the post share count
func updatePostShareCountInteral(c *gin.Context, userId string, postId string) {
	//Parse the request
	postShareCountRequest, err := parsePostShareCountRequest(c)
	if err != nil {
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	updatePostShareCountEndPoint := fmt.Sprintf(UpdatePostShareCountEndPoint, postId)

	// Send request
	utils.SendRequest(c, utils.SwarmService, updatePostShareCountEndPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, postShareCountRequest)

}

// parsePostShareCountRequest is used to parse the request body for post share count
func parsePostShareCountRequest(c *gin.Context) (*PostShareCountRequest, error) {
	//POST body params
	var pscr PostShareCountRequest

	if err := c.ShouldBindJSON(&pscr); err != nil {
		return nil, err
	}

	return &pscr, nil
}
