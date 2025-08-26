package httputil

import (
	"go-contracts/config"
	"go-contracts/util"
	"net/http"
	"time"
)

// Server HTTP服务器包装器
type Server struct {
	*http.Server
}

// New 从配置创建HTTP服务器
func New(cfg *config.HTTPServerConfig, handler http.Handler) (*Server, error) {
	// 使用配置的地址，默认使用:8080
	addr := cfg.Addr
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}

	util.Log.Info("HTTP server initialized", "address", addr)
	return &Server{server}, nil
}

// Start 启动HTTP服务器（非阻塞）
func (s *Server) Start() error {
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			util.Log.Error("HTTP server failed", "error", err)
		}
	}()
	return nil
}
