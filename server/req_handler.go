package server

import (
	"slices"
	"strconv"

	"github.com/hamidgh01/HttpServer-from-TCP/http"
)

func handleRequestAndGetResponse(req *http.Request) (*http.Response, error) {

	// 1. send request to backend Framework and receive response
	response := backendFrameworkSimulation(req)

	if response.Body == nil {
		return response, nil
	}

	// 2. compress response body if request accepts compression (has `Accept-Encoding` header)
	response = handleCompression(req, response)

	return response, nil
}

func backendFrameworkSimulation(req *http.Request) *http.Response {
	type requestComponents struct {
		Method        string      `json:"method"`
		Path          string      `json:"path"`
		Headers       http.Header `json:"headers"`
		ContentLength int         `json:"content_length"`
	}

	defaultBody := struct {
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

	response, _ := http.JSONResponse(200, defaultBody, req)
	return response
}

func handleCompression(req *http.Request, resp *http.Response) *http.Response {

	var contentLength int
	acceptEncodings := req.Headers.Get("Accept-Encoding")

	if len(acceptEncodings) == 0 {
		contentLength = len(resp.Body) // Note: response has a non-nil body
		resp.Headers.Set("Content-Length", strconv.Itoa(contentLength))
		return resp
	}

	var compressionAlgorithm CompressionAlg
	switch {
	case slices.Contains(acceptEncodings, "gzip"): // first priority
		compressionAlgorithm = Gzip
	case slices.Contains(acceptEncodings, "deflate"): // second priority
		compressionAlgorithm = Deflate
	default: // others (like `br` and `zstd`) are not supported, so:
		contentLength = len(resp.Body)
		resp.Headers.Set("Content-Length", strconv.Itoa(contentLength))
		return resp
	}

	compressedBody, err := Compress(resp.Body, compressionAlgorithm)
	if err != nil {
		// log.Error("failed to compress body. reason: %s", err.Error())

		contentLength = len(resp.Body)
		resp.Headers.Set("Content-Length", strconv.Itoa(contentLength))
		return resp
	}

	resp.Body = compressedBody
	contentLength = len(compressedBody)
	resp.Headers.Set("Content-Length", strconv.Itoa(contentLength))
	resp.Headers.Set("Content-Encoding", string(compressionAlgorithm))

	return resp
}
