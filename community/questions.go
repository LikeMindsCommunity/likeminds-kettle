package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/api_client"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type Question struct {
	Id                  int32  `json:"id"`
	QuestionTitle       string `json:"question_title"`
	Value               string `json:"value"`
	Optional            bool   `json:"optional"`
	State               int32  `json:"state"`
	HelpText            string `json:"help_text"`
	IsHidden            bool   `json:"is_hidden"`
	Field               bool   `json:"field"`
	Rank                int32  `json:"rank"`
	QuestionChangeState int32  `json:"question_change_state"`
	CanAddOptions       bool   `json:"can_add_options"`
	IsCompulsory        bool   `json:"is_compulsory"`
}

//EditQuestionsRequest
type EditQuestionsRequest struct {
	Questions []Question `json:"questions"`
}

//EditQuestions is used to edit Community Questions
func EditQuestions(c *gin.Context) {

	//Create internal API client
	client := api_client.NewAPIClient()

	//Call GET api/bot to get bot
	response := user.GetBotResponse(c, utils.GETMethod)
	if response == nil {
		return
	}

	//Body to be sent in the api/community/edit_questions POST request
	editQuestionsRequest, err := parseEditQuestionRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	options := api_client.PostRequestOptions{
		Url:           client.CoreServiceBaseURL + EditQuestionsEndPoint,
		Body:          editQuestionsRequest,
		CustomHeaders: utils.CreateHeaders(c, user.GetUserUniqueIDFromResponse(response)),
	}

	respBytes, err := client.PostRequest(&options, api_client.BodyTypeRaw)
	if err != nil {
		//If API fails or any other error
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Parse Response
	utils.ParseResponse(c, respBytes)
}

func parseEditQuestionRequest(c *gin.Context) (*EditQuestionsRequest, error) {
	//POST body params
	var eqr EditQuestionsRequest

	if err := c.ShouldBindJSON(&eqr); err != nil {
		return nil, err
	}

	return &eqr, nil
}
