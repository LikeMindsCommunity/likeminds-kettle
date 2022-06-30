package conversation

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/token"
	"github.com/nateshr/likeminds-authentication/utils"
)

type AddReactionRequest struct {
	ChatroomID     int64  `json:"chatroom_id"`
	Reaction       string `json:"reaction"`
	ConversationID int64  `json:"conversation_id"`
}

type RemoveReactionRequest struct {
	ChatroomID     int64 `json:"chatroom_id"`
	ConversationID int64 `json:"conversation_id"`
}

//AddReaction is used to add reaction to specific conversation
func AddReaction(c *gin.Context) {
	Conversation(c, utils.PUTMethod)
}

//RemoveReaction is used to delete reaction from a specific conversation
func RemoveReaction(c *gin.Context) {
	Conversation(c, utils.DELETEMethod)
}

//Reaction method handles reaction on a conversation object
func Reaction(c *gin.Context, method int) {
	//Create internal API client
	client := api_client.NewAPIClient()

	//Check if request has LTM token or not
	ltm, ok := c.MustGet(token.ParamLTM).(*token.LoginTokenMeta)
	if !ok {
		//If token is not available
		utils.GeneralAPIError(c, utils.ErrorInvalidLTM)
		return
	}

	//Send request
	var respBytes []byte
	switch method {
	case utils.PUTMethod:

		respBytes = addReactionInternal(c, client, ltm)

	case utils.DELETEMethod:

		respBytes = removeReactionInternal(c, client, ltm)
	}

	if respBytes == nil {
		return
	}

	//Parse response
	utils.ParseResponse(c, respBytes)
}

func parseAddReactionRequest(c *gin.Context) (*AddReactionRequest, error) {
	//POST body params
	var arr AddReactionRequest

	if err := c.ShouldBindJSON(&arr); err != nil {
		return nil, err
	}

	return &arr, nil
}

func parseRemoveReactionRequest(c *gin.Context) (*RemoveReactionRequest, error) {
	//POST body params
	var rrr RemoveReactionRequest

	if err := c.ShouldBindJSON(&rrr); err != nil {
		return nil, err
	}

	return &rrr, nil
}

func addReactionInternal(c *gin.Context, client *api_client.APIClient, ltm *token.LoginTokenMeta) []byte {
	//Body to be sent in the api/conversation/add_reaction POST request
	addReactionRequest, err := parseAddReactionRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + AddReactionEndPoint,
		Body:          addReactionRequest,
		CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
	}

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeFormUrlEncoded)

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}

func removeReactionInternal(c *gin.Context, client *api_client.APIClient, ltm *token.LoginTokenMeta) []byte {
	//Body to be sent in the api/conversation/remove_reaction POST request
	removeReactionRequest, err := parseRemoveReactionRequest(c)

	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + RemoveReactionEndPoint,
		Body:          removeReactionRequest,
		CustomHeaders: utils.CreateHeaders(c, ltm.UserUniqueID),
	}

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeFormUrlEncoded)

	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return nil
	}

	return respBytes
}
