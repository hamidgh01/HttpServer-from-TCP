package server

import (
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/hamidgh01/HttpServer-from-TCP/http"
)

func (s *Server) handleConnection(conn net.Conn) {

	defer func() {
		conn.Close()
		s.wg.Done()
		s.logger.Infof("connection from '%s' closed.", conn.RemoteAddr().String())
	}()

	conn.SetDeadline(time.Now().Add(s.connTimeout))
	reqCounter := 0
	for reqCounter <= s.maxRequestsPerConn {

		// 1. parse incoming bytes to http request
		request, err := parseRequest(conn)
		if err != nil {
			if errors.Is(err, io.EOF) ||
				errors.Is(err, io.ErrUnexpectedEOF) ||
				errors.Is(err, os.ErrDeadlineExceeded) {
				s.logger.Debugf("close connection due to: %s", err.Error())
				return
			}

			s.logger.Infof(
				"received malformed request from '%s'. error: %s",
				conn.RemoteAddr().String(), err.Error(),
			)

			// ToDo: send a HTTP Response with status = '400 Bad Request'
			// (define a defaultBadRequestResponse of type Response and send it)
			// s.logger.Info("malformed request from '%s' is served with '400 Bad Request'")

			return
		}

		s.logger.Infof(
			"received '%s %s' from '%s'",
			request.Method, request.Path, conn.RemoteAddr().String(),
		)

		// 2. analyze received request and build proper response -> implement later (ToDo)
		defaultBody := struct {
			Message string `json:"message"`
		}{Message: "Hello from HTTP Server..."}

		defaultResponse, err := http.JSONResponse(200, defaultBody, request)
		if err != nil {
			s.logger.Errorf("failed to build response: %s", err.Error())
			// ToDo: send a HTTP Response with status = '500 Internal Server Error'
			// (define a ServerErrorResponse of type Response and send it)
			return
		}

		// 3. encode response to raw bytes and send (write to conn)
		responseBytes := ToBytes(defaultResponse)
		s.logger.Debugf(
			"response for '%s %s' from '%s':\n%s\n",
			request.Method, request.Path, conn.RemoteAddr().String(), responseBytes,
		)

		if _, err := conn.Write(responseBytes); err != nil {
			s.logger.Errorf("failed to send response over connection: %s", err.Error())
			return
		}

		s.logger.Infof(
			"'%s %s' from '%s' is served with '%s'",
			request.Method, request.Path, conn.RemoteAddr().String(), defaultResponse.Status,
		)

		// 4. decide to keep alive or close connection (base on the request header and version)
		if !shouldKeepAlive(request) {
			return
		}

		// if keep alive: refresh deadline & increase reqCounter by 1
		conn.SetDeadline(time.Now().Add(s.connTimeout))
		reqCounter++
	}

}

func shouldKeepAlive(req *http.Request) bool {
	ConnectionHeader := req.Headers.Get("Connection")
	if req.Version >= 1.1 {
		// for HTTP version 1.1 and higher: default is `keep-alive` unless "Connection: close"
		switch {
		case len(ConnectionHeader) == 0:
		case strings.ToLower(ConnectionHeader[0]) == "close":
			return false
		}
		return true
	}

	// for HTTP version 1.0 and lower: default is `close` unless "Connection: keep-alive"
	switch {
	case len(ConnectionHeader) == 0:
	case strings.ToLower(ConnectionHeader[0]) == "keep-alive":
		return true
	}
	return false
}
