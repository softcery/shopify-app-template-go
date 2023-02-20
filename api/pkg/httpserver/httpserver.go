package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	"log"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultMaxHeaderBytes  = 1 << 20
	_defaultShutdownTimeout = 3 * time.Second
)

// Server - represents http server.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// Option - represents http server option.
type Option func(*Server)

// Port - configures http server port.
func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

// ReadTimeout - configures http server read timeout.
func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

// WriteTimeout - configures http server read timeout.
func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

// ShutdownTimeout - configures http server shutdown timeout.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

// New - creates instance of new http server.
func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Addr:           _defaultAddr,
		Handler:        handler,
		ReadTimeout:    _defaultReadTimeout,
		WriteTimeout:   _defaultWriteTimeout,
		MaxHeaderBytes: _defaultMaxHeaderBytes,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}

	// add custom options
	for _, opt := range opts {
		opt(s)
	}

	s.start()

	return s
}

// Start - bootstraps http server.
func (s *Server) start() {
	log.Printf("Starting HTTP server on port %s", s.server.Addr)
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify - returns error notification channel.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown - shuts down http server gracefully.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
