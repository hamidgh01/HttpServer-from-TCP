package http

/*
Autonomy of HTTP Request:

+----------------------------+
│ <method> <path> <version>  │ Request-Line
│ <key: value>               │
│ <key: value>               │ Headers
│ <key: value>               │
│ ...                        │
│                            │ 1 empty line
│ data. can be:              │
│ - text/plain               │
│ - text/html                │ Body
│ - application/json         │
│ - etc.                     │
+----------------------------+
*/

type Request struct {
	Method  string
	Path    string
	Version float64
	Headers Header
	Body    []byte
}
