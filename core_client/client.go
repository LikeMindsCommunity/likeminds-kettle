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
	CoreServiceBaseURL         string
	SubscriptionServiceBaseURL string
	HTTPClient                 *http.Client
}

type GetRequestOptions struct {
	Url           string
	Params        map[string]string
	CustomHeaders map[string]string
	Header        http.Header
}

type PostRequestOptions struct {
	Url     string
	Body    interface{}
	Headers map[string]string
	Header  http.Header
}

func GetCoreServiceBaseUrl() string {
	CoreServiceBaseURL := os.Getenv("CORE_BASE_URL")

	if len(CoreServiceBaseURL) == 0 {
		CoreServiceBaseURL = "https://beta.likeminds.community"
	}

	return CoreServiceBaseURL
}

func GetSubscriptionServiceBaseUrl() string {
	SubscriptionServiceBaseURL := os.Getenv("CORE_BASE_URL")

	if len(SubscriptionServiceBaseURL) == 0 {
		SubscriptionServiceBaseURL = "https://betasubscription.likeminds.community"
	}

	return SubscriptionServiceBaseURL
}

func NewClient() *Client {
	return &Client{
		CoreServiceBaseURL:         GetCoreServiceBaseUrl(),
		SubscriptionServiceBaseURL: GetSubscriptionServiceBaseUrl(),
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func AddHeaders(req *http.Request, headers map[string]string, header http.Header) {
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	if header != nil {
		req.Header = header
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

	headers := options.CustomHeaders
	if headers != nil {
		AddHeaders(req, headers, options.Header)
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
		AddHeaders(req, headers, options.Header)
	}

	res := make(map[string]interface{})

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}
