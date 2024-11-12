package requests

// CreateNotificationActivityRequest is used to create activity for a user | POST /user/activity
type CreateNotificationActivityRequest struct {
	Action       string   `json:"action" binding:"required"`
	ActionBy     string   `json:"action_by" binding:"required"`
	ActionOn     []string `json:"action_on" binding:"required"`
	EntityType   string   `json:"entity_type" binding:"required"`
	EntityId     string   `json:"entity_id"`
	ActivityText string   `json:"activity_text"`
}
