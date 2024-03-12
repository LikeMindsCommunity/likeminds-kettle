package cache

// Params
const (
	ParamRedisClient = "redis_client"
)

// Cache Keys
const (
	CommunityConfigurationsCacheKey   = "%s_community_configurations"
	ProfileMetaConfigurationsCacheKey = "%s_profile_meta_configurations"
	UserMetaCacheKey                  = "%s_user_meta"   // userUniqueId
	TopicMetaCacheKey                 = "%s_topic_meta"  // topicId
	UserTopicsCacheKey                = "%s_user_topics" // userUniqueId
	WidgetMetaCacheKey                = "%s_widget_meta" // widgetId
)

// Cache TTLs in hours
const (
	ProfileMetaConfigurationsCacheTTL = 6 // TODO: change it to 7 days
	UserMetaCacheTTL                  = 7 * 24
	TopicMetaCacheTTL                 = 7 * 24
	UserTopicsCacheTTL                = 7 * 24
	WidgetMetaCacheTTL                = 7 * 24
)
