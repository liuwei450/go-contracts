package httputil

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	DefaultReadTimeout  = 10 * time.Second
	DefaultWriteTimeout = 10 * time.Second
	DefaultIdleTimeout  = 30 * time.Second
)

type HTTPServer struct {
	listener net.Listener
	srv      *http.Server
	closed   atomic.Bool
}

type HTTPOption func(srv *HTTPServer) error

func StartServerWithDefaults(addr string, handler http.Handler) (*HTTPServer, error) {
	return StartHTTPServer(addr, handler,
		WithTimeouts(DefaultReadTimeout, DefaultWriteTimeout, DefaultIdleTimeout),
	)
}

func StartHTTPServer(addr string, handler http.Handler, opts ...HTTPOption) (*HTTPServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to address %q: %w", addr, err)
	}

	srvCtx, srvCancel := context.WithCancel(context.Background())
	srv := &http.Server{
		Handler: handler,
		BaseContext: func(listener net.Listener) context.Context {
			return srvCtx
		},
	}
	out := &HTTPServer{listener: listener, srv: srv}

	for _, opt := range opts {
		if err := opt(out); err != nil {
			srvCancel()
			return nil, errors.Join(fmt.Errorf("failed to apply HTTP option: %w", err), listener.Close())
		}
	}

	go func() {
		err := out.srv.Serve(listener)
		srvCancel()
		if errors.Is(err, http.ErrServerClosed) {
			out.closed.Store(true)
		} else {
			panic(fmt.Errorf("unexpected serve error: %w", err))
		}
	}()

	return out, nil
}

func WithTimeouts(read, write, idle time.Duration) HTTPOption {
	return func(srv *HTTPServer) error {
		srv.srv.ReadTimeout = read
		srv.srv.WriteTimeout = write
		srv.srv.IdleTimeout = idle
		return nil
	}
}

func (s *HTTPServer) Closed() bool {
	return s.closed.Load()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	if err := s.Shutdown(ctx); err != nil {
		if errors.Is(err, ctx.Err()) {
			return s.Close()
		}
		return err
	}
	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *HTTPServer) Close() error {
	return s.srv.Close()
}

func (s *HTTPServer) Addr() net.Addr {
	return s.listener.Addr()
}
