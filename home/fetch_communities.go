package home

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/core_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
	"net/http"
	"sync"
)

type FetchCommunitiesResponse struct {
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
		wg.Done()
	}()

	wg.Wait()
	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    fetchCommunities,
	})
}
