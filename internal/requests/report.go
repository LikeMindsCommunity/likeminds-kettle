package requests

type CloseReportsNewRequest struct {
	ReportIds   []int  `json:"report_ids" binding:"required"`
	ActionTaken string `json:"action_taken,omitempty"`
}
