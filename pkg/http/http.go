package http

import (
	"bytes"
	"changer/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"time"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"

	STATUS_OK         = 200
	STATUS_CREATED    = 201
	STATUS_NO_CONTENT = 204
)

type Request struct {
	Client  *http.Client
	Url     string
	Method  string
	Body    interface{}
	Headers map[string]string
	Timeout time.Duration
}

func NewRequest(method string, url string, body interface{}) *Request {
	return &Request{
		Method:  method,
		Client:  &http.Client{},
		Headers: make(map[string]string),
		Timeout: 30 * time.Second,
		Url:     url,
		Body:    body,
	}
}
func (r *Request) SetAuthorizationHeaders(token string) {
	r.SetHeaders("Authorization", "Bearer "+token)

}
func (r *Request) SetHeaders(key string, value string) {
	r.Headers[key] = value
}

func (r *Request) SetBody(body interface{}) error {
	// Check if method is POST or PUT
	if r.Method != POST && r.Method != PUT && r.Method != PATCH {
		return errors.New("body can not be set for this method: " + r.Method)
	}
	// Set Body
	if body != nil {

		b, err := json.Marshal(&body)
		if err != nil {
			return err
		}
		r.Body = b
	}
	return nil

}
func (r *Request) SetTimeout(timeout time.Duration) {
	r.Timeout = timeout

}

// Set the body in the request
func (r *Request) getBodyReader() io.Reader {
	if r.Body == nil {
		return nil
	}
	b, err := json.Marshal(r.Body)
	if err != nil {
		return nil
	}
	return bytes.NewReader(b)
}

func (r *Request) buildHeaders(req *http.Request) {
	// Render all the headers and sets them to the requestS
	req.Header.Set("Content-Type", "application/json")

	if r.Headers != nil {
		for k, v := range r.Headers {
			req.Header.Add(k, v)
		}
	}

}

// Check if the status code is 2xx
func (r *Request) checkStatusCode(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return errors.New("request failed, status code: " + fmt.Sprint(resp.StatusCode))
}

// Unmarshal the response body into a BaseResponse struct
func (r *Request) unmarshalBody(resp *http.Response) (*models.BaseResponse, error) {
	// Close the response body
	defer resp.Body.Close()

	response := &models.BaseResponse{}

	// Read the response body
	bpdy, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	fmt.Println(string(bpdy))
	// Unmarshal the response body into a Response struct
	if err := json.Unmarshal(bpdy, &response); err != nil {
		return response, err
	}

	return response, nil
}

// Build the request and send it
func (r *Request) do() (*http.Response, error) {
	r.Client.Timeout = r.Timeout

	// Create a new request using http
	req, err := http.NewRequest(r.Method, r.Url, r.getBodyReader())

	if err != nil {
		return nil, err
	}
	// Set headers
	r.buildHeaders(req)

	// Send the request and check for errors
	return r.Client.Do(req)

}

func (r *Request) Execute() (*models.BaseResponse, error) {

	// Parse the response body into a Response struct
	response := &models.BaseResponse{}

	// Execute the request
	resp, err := r.do()
	// If the request fails after the retry limit, return the error
	if err != nil {
		return response, err
	}

	// Unmarshal the response body
	response, err = r.unmarshalBody(resp)
	if err != nil {
		return response, err
	}
	// Check if the status code is not 2xx
	if err = r.checkStatusCode(resp); err != nil {
		return response, fmt.Errorf("request failed, status code: %s, message: %+v", fmt.Sprint(resp.StatusCode), response)
	}
	// Check if the request was successful

	return response, nil

}
