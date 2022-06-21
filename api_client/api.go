package api_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type BodyType int

const (
	BodyTypeRaw BodyType = iota
	BodyTypeFormUrlEncoded
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
	Params        map[string]string
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

func UpdateBody(pro *PostRequestOptions, body_type BodyType) (*http.Request, error) {

	var req *http.Request

	switch body_type {
	case BodyTypeRaw:

		data, err := json.Marshal(pro.Body)

		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest(http.MethodPost, pro.Url, bytes.NewBuffer(data))

		if err != nil {
			return nil, err
		}

	case BodyTypeFormUrlEncoded:

		var err error

		body, _ := json.Marshal(pro.Body)
		payload := convertToFormURLEncoded(&body)

		req, err = http.NewRequest(http.MethodPost, pro.Url, strings.NewReader(payload.Encode()))

		if err != nil {
			return nil, err
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(payload.Encode())))

	}

	return req, nil
}

func convertToFormURLEncoded(body *[]byte) url.Values {
	// datamap | converts incoming request body into a map
	var datamap map[(string)]interface{}
	json.Unmarshal(*body, &datamap)

	// payload | form-url-encode paylaod
	payload := url.Values{}

	// loop over map and fill payload data
	for key, value := range datamap {

		// if value is not type string then convert to string
		switch reflect.TypeOf(value).String() {

		case "int":
			payload.Set(key, strconv.Itoa(value.(int)))

		case "float64":
			var elementString string = strconv.FormatFloat(value.(float64), 'f', -1, 64)
			payload.Set(key, elementString)

		case "bool":
			payload.Set(key, strconv.FormatBool(value.(bool)))

		case "string":
			payload.Set(key, value.(string))

		}
	}

	return payload
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
	req, err := http.NewRequest(http.MethodGet, gro.Url, nil)
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

func (c *APIClient) PostRequest(pro *PostRequestOptions, body_type BodyType) ([]byte, error) {

	req, err := UpdateBody(pro, body_type)
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

	req, err := http.NewRequest(http.MethodPut, pro.Url, bytes.NewBuffer(jsonData))
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

func (c *APIClient) DeleteRequest(pro *PostRequestOptions) ([]byte, error) {
	jsonData, err := json.Marshal(pro.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, pro.Url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	params := pro.Params
	if params != nil {
		AddParams(req, params)
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
