package http

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
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

func newResponse(
	statusCode int16,
	body []byte,
	contentType string, // would be "" if body == nil
	r *Request,
) (*Response, error) {

	status, err := setStatus(statusCode)
	if err != nil {
		return nil, err
	}

	headers := make(Header)

	// default headers for all responses
	headers.Set("Server", "HTTP-Server-by-hamidgh01")
	headers.Set("Date", time.Now().Format(time.RFC1123Z))

	if body == nil {
		return &Response{
			Version: defaultHttpVersion,
			Status:  status,
			Headers: headers,
			Body:    nil,
			Request: r,
		}, nil
	}

	headers.Set("Content-Type", contentType) // + charset=utf-8
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
		return newResponse(statusCode, nil, "", r)
	}

	bodyBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode data to json. reason: %w", err)
	}

	return newResponse(statusCode, bodyBytes, "application/json", r)
}

func XMLResponse(statusCode int16, data any, r *Request) (*Response, error) {
	if data == nil {
		return newResponse(statusCode, nil, "", r)
	}

	bodyBytes, err := xml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode data to xml. reason: %w", err)
	}

	return newResponse(statusCode, bodyBytes, "application/xml", r)
}

func HTMLResponse(statusCode int16, data string, r *Request) (*Response, error) {
	bodyBytes := []byte(data)
	return newResponse(statusCode, bodyBytes, "text/html", r)
}

func StringResponse(statusCode int16, data string, r *Request) (*Response, error) {
	bodyBytes := []byte(data + "\r\n")
	return newResponse(statusCode, bodyBytes, "text/plain", r)
}
