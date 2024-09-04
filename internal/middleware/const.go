package middleware

const RPM = 1 // RPM

// Security Headers
const (

	// StrictTransportSecurityHeaderKey | strict transport security (TLS) header key
	StrictTransportSecurityHeaderKey = "Strict-Transport-Security"

	// StrictTransportSecurityHeaderValue | strict transport security (TLS) header value, max-age (in hours)
	StrictTransportSecurityHeaderValue = "max-age=720; includeSubDomains"

	// ContentTypeOptionsHeaderKey | api content type options header key
	ContentTypeOptionsHeaderKey = "X-Content-Type-Options"

	// ContentTypeOptionsHeaderValue | api content type options header value
	ContentTypeOptionsHeaderValue = "nosniff"

	// CacheControlHeaderKey | application cache control header key
	CacheControlHeaderKey = "Cache-Control"

	// CacheControlHeaderValue | application cache control header value
	CacheControlHeaderValue = "no-cache; no-store; must-revalidate"
)

// Context Headers
const (
	ContextApiKeyHeader = "X-Api-Key"
	ContextPlatformTypeHeader = "X-Platform-Type"
)
