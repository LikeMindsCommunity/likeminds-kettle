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

//Fetch is used to blacklist LTM and RTM tokens
func Fetch(c *gin.Context) {
	_, ok := c.MustGet("ltm").(*token.LoginTokenMeta)
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

	go func() {
		homeCommunities, errHomeCommunities := apiClient.GetRequest(&core_client.GetRequestOptions{
			Url:    "/api/community_member/home_communities",
			Header: c.Request.Header,
		})
		if errHomeCommunities != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success:      false,
				ErrorMessage: errHomeCommunities.Error(),
			})
			return
		}
		fetchCommunities.HomeCommunities = homeCommunities
		wg.Done()
	}()
	go func() {
		subscriptions, errSubscription := apiClient.GetRequest(&core_client.GetRequestOptions{
			Url:    "/api/subscription/fetch",
			Header: c.Request.Header,
		})
		if errSubscription != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Success:      false,
				ErrorMessage: errSubscription.Error(),
			})
			fetchCommunities.Subscriptions = subscriptions
			return
		}
		wg.Done()
	}()

	c.JSON(http.StatusOK, utils.Response{
		Success: true,
		Data:    fetchCommunities,
	})
}
