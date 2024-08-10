package token

// default auth token expiry in minutes
const (
	PROD_AUTH_TOKEN_EXPIRY = 15
	BETA_AUTH_TOKEN_EXPIRY = 60

	REFRESH_TOKEN_EXPIRY = 744 // 31 days
	DEFAULT_TOKEN_EXPIRY = -1
)

// Auth token data params
const (
	TokemAccessUUID   = "access_uuid"
	TokenAPIKey       = "api_key"
	TokenUserUniqueId = "user_unique_id"
	TokenIsGuest      = "is_guest"
	TokenExp          = "exp"
	TokenDeviceID     = "device_id"
)
