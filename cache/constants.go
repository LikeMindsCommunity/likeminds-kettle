package cache

// Params
const (
	ParamRedisClient = "redis_client"
)

// Cache Keys
const (
	CommunityConfigurationsCacheKey   = "%d_community_configurations"
	ProfileMetaConfigurationsCacheKey = "%d_profile_meta_configurations"
)

// Cache TTLs
const (
	ProfileMetaConfigurationsCacheTTL = 6
)
