package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/gllm-dev/molinets/internal/adapters/handlers/http/healthzhdl"
	"github.com/gllm-dev/molinets/pkg/logger"
)

type Server struct {
	srv    *http.Server
	logger *slog.Logger
}

func NewServer(cfg *ServerConfig, mcpServer *mcp.Server, healthz *healthzhdl.Handler) *Server {
	l := logger.New("http_server")
	return &Server{
		srv: &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.Port),
			Handler:           withRequestLog(l, NewHandler(mcpServer, healthz)),
			ReadHeaderTimeout: 10 * time.Second,
		},
		logger: l,
	}
}

func NewHandler(mcpServer *mcp.Server, healthz *healthzhdl.Handler) http.Handler {
	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return mcpServer },
		&mcp.StreamableHTTPOptions{Stateless: true, JSONResponse: true},
	)
	mux := http.NewServeMux()
	mux.Handle("/mcp", mcpHandler)
	mux.HandleFunc("GET /healthz", healthz.Healthz)
	return mux
}

func (s *Server) Start(_ context.Context) error {
	s.logger.Info("listening", slog.String("addr", s.srv.Addr))
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func withRequestLog(l *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		if r.URL.Path != "/healthz" {
			l.InfoContext(r.Context(), "request",
				slog.String("method", r.Method), slog.String("path", r.URL.Path), slog.Duration("duration", time.Since(start)))
		}
	})
}
