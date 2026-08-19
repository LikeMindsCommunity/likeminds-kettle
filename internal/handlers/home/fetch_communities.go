package home

import (
	"net/http"
	"sync"

	"github.com/LikeMindsCommunity/likeminds-kettle/internal/constants"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils"
	"github.com/LikeMindsCommunity/likeminds-kettle/internal/utils/api_client"
	"github.com/gin-gonic/gin"
)

type FetchCommunitiesResponse struct {
	HomeCommunities interface{} `json:"my_communities"`
	Subscriptions   interface{} `json:"subscriptions"`
}

const CommunitiesEndPoint = "/api/community_member/home_communities?page="
const SubscriptionEndPoint = "/api/subscription/fetch"
const ParamPage = "page"
const ResponseMyCommunities = "home_communities"
const ResponseSubscriptions = "subscriptions"

// FetchCommunities is used to blacklist LTM and RTM tokens
func FetchCommunities(c *gin.Context) {

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(constants.ParamLTM).(*constants.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralBadRequestError(c, utils.ErrorInvalidLTM)
		return
	}

	//GET Request params
	page := c.Query(ParamPage)
	if page == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	apiClient := api_client.NewAPIClient()
	wg := sync.WaitGroup{}
	wg.Add(2)
	headers := make(map[string]interface{})
	headers[utils.HeadersMemberId] = ltm.UserUniqueID
	resp := utils.Response{}
	resp.Data = make(map[string]interface{})

	utils.SafeGo(func() {
		respBytes, _, err := apiClient.GetRequest(&api_client.GetRequestOptions{
			Url:           apiClient.CoreServiceBaseURL + CommunitiesEndPoint + page,
			CustomHeaders: headers,
		})
		if err != nil {
			resp.ErrorMessage = err.Error()
			wg.Done()
		}
		var apiCR api_client.APIClientResponse
		err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
		if err != nil {
			resp.ErrorMessage = err.Error()
			wg.Done()
		}
		resp.Data.(map[string]interface{})[ResponseMyCommunities] = apiCR.Response
		wg.Done()
	})
	utils.SafeGo(func() {
		respBytes, _, err := apiClient.GetRequest(&api_client.GetRequestOptions{
			Url:           apiClient.SubscriptionServiceBaseURL + SubscriptionEndPoint,
			CustomHeaders: headers,
		})
		if err != nil {
			resp.ErrorMessage = err.Error()
			wg.Done()
		}
		var apiCR api_client.APIClientResponse
		err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
		if err != nil {
			resp.ErrorMessage = err.Error()
			wg.Done()
		}
		resp.Data.(map[string]interface{})[ResponseSubscriptions] = apiCR.Response[ResponseSubscriptions]
		wg.Done()
	})

	wg.Wait()

	if resp.ErrorMessage != "" {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
}
