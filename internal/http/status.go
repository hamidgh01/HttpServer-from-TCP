package http

import (
	"errors"
	"fmt"
)

var (
	ErrStatusNotSupported = errors.New("status not supported.")
)

var supportedStatuses = map[int16]string{
	// 1** : Informational
	100: "Continue",

	// 2** : Successful
	200: "OK",
	201: "Created",
	202: "Accepted",
	204: "No Content",

	// 3** : Redirection
	301: "Moved Permanently",
	304: "Not Modified",
	307: "Temporary Redirect",
	308: "Permanent Redirect",

	// 4** : Client Error
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	408: "Request Timeout",
	409: "Conflict",
	411: "Length Required",
	413: "Request Entity Too Large",
	422: "Unprocessable Entity",
	429: "Too Many Requests",

	// 5** : Server Error
	500: "Internal Server Error",
	502: "Bad Gateway",
}

type Status struct {
	Code    int16
	Message string
}

func setStatus(code int16) (Status, error) {
	message, ok := supportedStatuses[code]
	if !ok {
		return Status{}, ErrStatusNotSupported
	}

	return Status{Code: code, Message: message}, nil
}

func (s Status) String() string {
	return fmt.Sprintf("%d %s", s.Code, s.Message)
}
