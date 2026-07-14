package main

import (
	"fmt"
	"os"

	"github.com/hamidgh01/HttpServer-from-TCP/config"
	"github.com/hamidgh01/HttpServer-from-TCP/logging"
	"github.com/hamidgh01/HttpServer-from-TCP/server"
)

func main() {
	// init configurations
	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Printf("failed to init configurations. reason: %s", err)
		os.Exit(1)
	}

	logger := logging.NewLogger(cfg)

	// create and run server
	server := server.NewServer(cfg, logger)
	if err := server.Start(); err != nil {
		logger.Fatalf("failed to run Server. reason: %s", err)
	}
}
