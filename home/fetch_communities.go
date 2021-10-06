package home

import (
	"github.com/gin-gonic/gin"
<<<<<<< HEAD
	"github.com/nateshr/likeminds-authentication/api_client"
=======
	"github.com/nateshr/likeminds-authentication/core_client"
>>>>>>> fetch_communitites tested
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
	"sync"
)

<<<<<<< HEAD
type FetchCommunjitiesResponse struct {
=======
type FetchCommunitiesResponse struct {
>>>>>>> fetch_communitites tested
	HomeCommunities interface{} `json:"my_communities"`
	Subscriptions   interface{} `json:"subscriptions"`
}

//FetchCommunities is used to blacklist LTM and RTM tokens
func FetchCommunities(c *gin.Context) {
	ltm, ok := c.MustGet("ltm").(*token.LoginTokenMeta)
	if !ok {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Success:      false,
			ErrorMessage: "Something went wrong! Please try after sometime",
		})
		return
	}
<<<<<<< HEAD
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
	headers["x-member-id"] = ltm.UserID
	resp := utils.Response{}

	go func() {
		respBytes, err := apiClient.GetRequest(&api_client.GetRequestOptions{
			Url:           apiClient.CoreServiceBaseURL + "/api/community_member/home_communities?page=" + page,
			CustomHeaders: headers,
		})
		if err != nil {
			resp.ErrorMessage = err.Error()
			wg.Done()
			return
		}
		var apiCR api_client.APIClientResponse
		err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
		if err != nil {
			resp.ErrorMessage = "Something went wrong! Please try after sometime"
			wg.Done()
		}
		resp.Data.(map[string]interface{})["my_communities"] = apiCR.Response["my_communities"]
		wg.Done()
	}()
	go func() {
		respBytes, err := apiClient.GetRequest(&api_client.GetRequestOptions{
			Url:           apiClient.SubscriptionServiceBaseURL + "/api/subscription/fetch",
			CustomHeaders: headers,
		})
		if err != nil {
			resp.ErrorMessage = err.Error()
			wg.Done()
			return
		}
		var apiCR api_client.APIClientResponse
		err = api_client.UnmarshalAPIClientResponse(respBytes, &apiCR)
		if err != nil {
			resp.ErrorMessage = "Something went wrong! Please try after sometime"
			wg.Done()
			return
		}
		resp.Data.(map[string]interface{})["subscriptions"] = apiCR.Response["subscriptions"]
=======

	apiClient := core_client.NewClient()
	wg := sync.WaitGroup{}
	wg.Add(2)
	fetchCommunities := FetchCommunitiesResponse{}
	headers := make(map[string]string)
	headers["x-member-id"] = ltm.UserID

	go func() {
		homeCommunities, errHomeCommunities := apiClient.GetRequest(&core_client.GetRequestOptions{
			Url:           apiClient.CoreServiceBaseURL + "/api/community_member/home_communities?page=0",
			CustomHeaders: headers,
		})
		if errHomeCommunities != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success:      false,
				ErrorMessage: errHomeCommunities.Error(),
			})
			wg.Done()
			return
		}
		fetchCommunities.HomeCommunities = homeCommunities
		wg.Done()
	}()
	go func() {
		subscriptions, errSubscription := apiClient.GetRequest(&core_client.GetRequestOptions{
			Url:           apiClient.SubscriptionServiceBaseURL + "/api/subscription/fetch",
			CustomHeaders: headers,
		})
		if errSubscription != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success:      false,
				ErrorMessage: errSubscription.Error(),
			})
			wg.Done()
			return
		}
		fetchCommunities.Subscriptions = subscriptions
>>>>>>> fetch_communitites tested
		wg.Done()
	}()

	wg.Wait()
<<<<<<< HEAD

	if resp.ErrorMessage != "" {
		c.JSON(http.StatusInternalServerError, resp)
	}
	c.JSON(http.StatusOK, resp)
=======
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    fetchCommunities,
	})
>>>>>>> fetch_communitites tested
}
