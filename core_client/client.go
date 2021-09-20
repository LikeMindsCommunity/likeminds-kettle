package core_client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func GetBaseUrl() string {
	BaseURL := os.Getenv("CORE_BASE_URL")

	if len(BaseURL) == 0 {
		BaseURL = "https://beta.likeminds.community/"
	}

	return BaseURL
}

func NewClient() *Client {
	return &Client{
		BaseURL: GetBaseUrl(),
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func AddHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Add(k, v)
	}
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	response := v
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return err
	}

	return nil
}
