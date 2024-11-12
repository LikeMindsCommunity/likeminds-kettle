package constants

const CommunityReportV1EndPoint = "/api/community/report"

// Swarm API endpoints
const (
	// NotificationFeedEndpoint | swarm API endpoint for user activity list
	NotificationFeedEndpoint = "/user/activity"

	// NotificationFeedUnreadCountEndPoint | swarm API endpoint for user activity unread count
	NotificationFeedUnreadCountEndPoint = "/user/activity/unread_count"

	// NotificationActivityMarkReadEndPoint | swarm API endpoint for user activity mark read
	NotificationActivityMarkReadEndPoint = "/user/activity/%s/mark_read"

	// UserProfileActivityEndPoint | swarm API endpoint for user profile activity
	UserProfileActivityEndPoint = "/user/%s/activity"
)
