package server

import (
	"errors"
	"fmt"
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
		request, err, errPosition := parseRequest(conn)
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

			s.serveWith400BadRequest(conn, errPosition)
			return
		}

		s.logger.Infof(
			"received '%s %s' from '%s'",
			request.Method, request.Path, conn.RemoteAddr().String(),
		)

		// 2. analyze received request and build proper response -> implement later (ToDo)
		response, err := handleRequestAndGetResponse(request)
		if err != nil {
			s.logger.Errorf("failed to build response: %s", err.Error())
			s.serveWith500InternalServerError(conn, request)
			return
		}

		// 3. encode response to raw bytes and send (write to conn)
		responseBytes := ToBytes(response)
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
			request.Method, request.Path, conn.RemoteAddr().String(), response.Status,
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

func (s *Server) serveWith400BadRequest(conn net.Conn, errPosition ErrorPosition) {
	respBody := fmt.Sprintf(
		"400 Bad Request (malformed request format)\r\napproximate malformed position: %s",
		errPosition,
	)
	response, _ := http.StringResponse(400, respBody, nil)

	if _, err := conn.Write(ToBytes(response)); err != nil {
		s.logger.Errorf("failed to send response (400) over connection: %s", err.Error())
		return
	}

	s.logger.Infof(
		"malformed request from '%s' is served with '400 Bad Request'", conn.RemoteAddr().String(),
	)
}

func (s *Server) serveWith500InternalServerError(conn net.Conn, req *http.Request) {
	if _, err := conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n")); err != nil {
		s.logger.Errorf("failed to send response (500) over connection: %s", err.Error())
		return
	}

	s.logger.Infof(
		"'%s %s' from '%s' is served with '500 Internal Server Error'",
		req.Method, req.Path, conn.RemoteAddr().String(),
	)
}
