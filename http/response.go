package http

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
)

/*
Autonomy of HTTP Response:

+----------------------------------+
│ <version> <statusCode> <message> │ Status-Line
│ <key: value>                     │
│ <key: value>                     │ Headers
│ <key: value>                     │
│ ...                              │
│                                  │ 1 empty line
│ data. can be:                    │
│ - text/plain                     │
│ - text/html                      │ Body
│ - application/json               │
│ - etc.                           │
+----------------------------------+
*/

type Response struct {
	Version float32
	Status  Status

	Headers Header

	Body []byte

	Request *Request
}

var defaultHttpVersion float32 = 1.1

func newResponse(statusCode int16, headers Header, body []byte, r *Request) (*Response, error) {
	status, err := setStatus(statusCode)
	if err != nil {
		return nil, err
	}

	// default headers for all responses
	headers.Set("Server", "HTTP-Server-by-hamidgh01")
	headers.Set("Date", time.Now().Format(time.RFC1123Z))

	return &Response{
		Version: defaultHttpVersion,
		Status:  status,
		Headers: headers,
		Body:    body,
		Request: r,
	}, nil
}

func JSONResponse(statusCode int16, data any, r *Request) (*Response, error) {
	if data == nil {
		return newResponse(statusCode, make(Header), nil, r)
	}

	bodyBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode data to json. reason: %w", err)
	}

	headers := make(Header)
	contentLength := len(bodyBytes)
	headers.Set("Content-Length", strconv.Itoa(contentLength))
	headers.Set("Content-Type", "application/json") // charset=utf-8

	return newResponse(statusCode, headers, bodyBytes, r)
}

func XMLResponse(statusCode int16, data any, r *Request) (*Response, error) {
	if data == nil {
		return newResponse(statusCode, make(Header), nil, r)
	}

	bodyBytes, err := xml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode data to xml. reason: %w", err)
	}

	headers := make(Header)
	contentLength := len(bodyBytes)
	headers.Set("Content-Length", strconv.Itoa(contentLength))
	headers.Set("Content-Type", "application/xml") // charset=utf-8

	return newResponse(statusCode, headers, bodyBytes, r)
}

func HTMLResponse(statusCode int16, data string, r *Request) (*Response, error) {
	bodyBytes := []byte(data)

	headers := make(Header)
	contentLength := len(bodyBytes)
	headers.Set("Content-Length", strconv.Itoa(contentLength))
	headers.Set("Content-Type", "text/html") // charset=utf-8

	return newResponse(statusCode, headers, bodyBytes, r)
}

func StringResponse(statusCode int16, data string, r *Request) (*Response, error) {
	bodyBytes := []byte(data)

	headers := make(Header)
	contentLength := len(bodyBytes)
	headers.Set("Content-Length", strconv.Itoa(contentLength))
	headers.Set("Content-Type", "text/plain") // charset=utf-8

	return newResponse(statusCode, headers, bodyBytes, r)
}
