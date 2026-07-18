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
	fmt.Fprintf(buf, "HTTP/%v %s%s", resp.Version, resp.Status, CRLF)
}

func mergeHeaderLines(resp *http.Response, buf *bytes.Buffer) {
	for key, values := range resp.Headers {
		buf.WriteString(key)
		buf.WriteString(": ")
		for idx, val := range values {
			if idx == 0 {
				buf.WriteString(val)
			} else {
				buf.WriteString(", ")
				buf.WriteString(val)
			}
		}
		buf.WriteString(CRLF)
	}

	buf.WriteString(CRLF) // an empty-line (marks end of headers, and start of body)
}

func mergeBody(resp *http.Response, buf *bytes.Buffer) {
	if resp.Body == nil {
		return
	}

	buf.Write(resp.Body)
}
