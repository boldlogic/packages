package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Server оборачивает стандартный HTTP-сервер.
type Server struct {
	httpServer *http.Server
}

func (conf *ServerConfig) getAddr() string {
	return fmt.Sprintf("%s:%d", conf.ListenIp, conf.ListenPort)
}

// NewServer создаёт HTTP-сервер с переданным обработчиком и настройками.
func NewServer(handler http.Handler, conf ServerConfig) *Server {
	srv := &Server{
		httpServer: &http.Server{
			Addr:              conf.getAddr(),
			Handler:           handler,
			ReadTimeout:       time.Duration(conf.Opts.ReadTimeout) * time.Second,
			ReadHeaderTimeout: time.Duration(conf.Opts.ReadHeaderTimeout) * time.Second,
			WriteTimeout:      time.Duration(conf.Opts.WriteTimeout) * time.Second,
			IdleTimeout:       time.Duration(conf.Opts.IdleTimeout) * time.Second,
		},
	}

	return srv
}

// ListenAndServe запускает HTTP-сервер.
func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown корректно завершает HTTP-сервер с учётом переданного контекста.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
