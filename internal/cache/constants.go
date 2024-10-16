package cache

// Params
const (
	ParamRedisClient = "redis_client"
)

// Cache Keys
const (
	CommunityConfigurationsCacheKey    = "%s_community_configurations"
	ProfileMetaConfigurationsCacheKey  = "%d_profile_meta_configurations"  // communityId
	FeedMetadataConfigurationsCacheKey = "%d_feed_metadata_configurations" // communityId
	CommunityIdAgainstApiKeyCacheKey   = "%s_community_id"                 // apiKey
	CommunitySettingsCacheKey          = "%d_community_settings"           // communityId
	UserMetaCacheKey                   = "%d_%s_user_meta"                 // communityId userUniqueId
	TopicMetaCacheKey                  = "%d_%s_topic_meta"                // communityId topicId
	UserTopicsCacheKey                 = "%d_%s_user_topics"               // communityId userUniqueId
	WidgetMetaCacheKey                 = "%d_%s_widget_meta"               // communityId widgetId
	CommunityBillingDataKey            = "%d_community_billing_data"       // communityId
	TierDataKey                        = "%d_tier_data"
	ChatroomParticipantsKey            = "chatroom_participants_%s" // chatroom participants
)

// Cache TTLs in hours
const (
	ProfileMetaConfigurationsCacheTTL = 7 * 24
	CommunityIdAgainstApiKeyCacheTTL  = 30 * 24
	CommunitySettingsCacheTTL         = 30 * 24
	UserMetaCacheTTL                  = 7 * 24
	TopicMetaCacheTTL                 = 7 * 24
	UserTopicsCacheTTL                = 7 * 24
	WidgetMetaCacheTTL                = 7 * 24
	CommunityBillingDataTTL           = 7 * 24
	TierDataTTL                       = 30 * 24
	ChatroomParticipantsTTL           = 7 * 24
	ChatroomTTL                       = 30 * 24
)
