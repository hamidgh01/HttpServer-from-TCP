package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/hamidgh01/HttpServer-from-TCP/http"
)

func parseRequest(conn net.Conn) (*http.Request, error) {

	bufferedReader := bufio.NewReader(conn)

	// 1. read and parse Request Line
	method, path, version, err := parseRequestLine(bufferedReader)
	if err != nil {
		return nil, err
	}

	// 2. read and parse Headers
	headers, contentLength, err := parseHeaders(bufferedReader)
	if err != nil {
		return nil, err
	}

	// 3. read body data
	body, err := readBody(bufferedReader, contentLength)
	if err != nil {
		return nil, err
	}

	return &http.Request{
		Method:  method,
		Path:    path,
		Version: version,

		Headers: headers,

		ContentLength: contentLength,
		Body:          body,
	}, nil
}

func parseRequestLine(r *bufio.Reader) (method string, path string, version float64, err error) {

	requestLine, err := r.ReadString('\n')
	if err != nil {
		return "", "", 0, err
	}

	result := strings.Split(strings.TrimSpace(requestLine), " ")

	if len(result) != 3 {
		return "", "", 0, fmt.Errorf("malformed request line: '%s'", requestLine)
	}

	method, path, protocol := result[0], result[1], result[2]

	versionStr := strings.TrimPrefix(strings.ToUpper(protocol), "HTTP/")
	version, err = strconv.ParseFloat(versionStr, 64)

	return
}

func parseHeaders(r *bufio.Reader) (headers http.Header, contentLength int, err error) {
	headers = make(http.Header)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, 0, err
		}

		// empty line -> end of headers
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		if !strings.Contains(line, ":") {
			continue // skipping malformed header lines
		}

		result := strings.SplitN(line, ":", 2)
		key, valuesStr := strings.TrimSpace(result[0]), result[1]

		var values []string
		for value := range strings.SplitSeq(valuesStr, ",") {
			values = append(values, strings.TrimSpace(value))
		}

		headers.Set(key, values...) // overwrite if redeclared

		if key == "Content-Length" {
			contentLength, err = strconv.Atoi(values[0])
			if err != nil {
				return nil, 0, err
			} else if contentLength < 0 {
				return nil, 0, fmt.Errorf("negative Content-Length value: %d", contentLength)
			}
		}
	}

	return
}

func readBody(r *bufio.Reader, contentLength int) (body []byte, err error) {
	if contentLength == 0 {
		return nil, nil
	}

	body = make([]byte, 0, contentLength)
	_, err = io.ReadFull(r, body)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
