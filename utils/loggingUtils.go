package utils

func GeneratePlatformString(platform_code string) string{
	switch platform_code {
	case PlatformAndroid:
		return "android"
	case PlatformWeb:
		return "web"
	case PlatformIoS:
		return "ios"
	case PlatformFlutter:
		return "flutter"
	case PlatformReact:
		return "react"
	case PlatformReactNative:
		return "react-native"
	default:
		return "unknown"
	}
}
