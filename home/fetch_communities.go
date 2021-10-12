package home

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
	"sync"
)

type FetchCommunjitiesResponse struct {
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
	//GET Request params
	page := c.Query("page")
	if page == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			Success:      false,
			ErrorMessage: "Query params missing!",
		})
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
		wg.Done()
	}()

	wg.Wait()

	if resp.ErrorMessage != "" {
		c.JSON(http.StatusInternalServerError, resp)
	}
	c.JSON(http.StatusOK, resp)
}
