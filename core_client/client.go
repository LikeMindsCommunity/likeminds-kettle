package core_client

import (
	"bytes"
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

type GetRequestOptions struct {
	Url     string
	Params  map[string]string
	Headers map[string]string
}

type PostRequestOptions struct {
	Url     string
	Body    interface{}
	Headers map[string]string
}

func GetBaseUrl() string {
	BaseURL := os.Getenv("CORE_BASE_URL")

	if len(BaseURL) == 0 {
		BaseURL = "https://beta.likeminds.community"
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

func AddParams(req *http.Request, params map[string]string) {
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetRequest(options *GetRequestOptions) (interface{}, error) {
	req, err := http.NewRequest("GET", options.Url, nil)
	if err != nil {
		return nil, err
	}

	params := options.Params
	if params != nil {
		AddParams(req, params)
	}

	headers := options.Headers
	if headers != nil {
		AddHeaders(req, headers)
	}

	res := make(map[string]interface{})

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) PostRequest(options *PostRequestOptions) (interface{}, error) {
	jsonData, err := json.Marshal(options.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", options.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	headers := options.Headers
	if headers != nil {
		AddHeaders(req, headers)
	}

	res := make(map[string]interface{})

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}
