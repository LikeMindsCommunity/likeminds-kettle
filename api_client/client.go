package api_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type APIClient struct {
	CoreServiceBaseURL         string
	SubscriptionServiceBaseURL string
	HTTPClient                 *http.Client
}

type GetRequestOptions struct {
	Url           string
	Params        map[string]string
	CustomHeaders map[string]string
}

type PostRequestOptions struct {
	Url     string
	Body          interface{}
	CustomHeaders map[string]string
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

func NewAPIClient() *APIClient {
	return &APIClient{
		CoreServiceBaseURL:         GetCoreServiceBaseUrl(),
		SubscriptionServiceBaseURL: GetSubscriptionServiceBaseUrl(),
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

func (c *APIClient) sendRequest(req *http.Request, v interface{}) error {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(res.Body)

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	if err != nil {
		return err
	}
	return nil
}

func (c *APIClient) GetRequest(gro *GetRequestOptions) (interface{}, error) {
	req, err := http.NewRequest("GET", gro.Url, nil)
	if err != nil {
		return nil, err
	}

	params := gro.Params
	if params != nil {
		AddParams(req, params)
	}

	headers := gro.CustomHeaders
	if headers != nil {
		AddHeaders(req, headers)
	}

	res := make(map[string]interface{})

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *APIClient) PostRequest(pro *PostRequestOptions) (interface{}, error) {
	jsonData, err := json.Marshal(pro.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", pro.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	headers := pro.CustomHeaders
	if headers != nil {
		AddHeaders(req, headers)
	}

	res := make(map[string]interface{})

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}
