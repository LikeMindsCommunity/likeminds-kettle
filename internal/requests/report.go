package requests

type CloseReportsNewRequest struct {
	ReportIds []int  `json:"report_ids" binding:"required"`
	Status    string `json:"status,omitempty"`
}
