package api_client

import (
	"encoding/json"
)

// APIClientResponse Used only for internal API calls
type APIClientResponse struct {
	Success      bool                   `json:"success"`
	ErrorMessage string                 `json:"error_message"`
	ErrorMeta    map[string]interface{} `json:"error_meta,omitempty"`
	Response     map[string]interface{} `json:"-"`
}

// UnmarshalAPIClientResponse used to unmarshal APIClientResponse i.e internal API call response
func UnmarshalAPIClientResponse(resp []byte, apiCR *APIClientResponse) error {
	if err := json.Unmarshal(resp, &apiCR); err != nil {
		return err
	}

	if err := json.Unmarshal(resp, &apiCR.Response); err != nil {
		return err
	}
	delete(apiCR.Response, "success")
	delete(apiCR.Response, "error_message")
	return nil
}
