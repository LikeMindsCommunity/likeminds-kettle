package cache

// Params
const (
	ParamRedisClient = "redis_client"
)

// Cache Keys
const (
	CommunityConfigurationsCacheKey   = "%s_community_configurations"
	ProfileMetaConfigurationsCacheKey = "%s_profile_meta_configurations"
)

// Cache TTLs
const (
	ProfileMetaConfigurationsCacheTTL = 6
)
