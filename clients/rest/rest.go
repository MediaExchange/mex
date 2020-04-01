/*
   Copyright 2019 Paul Howes

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/MediaExchange/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type RestRequest struct {
	*http.Request
	query url.Values
	restError error
	replyBody interface{}
}

// Client returns a new REST client that can be used to call a RESTful web service.
func NewRequest() *RestRequest {
	req := new(RestRequest)
	req.Request = new(http.Request)
	req.Proto = "HTTP/1.1"
	req.ProtoMajor = 1
	req.ProtoMinor = 1
	req.Header = make(http.Header)
	req.Header.Set("Accept", "application/json")
	req.query = make(url.Values)
	return req
}

// AddQuery adds a key/value pair to the request's query string.
func (req *RestRequest) AddQuery(key string, value string) *RestRequest {
	// Propagate previous errors.
	if req.restError != nil {
		return req
	}

	if len(key) == 0 {
		req.restError = errors.New("rest: query key must be provided")
		return req
	}

	if len(value) == 0 {
		req.restError = errors.New("rest: query value must be provided")
		return req
	}

	req.query.Add(key, value)
	return req
}

// SetBearerAuth sets the Authorization header with a bearer token.
func (req *RestRequest) SetBearerAuth(t string) *RestRequest {
	// Propagate previous errors.
	if req.restError != nil {
		return req
	}

	if len(t) == 0 {
		req.restError = errors.New("rest: token must be provided")
		return req
	}

	req.Header.Set("Authorization", "Bearer " + t)
	return req
}

// Body sets the request body to a serialized JSON object.
func (req *RestRequest) SetBody(b interface{}) *RestRequest {
	// Propagate previous errors.
	if req.restError != nil {
		return req
	}

	if b == nil {
		req.restError = errors.New("rest: body must be provided")
		return req
	}

	j, err := json.Marshal(b)
	if err != nil {
		req.restError = err
		return req
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(j))
	req.ContentLength = int64(len(j))
	req.GetBody = func() (io.ReadCloser, error) {
		return req.Body, nil
	}

	req.Header.Set("Content-Type", "application/json")
	return req
}

func (req *RestRequest) SetReplyBody(t interface{}) *RestRequest {
	if t == nil {
		req.restError = errors.New("rest: reply body must be provided.")
	}
	req.replyBody = t
	return req
}

func (req *RestRequest) Get(url string) (*http.Response, error) {
	req.Method = "GET"
	return req.do(url)
}

// Post sends a POST request to the remote server.
func (req *RestRequest) Post(url string) (*http.Response, error) {
	req.Method = "POST"
	return req.do(url)
}

func (req *RestRequest) do(rawurl string) (*http.Response, error) {
	// Propagate earlier errors.
	if req.restError != nil {
		return nil, req.restError
	}

	// Parse the URL string to properly configure the request.
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	req.Host = u.Host
	req.URL = u
	req.URL.RawQuery = req.query.Encode()

	// Run the request using the built-in http library. This handles all
	// buffering and 3xx redirect responses.
	log.Info("Calling REST service", log.String("url", req.URL.String()))
	res, err :=  http.DefaultClient.Do(req.Request)

	if err != nil {
		log.Error("rest.do: Unexpected error", log.Err(err))
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return res, errors.New(res.Status)
	}

	// Attempt to deserialize the body if no error occurred.
	if req.replyBody != nil {
		err := json.NewDecoder(res.Body).Decode(req.replyBody)
		return res, err
	}

	return res, err
}
