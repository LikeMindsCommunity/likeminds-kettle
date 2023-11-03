package utils

const HeadersMemberId = "x-member-id"
const HeadersVersionCode = "x-version-code"
const HeadersPlatformCode = "x-platform-code"
const HeadersPlatformType = "x-platform-type"
const HeadersSdkSource = "x-sdk-source"
const HeadersDeviceId = "x-device-id"
const HeadersApiKey = "x-api-key"
const HeadersAcceptVersion = "x-accept-version"

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
	PlatformDashboard PlatformType = "dashboard"
)

const (
	PlatformAndroid     string = "an"
	PlatformWeb         string = "web"
	PlatformIoS         string = "ios"
	PlatformFlutter     string = "fl"
	PlatformReact       string = "rt"
	PlatformReactNative string = "rn"
)

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
