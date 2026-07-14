package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/hamidgh01/HttpServer-from-TCP/config"
	"github.com/hamidgh01/HttpServer-from-TCP/logging"
)

type Server struct {
	host                 string
	port                 int
	listener             net.Listener
	connTimeout          time.Duration
	maxRequestsPerConn   int
	wg                   sync.WaitGroup
	shutdownNotification chan struct{}
	logger               *logging.Logger
}

func NewServer(cfg *config.Config, logger *logging.Logger) *Server {

	return &Server{
		host:                 cfg.ServerHost,
		port:                 cfg.ServerPort,
		connTimeout:          time.Second * time.Duration(cfg.TCPConnectionTimeout),
		maxRequestsPerConn:   cfg.MaxRequestsPerConn,
		shutdownNotification: make(chan struct{}),
		logger:               logger,
	}
}

func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%d", s.host, s.port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.listener = listener
	s.logger.Infof("server is running on '%s'", address)

	// Start accepting connections
	go s.acceptConnections()

	// Wait for shutdown signal
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	<-shutdownSignal
	s.logger.Info("shutdown signal received")

	return s.shutdown() // perform graceful shutdown
}

// infinitely accept connections, until receiving shutdown notification
func (s *Server) acceptConnections() {
	for {
		select {
		case <-s.shutdownNotification:
			return
		default:
		}

		if conn, err := s.listener.Accept(); err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			s.logger.Error(err.Error())
		} else {
			s.wg.Add(1)
			s.logger.Infof("connection accepted from: %s", conn.RemoteAddr().String())
			fmt.Println("connections will be handled concurrently here")
			// go s.handleConnection(conn)
			conn.Close()
			s.wg.Done()
		}
	}

}

func (s *Server) shutdown() error {

	close(s.shutdownNotification)
	if err := s.listener.Close(); err != nil {
		s.logger.Infof("Error closing listener: %s", err.Error())
	}

	// Wait for all connections to close
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("all connections closed.")
		s.logger.Info("server shutdown gracefully!")
	case <-time.After(10 * time.Second):
		s.logger.Info("shutdown timeout reached, forcing exit.")
		s.logger.Info("server shutdown!")
	}

	return nil
}
