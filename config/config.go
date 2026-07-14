package config

import (
	"flag"
	"fmt"
	"regexp"
)

// Default configuration settings
const (
	DEFAULT_SERVER_HOST            string = "127.0.0.1"
	DEFAULT_SERVER_PORT            int    = 8000
	DEFAULT_TCP_CONNECTION_TIMEOUT int    = 5 // seconds
)

type Config struct {
	ServerHost           string
	ServerPort           int
	TCPConnectionTimeout int
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
		"Deadline to associated with each connection (for read & write operations)",
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

	return nil
}
