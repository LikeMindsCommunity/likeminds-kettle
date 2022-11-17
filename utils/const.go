package utils

const HeadersMemberId = "x-member-id"
const HeadersVersionCode = "x-version-code"
const HeadersPlatformCode = "x-platform-code"
const HeadersPlatformType = "x-platform-type"
const HeadersDeviceId = "x-device-id"
const HeadersApiKey = "x-api-key"

const GETMethod = 0
const POSTMethod = 1
const PUTMethod = 2
const DELETEMethod = 3

type ServiceType int

const (
	CoreService ServiceType = iota
	SubscriptionService
)

type RequestType int

const (
	GETRequest RequestType = iota
	POSTRequestRawBody
	POSTRequestFormUrlEncodedBody
	PUTRequest
	DELETERequest
)

type PlatformType string

const (
	PlatformDashboard PlatformType = "dashboard"
)

const (
	PlatformAndroid string = "an"
	PlatformWeb     string = "web"
	PlatformIoS     string = "ios"
)

var MaxVersion int = 9999

var QuestionIdVersions = map[string]int{
	PlatformAndroid: MaxVersion,
	PlatformWeb:     MaxVersion,
	PlatformIoS:     363,
}
