package server

import (
	"bytes"
	"fmt"

	"github.com/hamidgh01/HttpServer-from-TCP/http"
)

const CRLF = "\r\n"

func ToBytes(resp *http.Response) []byte {
	buf := new(bytes.Buffer)
	mergeStatusLine(resp, buf)
	mergeHeaderLines(resp, buf)
	mergeBody(resp, buf)

	return buf.Bytes()
}

func mergeStatusLine(resp *http.Response, buf *bytes.Buffer) {
	fmt.Fprintf(buf, "HTTP/%v %s %s", resp.Version, resp.Status, CRLF)
}

func mergeHeaderLines(resp *http.Response, buf *bytes.Buffer) {
	for key, values := range resp.Headers {
		fmt.Fprintf(buf, "%s:", key)
		for _, v := range values {
			fmt.Fprintf(buf, " %s", v)
		}
		fmt.Fprintf(buf, "%s", CRLF)
	}

	fmt.Fprintf(buf, "%s", CRLF) // an empty-line (marks end of headers, and start of body)
}

func mergeBody(resp *http.Response, buf *bytes.Buffer) {
	if resp.Body == nil {
		return
	}

	buf.Write(resp.Body)
}
