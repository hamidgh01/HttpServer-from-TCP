package server

import "github.com/hamidgh01/HttpServer-from-TCP/http"

func handleRequestAndGetResponse(req *http.Request) (*http.Response, error) {

	type requestComponents struct {
		Method        string      `json:"method"`
		Path          string      `json:"path"`
		Headers       http.Header `json:"headers"`
		ContentLength int         `json:"content_length"`
	}

	body := struct {
		Message string            `json:"message"`
		Request requestComponents `json:"your_request_components"`
	}{
		Message: "your request is ready to be processed by a backend framework...",
		Request: requestComponents{
			Method:        req.Method,
			Path:          req.Path,
			Headers:       req.Headers,
			ContentLength: req.ContentLength,
		},
	}

	return http.JSONResponse(200, body, req)
}
