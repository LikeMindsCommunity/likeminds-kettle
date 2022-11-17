package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type QuestionAnswer struct {
	QuestionId interface{} `json:"question_id"`
	Answer     string      `json:"answer"`
}

type MemberProfileRequest struct {
	QuestionAnswers []QuestionAnswer `json:"question_answers"`
	ImageUrl        string           `json:"image_url"`
}

//GetProfile is used to get member profile
func GetMemberProfile(c *gin.Context) {
	Profile(c, utils.GETMethod)
}

//EditProfile is used to update a member profile in community
func EditMemberProfile(c *gin.Context) {
	Profile(c, utils.PUTMethod)
}

//Profile method handles community member profile
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
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, FetchMemberProfileEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), requestParams, nil)

	case utils.PUTMethod:

		//Body to be sent in the edit profile api internally
		memberProfileRequest, err := parseMemberProfileRequest(c)
		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		//Send Request
		utils.SendRequest(c, utils.CoreService, EditMemberProfileEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, memberProfileRequest)

	}
}

func parseMemberProfileRequest(c *gin.Context) (*MemberProfileRequest, error) {
	//POST body params
	var mpr MemberProfileRequest

	if err := c.ShouldBindJSON(&mpr); err != nil {
		return nil, err
	}

	return &mpr, nil
}
