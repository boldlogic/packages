package middleware

import (
	"github.com/boldlogic/packages/transport/httpserver/httpmetrics"
	"go.uber.org/zap"
)

// Middleware связывает логгер и HTTP-метрики для обёрток вокруг http.Handler.
type Middleware struct {
	logger  *zap.Logger
	metrics *httpmetrics.HTTPMetrics
}

// NewMiddleware создаёт цепочку middleware с заданными метриками и логгером.
func NewMiddleware(metrics *httpmetrics.HTTPMetrics, logger *zap.Logger) *Middleware {
	return &Middleware{
		logger:  logger,
		metrics: metrics,
	}
}
