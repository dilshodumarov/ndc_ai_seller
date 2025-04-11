// Package httpserver implements HTTP server.
package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	Engine *gin.Engine
	Server *http.Server
	notify chan error

	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration
}

// New -.
func New(opts ...Option) *Server {
	s := &Server{
		Engine:          nil,
		Server:          nil,
		notify:          make(chan error, 1),
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	// Create a Gin engine
	engine := gin.New()

	// Optional middlewares
	engine.Use(gin.Logger(), gin.Recovery())

	// Create an http.Server wrapping gin
	server := &http.Server{
		Addr:         s.address,
		Handler:      engine,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}

	s.Engine = engine
	s.Server = server

	return s
}

// Start -.
func (s *Server) Start() {
	go func() {
		s.notify <- s.Server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.Server.Shutdown(ctx)
}
