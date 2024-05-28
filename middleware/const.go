package middleware

const RPM = 1 // RPM

// StrictTransportSecurityHeaderKey | strict transport security (TLS) header key
const StrictTransportSecurityHeaderKey = "Strict-Transport-Security"

// StrictTransportSecurityHeaderValue | strict transport security (TLS) header value
const StrictTransportSecurityHeaderValue = "max-age=31536000; includeSubDomains"

// ContentTypeOptionsHeaderKey | api content type options header key
const ContentTypeOptionsHeaderKey = "X-Content-Type-Options"

// ContentTypeOptionsHeaderValue | api content type options header value
const ContentTypeOptionsHeaderValue = "nosniff"

// CacheControlHeaderKey | application cache control header key
const CacheControlHeaderKey = "Cache-Control"

// CacheControlHeaderValue | application cache control header value
const CacheControlHeaderValue = "no-cache; no-store; must-revalidate"
