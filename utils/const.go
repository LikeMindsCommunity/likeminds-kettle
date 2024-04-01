package utils

const HeadersMemberId = "x-member-id"
const HeadersVersionCode = "x-version-code"
const HeadersPlatformCode = "x-platform-code"
const HeadersPlatformType = "x-platform-type"
const HeadersSdkSource = "x-sdk-source"
const HeadersDeviceId = "x-device-id"
const HeadersApiKey = "x-api-key"
const HeadersAcceptVersion = "x-accept-version"
const HeadersApiVersion = "x-api-version"
const HeaderMemberRole = "x-member-role"

const GETMethod = 0
const POSTMethod = 1
const PUTMethod = 2
const DELETEMethod = 3
const PatchMethod = 4

type ServiceType int

const (
	CoreService ServiceType = iota
	SubscriptionService
	SwarmService
)

type RequestType int

const (
	GETRequest RequestType = iota
	POSTRequestRawBody
	POSTRequestFormUrlEncodedBody
	PUTRequest
	DELETERequest
	PATCHRequest
)

type PlatformType string

const (
	PlatformDashboard      PlatformType = "dashboard"
	PlatformCaravanService PlatformType = "caravan-service"
	PlatformSwarmService   PlatformType = "swarm-service"
	PlatformKettleService  PlatformType = "kettle-service"
)

const (
	PlatformAndroid     string = "an"
	PlatformWeb         string = "web"
	PlatformIoS         string = "ios"
	PlatformFlutter     string = "fl"
	PlatformReact       string = "rt"
	PlatformReactNative string = "rn"
)

// Valid Platform Codes
var ValidFrontendPlatformCodes = []string{
	PlatformAndroid,
	PlatformWeb,
	PlatformIoS,
	PlatformFlutter,
	PlatformReact,
	PlatformReactNative,
}

var MaxVersion int = 9999

var QuestionIdVersions = map[string]int{
	PlatformAndroid:     202,
	PlatformWeb:         MaxVersion,
	PlatformIoS:         363,
	PlatformFlutter:     MaxVersion,
	PlatformReact:       MaxVersion,
	PlatformReactNative: MaxVersion,
}

const (
	ApiRevampV1 string = "v1"
)

// Context Params | These are the params that are passed in the context
const (
	ParamCommunityId = "community_id"
)

// Community Configurations | community configuration types
const (
	CommunityConfigurationMediaLimits     = "media_limits"
	CommunityConfigurationFeedMetadata    = "feed_metadata"
	CommunityConfigurationProfileMetadata = "profile_metadata"
)

// profile metadata configurations s| profile metadata types
const (
	ConfigurationsProfileMetaWidgetsEnabled = "widgets_enabled"
)

// Params
const ParamConfigurationTypes = "configuration_types"
const ParamWidgetIds = "widget_ids"
const ParamMemberIds = "member_ids"
const ParamParentIds = "parent_ids"
const ParamUUIDs = "uuids"
const ParamUUID = "uuid"
const ParamPageSize = "page_size"
const ParamCommunityID = "community_id"
const ParamTierType = "tier_type"

// Endpoints
const FetchCommunityConfigurationsEndpoint = "/api/community/configurations"
const WidgetEndPoint = "/widget"
const FetchCommunitySettingsEndpoint = "/api/community/fetch_community_settings"

const FetchMembersMetaEndPoint = "/api/community/fetch_members_meta"
const FetchTopicsEndpoint = "/topic"
const FetchUserTopicsEndpoint = "/user/topics"
const SDKAuthenticateEndPoint = "/api/sdk/authenticate"

// Skulk Endpoints
const TierEndpoint = "api/subscription/plan/tiers"
const BillingPlanEnpoint = "api/subscription/plan/billing"

// Community settings type
const (
	FeedRepostCommunitySettingType  = "feed_repost"
	UserTopicsConnectionSettingType = "user_topics_connection"
)

// Member Roles
const (
	GuestRole string = "GUEST"
)
