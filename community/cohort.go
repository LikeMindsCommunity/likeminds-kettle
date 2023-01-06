package community

import (
	"github.com/gin-gonic/gin"
	"github.com/nateshr/likeminds-authentication/user"
	"github.com/nateshr/likeminds-authentication/utils"
)

type CohortFilter struct {
	QuestionID    int    `json:"question_id" binding:"required"`
	QuestionTitle string `json:"question_title"`
	Value         string `json:"value" binding:"required"`
}

type CreateCohortRequest struct {
	Name      string         `json:"name" binding:"required"`
	MemberIDs []int          `json:"member_ids"  binding:"required"`
	Filter    []CohortFilter `json:"filter"`
}

type CohortRights struct {
	RightID    int    `json:"id" binding:"required"`
	IsLocked   bool   `json:"is_locked"`
	IsSelected bool   `json:"is_selected" binding:"required"`
	State      int    `json:"state" binding:"required"`
	Title      string `json:"title"`
}

type EditCohortRequest struct {
	CohortID   int            `json:"cohort_id" binding:"required"`
	Name       string         `json:"name"`
	MemberIDs  []int          `json:"member_ids"`
	Rights     []CohortRights `json:"rights"`
	CohortType int            `json:"type"`
	TypeID     string         `json:"type_id"`
	FilterList []CohortFilter `json:"filter"`
}

// Create a cohort in community
func CreateCohort(c *gin.Context) {
	Cohort(c, utils.POSTMethod)
}

// Get a cohort data in community
func GetCohort(c *gin.Context) {
	Cohort(c, utils.GETMethod)
}

// Delete a cohort data in community
func DeleteCohort(c *gin.Context) {
	Cohort(c, utils.DELETEMethod)
}

// Edit a cohort data in community
func EditCohort(c *gin.Context) {
	Cohort(c, utils.PUTMethod)
}

// Cohort method that will handle cohorts in community
func Cohort(c *gin.Context, method int) {

	// Authorize User
	userId := user.GetRequestingUserId(c)
	if userId == "" {
		return
	}

	botId := user.GetBotId(c)
	if botId != "" {
		userId = botId
	}

	// Send request
	switch method {
	case utils.POSTMethod:

		// Body to be sent in the create cohort api internally
		createCohortRequest, err := parseCreateCohortRequest(c)

		if err != nil {
			// If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, CreateCohortEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, createCohortRequest)

	case utils.GETMethod:

		var params map[string]string
		var urlEndpoint string = GetCommunityCohortsEndPoint

		cohortID := c.Query(ParamCohortID)

		if cohortID != "" {
			// Params to be sent in fetch cohort api internally
			params = map[string]string{
				ParamCohortID: c.Query(ParamCohortID),
			}

			urlEndpoint = GetCohortEndPoint
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, urlEndpoint, utils.GETRequest, utils.CreateHeaders(c, userId), params, nil)

	case utils.DELETEMethod:

		// Body to be sent in the delete cohort api internally
		editCohortRequest, err := parseEditCohortRequest(c)

		if err != nil {
			// If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, DeleteCohortEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, editCohortRequest)

	case utils.PUTMethod:

		// Body to be sent in the edit member api internally
		editCohortRequest, err := parseEditCohortRequest(c)

		if err != nil {
			// If POST body params are missing
			utils.GeneralAPIError(c, err.Error())
			return
		}

		// Send Request
		utils.SendRequest(c, utils.CoreService, EditCohortEndPoint, utils.POSTRequestRawBody, utils.CreateHeaders(c, userId), nil, editCohortRequest)

	}
}

func parseCreateCohortRequest(c *gin.Context) (*CreateCohortRequest, error) {

	// POST body params
	var ccr CreateCohortRequest

	if err := c.ShouldBindJSON(&ccr); err != nil {
		return nil, err
	}

	return &ccr, nil
}

func parseEditCohortRequest(c *gin.Context) (*EditCohortRequest, error) {

	// POST body params
	var ecr EditCohortRequest

	if err := c.ShouldBindJSON(&ecr); err != nil {
		return nil, err
	}

	return &ecr, nil
}
