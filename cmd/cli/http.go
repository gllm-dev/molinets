package cli

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/gllm-dev/molinets/di"
	"github.com/gllm-dev/molinets/pkg/logger"
)

const shutdownTimeout = 10 * time.Second

func NewCmdHTTP() *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Serve MCP over Streamable HTTP",
		Long:  `Serve MCP over Streamable HTTP on /mcp (stateless, JSON responses) with a health check on /healthz. The port comes from PORT (default 8080).`,
		RunE:  runHTTPServer,
	}
}

func runHTTPServer(cmd *cobra.Command, _ []string) error {
	level, _ := cmd.Flags().GetString("log-level")
	logger.Configure(os.Stdout, level, true)
	l := logger.New("cli_http")

	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server, err := di.InitializeHTTPServer(ctx)
	if err != nil {
		l.Error("failed to initialize http server", logger.Error(err))
		return err
	}

	errCh := make(chan error, 1)
	go func() { errCh <- server.Start(ctx) }()

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("http server error", logger.Error(err))
			return err
		}
	case <-ctx.Done():
		l.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := server.Stop(shutdownCtx); err != nil {
			l.Error("graceful shutdown failed", logger.Error(err))
			return err
		}
	}
	l.Info("http server stopped")
	return nil
}
