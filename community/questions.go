package community

import (
	"github.com/gin-gonic/gin"
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

	//Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	//Body to be sent in the api/community/edit_questions POST request
	editQuestionsRequest, err := parseEditQuestionRequest(c)
	if err != nil {
		//If POST body params are missing
		utils.GeneralAPIError(c, err.Error())
		return
	}

	//Send Request
	utils.SendRequest(c, utils.CoreService, EditQuestionsEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, editQuestionsRequest)
}

func parseEditQuestionRequest(c *gin.Context) (*EditQuestionsRequest, error) {
	//POST body params
	var eqr EditQuestionsRequest

	if err := c.ShouldBindJSON(&eqr); err != nil {
		return nil, err
	}

	return &eqr, nil
}
