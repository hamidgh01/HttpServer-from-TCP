package http

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
