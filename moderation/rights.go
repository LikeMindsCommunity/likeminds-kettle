package moderation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
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
	UserId      int64   `json:"user_id"`
	CommunityId int64   `json:"community_id"`
	CustomTitle string  `json:"custom_title"`
	Rights      []Right `json:"rights"`
	IsCM        bool    `json:"is_cm"`
}

//EditRights is used to edit community rights for members
func EditRights(c *gin.Context) {
	Rights(c, utils.PUTMethod)
}

//GetRights is used to get community rights for members
func GetRights(c *gin.Context) {
	Rights(c, utils.GETMethod)
}

//Rigths method handles community rights for members
func Rights(c *gin.Context, method int) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Call GET api/bot to get bot
	response := user.GetBotResponse(c, utils.GETMethod)
	if response == nil {
		return
	}

	//Send request
	var respBytes []byte
	var err error
	switch method {
	case utils.GETMethod:

		var options api_client.GetRequestOptions

		//Params to be sent in the fetch rights request
		params := map[string]string{
			ParamCommunityId: c.Query(ParamCommunityId),
			ParamUserId:      c.Query(ParamUserId),
		}

		//GET Request params
		is_cm := c.Query(ParamIsCm)

		if is_cm == "" || is_cm == "false" {
			//If is_cm is missing or false, call api/fetch_member_rights api internally

			options = api_client.GetRequestOptions{
				Url:           client.CoreServiceBaseURL + FetchMemberRights,
				CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
				Params:        params,
			}

		} else {
			//else, call api/fetch_community_manager_rights api internally

			options = api_client.GetRequestOptions{
				Url:           client.CoreServiceBaseURL + FetchCMRights,
				CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
				Params:        params,
			}
		}

		respBytes, err = client.GetRequest(&options)
		if err != nil {
			//If API fails or any other error
			utils.GeneralAPIError(c, err.Error())
			return
		}

	case utils.PUTMethod:

		var options api_client.PostRequestOptions

		//Body to be sent in the update rights POST request
		rightsRequest, err := parseRightsRequest(c)

		if err != nil {
			//If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		is_cm := rightsRequest.IsCM

		if !is_cm {
			//If is_cm is missing or false, call api/update_member_rights api internally

			options = api_client.PostRequestOptions{
				Url:           client.CoreServiceBaseURL + UpdateMemberRights,
				Body:          rightsRequest,
				CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
			}

		} else {
			//else, call api/update_community_manager_rights api internally

			options = api_client.PostRequestOptions{
				Url:           client.CoreServiceBaseURL + UpdateCMRights,
				Body:          rightsRequest,
				CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
			}
		}

		respBytes, err = client.PostRequest(&options)

		if err != nil {
			//If API fails or any other error
			utils.GeneralAPIError(c, err.Error())
			return
		}
	}

	//Parse response
	utils.ParseResponse(c, respBytes)
}

func parseRightsRequest(c *gin.Context) (*RightsRequest, error) {
	//POST body params
	var rr RightsRequest

	if err := c.ShouldBindJSON(&rr); err != nil {
		return nil, err
	}

	return &rr, nil
}
