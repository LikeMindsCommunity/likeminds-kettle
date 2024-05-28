package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/handlers/user"
	"github.com/nateshr/likeminds-authentication/internal/constants"
	"github.com/nateshr/likeminds-authentication/utils"
)

type Question struct {
	Id                     int32       `json:"id,omitempty"`
	QuestionTitle          string      `json:"question_title"`
	Value                  interface{} `json:"value"`
	Optional               bool        `json:"optional"`
	State                  int32       `json:"state"`
	HelpText               string      `json:"help_text"`
	IsHidden               bool        `json:"is_hidden"`
	Field                  bool        `json:"field"`
	Rank                   int32       `json:"rank"`
	ImageUrl               string      `json:"image_url,omitempty"`
	QuestionChangeState    int32       `json:"question_change_state"`
	CanAddOptions          bool        `json:"can_add_options"`
	IsCompulsory           bool        `json:"is_compulsory"`
	IsAnswerEditable       bool        `json:"is_answer_editable"`
	OptionsOnlyForSelf     *bool       `json:"options_only_for_self,omitempty"`
	Tag                    string      `json:"tag,omitempty"`
	DropdownSelectionLimit int32       `json:"dropdown_selection_limit"`
}

// EditQuestionsRequest
type EditQuestionsRequest struct {
	Questions []Question `json:"questions"`
}

// EditQuestions is used to edit Community Questions
func EditQuestions(c *gin.Context) {
	Questions(c, utils.PUTMethod)
}

// GetQuestions is used to get Community Questions
func GetQuestions(c *gin.Context) {
	Questions(c, utils.GETMethod)
}

// Questions is used to for Community Questions
func Questions(c *gin.Context, method int) {
	var userId string

	ltm, _ := c.Get(constants.ParamLTM)

	if ltm != nil {
		// Authorize User
		userId = user.GetRequestingUserId(c)
		if userId == "" {
			return
		}

		botId := user.GetBotId(c)
		if botId != "" {
			userId = botId
		}

	}

	switch method {
	case utils.PUTMethod:

		// Body to be sent in the api/community/edit_questions POST request
		editQuestionsRequest, err := parseEditQuestionRequest(c)
		if err != nil {
			// If POST body params are missing
			utils.GeneralBadRequestError(c, err.Error())
			return
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, EditQuestionsEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, editQuestionsRequest)

	case utils.GETMethod:
		// Send Request
		utils.SendRequest(c, utils.CoreService, FetchQuestionsEndPoint, utils.GETRequest, utils.CreateHeaders(c, userId), nil, nil)
	}
}

func parseEditQuestionRequest(c *gin.Context) (*EditQuestionsRequest, error) {
	// POST body params
	var eqr EditQuestionsRequest

	if err := c.ShouldBindJSON(&eqr); err != nil {
		return nil, err
	}

	return &eqr, nil
}
