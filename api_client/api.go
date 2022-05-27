package api_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	CustomHeaders map[string]interface{}
}

type PostRequestOptions struct {
	Url           string
	Body          interface{}
	CustomHeaders map[string]interface{}
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

func AddHeaders(req *http.Request, headers map[string]interface{}) {
	for k, v := range headers {
		req.Header.Add(k, v.(string))
	}
}

func AddParams(req *http.Request, params map[string]string) {
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
}

func (c *APIClient) sendRequest(req *http.Request) ([]byte, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("unknown error, status code: %d", resp.StatusCode)
	}
	//Defer close error
	if err != nil {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBytes, nil
}

func (c *APIClient) GetRequest(gro *GetRequestOptions) ([]byte, error) {
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

	respBytes, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}

func (c *APIClient) PostRequest(pro *PostRequestOptions) ([]byte, error) {
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

	respBytes, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}

func (c *APIClient) PutRequest(pro *PostRequestOptions) ([]byte, error) {
	jsonData, err := json.Marshal(pro.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", pro.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	headers := pro.CustomHeaders
	if headers != nil {
		AddHeaders(req, headers)
	}

	respBytes, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}
