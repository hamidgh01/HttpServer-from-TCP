package config

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
)

// Default configuration settings
const (
	// default server configurations
	DEFAULT_SERVER_HOST            string = "127.0.0.1"
	DEFAULT_SERVER_PORT            int    = 8000
	DEFAULT_TCP_CONNECTION_TIMEOUT int    = 5 // seconds
	DEFAULT_MAX_REQUESTS_PER_CONN  int    = 50

	// default logger configurations
	DEFAULT_LOG_LEVEL       string = "info"
	DEFAULT_LOG_OUTPUT_FILE string = "" // "" means `os.Stdout`
)

type Config struct {
	// server conf
	ServerHost           string
	ServerPort           int
	TCPConnectionTimeout int
	MaxRequestsPerConn   int

	// logger conf
	LogLevel      string
	LogOutputFile string
}

func InitConfig() (*Config, error) {
	cfg := &Config{}

	if err := parseCLI(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

var ipv4Pattern = regexp.MustCompile(`^(\b25[0-5]|\b2[0-4][0-9]|\b[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$`)

func parseCLI(c *Config) error {
	flag.StringVar(
		&c.ServerHost, "host", DEFAULT_SERVER_HOST, "Host to bind the server",
	)
	flag.IntVar(
		&c.ServerPort, "port", DEFAULT_SERVER_PORT, "Port to bind the server",
	)
	flag.IntVar(
		&c.TCPConnectionTimeout,
		"timeout",
		DEFAULT_TCP_CONNECTION_TIMEOUT,
		"Deadline to associated with each connection (for read & write operations) (min: 0 , max: 15)",
	)
	flag.IntVar(
		&c.MaxRequestsPerConn,
		"max-req-per-conn",
		DEFAULT_MAX_REQUESTS_PER_CONN,
		"Maximum number of requests per connection (when keep-alive) (min: 0 , max: 200)",
	)

	flag.StringVar(
		&c.LogLevel, "log-level", DEFAULT_LOG_LEVEL, "Logging level (OPTIONS: debug, info, warning, error, fatal)",
	)
	flag.StringVar(
		&c.LogOutputFile, "log-output", DEFAULT_LOG_OUTPUT_FILE, "Path to the log file. e.g. `./app.log` (default: os.Stdout)",
	)

	flag.Parse()

	// validate host input
	switch {
	case c.ServerHost == "localhost":
	case ipv4Pattern.MatchString(c.ServerHost):
	default:
		return fmt.Errorf("invalid -host input: '%s'. host should be 'localhost' or a valid IPv4.", c.ServerHost)
	}

	// validate port input
	if c.ServerPort <= 0 || c.ServerPort > 65535 {
		return fmt.Errorf("invalid -port input: '%d'. \nport must be between 1 and 65535", c.ServerPort)
	}

	// check timeout input
	if c.TCPConnectionTimeout <= 0 {
		c.TCPConnectionTimeout = 0
	} else if c.TCPConnectionTimeout > 15 {
		c.TCPConnectionTimeout = 15
	}

	// check max-req-per-conn input
	if c.MaxRequestsPerConn <= 0 {
		c.MaxRequestsPerConn = 0
	} else if c.MaxRequestsPerConn > 200 {
		c.MaxRequestsPerConn = 200
	}

	// validate logging level input
	switch strings.ToLower(c.LogLevel) {
	case "debug", "info", "warning", "error", "fatal":
	default:
		return fmt.Errorf(
			"invalid -log-level input: '%s'. \nvalid options: 'debug', 'info', 'warning', 'error', 'fatal'",
			c.LogLevel,
		)
	}

	// validate log output file input
	if c.LogOutputFile != "" {
		if !(strings.HasSuffix(c.LogOutputFile, ".log") || strings.HasSuffix(c.LogOutputFile, ".logs")) {
			return fmt.Errorf(
				"it's better the logger's output file ends with '.log' or '.logs'. \ncurrent input: '%s'",
				c.LogOutputFile,
			)
		}
	} // skip if LogOutputFile == "" -> output will be set to `os.Stdout`

	return nil
}
