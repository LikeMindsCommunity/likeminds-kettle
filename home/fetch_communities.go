package home

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
	"sync"
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

//FetchCommunities is used to blacklist LTM and RTM tokens
func FetchCommunities(c *gin.Context) {
	//Check if request has valid login token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.SomethingWentWrongError(c)
		return
	}
	//GET Request params
	page := c.Query(ParamPage)
	if page == "" {
		//If GET params are missing
		utils.GETQueryParamsMissingError(c)
		return
	}

	//GET Request params
	page := c.Query("page")
	if page == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success:      false,
			ErrorMessage: "Query params missing!",
		})
		return
	}

	apiClient := api_client.NewAPIClient()
	wg := sync.WaitGroup{}
	wg.Add(2)
	headers := make(map[string]interface{})
	headers[utils.HeadersMemberId] = ltm.UserID
	resp := utils.Response{}
	resp.Data = make(map[string]interface{})

	go func() {
		respBytes, err := apiClient.GetRequest(&api_client.GetRequestOptions{
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
			resp.ErrorMessage = "Something went wrong! Please try after sometime"
			wg.Done()
		}
		resp.Data.(map[string]interface{})[ResponseMyCommunities] = apiCR.Response
		wg.Done()
	}()
	go func() {
		respBytes, err := apiClient.GetRequest(&api_client.GetRequestOptions{
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
			resp.ErrorMessage = "Something went wrong! Please try after sometime"
			wg.Done()
		}
		resp.Data.(map[string]interface{})[ResponseSubscriptions] = apiCR.Response[ResponseSubscriptions]
		wg.Done()
	}()

	wg.Wait()

	if resp.ErrorMessage != "" {
		c.JSON(http.StatusInternalServerError, resp)
	}
	c.JSON(http.StatusOK, resp)
}
