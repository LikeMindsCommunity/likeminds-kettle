package api_client

import (
	"encoding/json"
)

type APIClientResponse struct {
	Success      bool                   `json:"success"`
	ErrorMessage string                 `json:"error_message"`
	Response     map[string]interface{} `json:"-"`
}

func unmarshalAPIClientResponse(resp string) APIClientResponse {
	apiCR := APIClientResponse{}
	if err := json.Unmarshal([]byte(resp), &apiCR); err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(resp), &apiCR.Response); err != nil {
		panic(err)
	}
	delete(apiCR.Response, "success")
	delete(apiCR.Response, "error_message")

	return apiCR
}
