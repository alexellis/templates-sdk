// Copyright (c) Alex Ellis 2018. All rights reserved.
// Copyright (c) OpenFaaS Author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package handler

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
)

type Response interface {
	GetHeader() http.Header

	GetBody() []byte

	GetStatusCode() int
}

type Request interface {
	GetHeader() http.Header

	GetBody() []byte

	GetStatusCode() int

	GetQueryString() string
	GetMethod() string
	GetHost() string
	Context() context.Context
}

type FunctionResponse struct {
	// Body the body will be written back
	body []byte

	// StatusCode needs to be populated with value such as http.StatusOK
	statusCode int

	// Header is optional and contains any additional headers the function response should set
	header http.Header
}

func NewFunctionResponse(body []byte, statusCode int, header http.Header) *FunctionResponse {
	return &FunctionResponse{
		body:       body,
		statusCode: statusCode,
		header:     header,
	}
}

func (r *FunctionResponse) GetHeader() http.Header {
	return r.header
}

func (r *FunctionResponse) GetBody() []byte {
	return r.body
}

func (r *FunctionResponse) GetStatusCode() int {
	return r.statusCode
}

type FunctionRequest struct {
	body        []byte
	header      http.Header
	queryString string
	method      string
	host        string
	ctx         context.Context
}

func NewFunctionRequest(r *http.Request) *FunctionRequest {

	var body []byte
	if r.Body != nil {
		defer r.Body.Close()
		var err error
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %s", err)
		}
	}

	return &FunctionRequest{
		body:        body,
		header:      r.Header,
		queryString: r.URL.RawQuery,
		method:      r.Method,
		host:        r.Host,
		ctx:         r.Context(),
	}
}

func (req *FunctionRequest) GetHeader() http.Header {
	return req.header
}

func (req *FunctionRequest) GetBody() []byte {
	return req.body
}

func (req *FunctionRequest) GetQueryString() string {
	return req.queryString
}

func (req *FunctionRequest) GetMethod() string {
	return req.method
}

func (req *FunctionRequest) GetHost() string {
	return req.host
}

// Context is set for optional cancellation inflight requests.
func (req *FunctionRequest) GetContext() context.Context {
	return req.ctx
}

// WithContext overides the context for the Request struct
func (req *FunctionRequest) WithContext(ctx context.Context) {
	// AE: Not keen on panic mid-flow in user-code, however stdlib also appears to do
	// this. https://golang.org/src/net/http/request.go
	// This is not setting a precedent for broader use of "panic" to handle errors.
	if ctx == nil {
		panic("nil context")
	}
	req.ctx = ctx
}

// FunctionHandler used for a serverless Go method invocation
type FunctionHandler interface {
	Handle(req FunctionRequest) (FunctionResponse, error)
}

func init() {

}
