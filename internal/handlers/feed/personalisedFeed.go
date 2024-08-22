package feed

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/internal/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/logging"
	"github.com/nateshr/likeminds-authentication/internal/utils"
)

// Exposed method to recompute personalised feed
func RecomputePersonalisedFeed(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	headers := utils.CreateHeaders(c, userId)

	// Check if personalised feed settings is enabled
	if !utils.IsPersonalisedFeedEnabled(utils.GetRedisClientFromContext(c), headers) {
		utils.GeneralBadRequestError(c, utils.PersonalisedFeedDisabledError)
		return
	}

	// Send request to /personalised/recompute
	utils.SendRequest(c, utils.SwarmService, RecomputePersonalisedFeedEndPoint, utils.POSTRequestRawBody, headers, nil, nil)
}

// Exposed method to reorder personalised feed
func ReorderPersonalisedFeed(c *gin.Context) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	headers := utils.CreateHeaders(c, userId)

	// Check if personalised feed settings is enabled
	if !utils.IsPersonalisedFeedEnabled(utils.GetRedisClientFromContext(c), headers) {
		utils.GeneralBadRequestError(c, utils.PersonalisedFeedDisabledError)
		return
	}

	// Send request to /personalised/reorder
	utils.SendRequest(c, utils.SwarmService, ReorderPersonalisedFeedEndPoint, utils.POSTRequestRawBody, headers, nil, nil)
}

// Exposed method to fetch personalised feed
func FetchPersonalisedFeed(c *gin.Context) {

	userId, access := validateFetchPersonalisedFeed(c)
	if !access {
		return
	}

	headers := utils.CreateHeaders(c, userId)

	// Check if personalised feed settings is enabled
	if !utils.IsPersonalisedFeedEnabled(utils.GetRedisClientFromContext(c), headers) {
		utils.GeneralBadRequestError(c, utils.PersonalisedFeedDisabledError)
		return
	}

	shouldReorder := utils.GetBooleanFromString(c.Query(ParamShouldReorder))
	shouldRecompute := utils.GetBooleanFromString(c.Query(ParamShouldRecompute))

	// Reorder the metrics if should_reorder is true
	if shouldReorder {
		// Send request to /personalised/reorder
		respBytes, _ := utils.GetRequestResponse(c, utils.SwarmService, ReorderPersonalisedFeedEndPoint, utils.POSTRequestRawBody, headers, nil, nil)
		if respBytes == nil {
			utils.GenerateResponse(c, nil, false)
			return
		}
	}

	// Recompute the metrics in background if should_recompute is true
	if shouldRecompute {
		// Send request to /personalised/recompute
		go func() {
			respBytes, statusCode, err := utils.GetRequestResponseWithoutContext(utils.SwarmService, RecomputePersonalisedFeedEndPoint, utils.POSTRequestRawBody, headers, nil, nil)
			if err != nil {
				logging.Error(fmt.Sprintf("Error in recomputing personalised feed: %s", err.Error()))
			}

			apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
			convertedApiCR, _ := json.Marshal(apiCR)
			logging.Error(fmt.Sprintf("Response of recompute personalised feed: %s", convertedApiCR))
		}()
	}

	params := map[string]string{
		ParamPage:     c.Query(ParamPage),
		ParamPageSize: c.Query(ParamPageSize),
	}

	// Send request to /feed/personalised
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, FetchPersonalisedFeedEndPoint, utils.GETRequest, headers, params, nil)

	// Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	dataResponse, err := utils.PopulateDataResponseForFeed(headers, utils.GetRedisClientFromContext(c), apiCR.Response)
	if err != nil {
		utils.GenerateResponse(c, nil, false)
		return
	}

	//Send response
	utils.GenerateResponse(c, dataResponse, true)
}

func validateFetchPersonalisedFeed(c *gin.Context) (string, bool) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return "", false
	}

	// Check if personalised feed settings is enabled
	if !utils.IsPersonalisedFeedEnabled(utils.GetRedisClientFromContext(c), utils.CreateHeaders(c, userId)) {
		utils.GeneralBadRequestError(c, utils.PersonalisedFeedDisabledError)
		return "", false
	}

	// Fetch member access to check whether member is CM or not
	success, response := user.FetchMemberAccess(c, VIEW_POST_ACTION, userId)
	if !success {
		return "", false
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return "", false
	}

	//add Admin role in headers if user is cm
	if response.IsCm {
		headers := map[string]string{
			utils.HeaderMemberRole: utils.CMRole,
		}

		utils.AddHeaders(c, headers)
	}

	return userId, true
}
