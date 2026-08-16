package cli

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/gllm-dev/molinets/di"
	"github.com/gllm-dev/molinets/pkg/logger"
)

func NewCmdStdio() *cobra.Command {
	return &cobra.Command{
		Use:   "stdio",
		Short: "Serve MCP over stdin and stdout",
		Long:  `Serve MCP over stdin and stdout, for MCP clients that launch the server as a local process. Logs go to stderr.`,
		RunE:  runStdioServer,
	}
}

func runStdioServer(cmd *cobra.Command, _ []string) error {
	level, _ := cmd.Flags().GetString("log-level")
	logger.Configure(os.Stderr, level, false)
	l := logger.New("cli_stdio")

	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server, err := di.InitializeStdioServer(ctx)
	if err != nil {
		l.Error("failed to initialize stdio server", logger.Error(err))
		return err
	}
	if err := server.Run(ctx); err != nil {
		l.Error("stdio server error", logger.Error(err))
		return err
	}
	return nil
}
