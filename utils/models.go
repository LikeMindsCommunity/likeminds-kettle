package utils

type AuthenticationResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message"`
	Data         interface{} `json:"data"`
}
