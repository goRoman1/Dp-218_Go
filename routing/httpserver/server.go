package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"
)

const(
	defaultReadTimeout = 5*time.Second
	defaultWriteTimeout = 5* time.Second
	defaultIdleTimeout = 30*time.Second
	defaultShutdownTimeout = 3 * time.Second
	defaultAddr = ":8080"
)

type Server struct{
	server *http.Server
	notify chan error
	shutdownTimeout time.Duration
}

type Option func(*Server)

func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
		Addr:         defaultAddr,
	}

	server := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt:=range opts {
		opt(server)
	}

	server.start()

	return server
}

func (s *Server) start(){
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error{
	return s.notify
}

func (s *Server) Shutdown() error{
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}



func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func IdleTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.IdleTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}