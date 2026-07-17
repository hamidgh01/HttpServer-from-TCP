package server

import (
	"errors"
	"io"
	"net"
	"os"
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
	for {

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
	}

}
