package constants

const (
	HeaderAuthorization = "Authorization"
	ParamAccessToken    = "access_token"
	ParamRefreshToken   = "refresh_token"
)

const (
	ParamLTM             = "ltm"
	ParamVTM             = "vtm"
	ParamOTM             = "otm"
	ParamRTM             = "rtm"
	ErrorInvalidVTM      = "Invalid VTM."
	ErrorInvalidOTM      = "Invalid OTM."
	ErrorInvalidLTM      = "Invalid LTM."
	ErrorInvalidRTM      = "Invalid RTM."
	ErrorInvalidLTMorVTM = "Invalid LTM or VTM."
)

type OnboardingTokenMeta struct {
	AccessUuid         string
	AccessTokenExpires int64
	AccessToken        string
	ApiKey             string
}

type LoginTokenMeta struct {
	AccessUuid         string
	AccessToken        string
	AccessTokenExpires int64
	UserUniqueID       string
	ApiKey             string
	IsGuest            bool
}

type RefreshTokenMeta struct {
	RefreshUuid         string
	RefreshToken        string
	RefreshTokenExpires int64
	UserUniqueID        string
	ApiKey              string
	IsGuest             bool
}

type VerifyTokenMeta struct {
	AccessUuid         string
	AccessTokenExpires int64
	AccessToken        string
	ApiKey             string
	EmailID            string
	MobileNo           string
	CountryCode        string
}
