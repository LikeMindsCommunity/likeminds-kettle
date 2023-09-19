package poll

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/feed"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CreatePollVoteRequest struct {
	Votes []string `json:"votes"`
}

func parseCreatePollVoteRequest(c *gin.Context) (*CreatePollVoteRequest, error) {
	//POST body params
	var cpvr CreatePollVoteRequest

	if err := c.ShouldBindJSON(&cpvr); err != nil {
		return nil, err
	}

	return &cpvr, nil
}

// CreatePollVote is used to vote on a poll
func CreatePollVote(c *gin.Context) {
	PollVote(c, utils.PUTMethod)
}

// GetPollVotes is used to get votes of a specific poll
func GetPollVotes(c *gin.Context) {
	PollVote(c, utils.GETMethod)
}

// PollVote method handles poll vote objects
func PollVote(c *gin.Context, method int) {

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	//Access query params and url generation
	pollId := c.Param("poll_id")
	PollVoteEndPoint := fmt.Sprintf(PollVoteEndPoint, pollId)

	//Send request
	switch method {
	case utils.GETMethod:
		getPollVotesInternal(c, userId, PollVoteEndPoint)

	case utils.PUTMethod:
		createPollVotesInternal(c, userId, PollVoteEndPoint)

	}
}

func getPollVotesInternal(c *gin.Context, userId string, endPoint string) {
	//Params to be sent in the /poll/<poll_id>/vote request
	params := map[string]string{
		ParamVotes: c.Query(ParamVotes),
	}

	//Fetch member access to view poll votes
	success, response := user.FetchMemberAccess(c, feed.IS_MEMBER, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	respBytes, statusCode := utils.GetRequestResponse(c, utils.SwarmService, endPoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	//Validate response
	apiCR := utils.ValidateClientResponse(c, respBytes, statusCode)
	if apiCR == nil {
		return
	}

	//If flow succeeds
	dataResponse := apiCR.Response
	if value, ok := dataResponse["votes"]; ok {
		votesData := value.([]map[string]interface{})
		userIds := []string{}

		//Fetch user ids
		for _, voteData := range votesData {
			if userUniqueIds, ok := voteData["users"]; ok {
				userIds = append(userIds, userUniqueIds.([]string)...)
			}
		}

		//Fetch user data for given user_unique_ids
		userData, err := user.FetchMemberMeta(utils.CreateHeaders(c, userId), userIds)
		if err != nil {
			utils.GeneralBadRequestError(c, utils.ErrorFetchingUserData)
			return
		}

		//Update user data in dataResponse
		dataResponse["users"] = userData
	}

	//Send response
	utils.GenerateResponse(c, dataResponse)
}

func createPollVotesInternal(c *gin.Context, userId string, endPoint string) {
	//Body to be sent in the /poll/poll_id/vote PUT request
	createPollVoteRequest, err := parseCreatePollVoteRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralBadRequestError(c, err.Error())
		return
	}

	//Fetch member access to create poll vote
	success, response := user.FetchMemberAccess(c, feed.IS_MEMBER, userId)
	if !success {
		return
	}

	//If not access
	if !response.Access {
		utils.MemberAccessFailError(c)
		return
	}

	//Send Request
	utils.SendRequest(c, utils.SwarmService, endPoint, utils.PUTRequest, utils.CreateHeaders(c, userId), nil, createPollVoteRequest)
}
