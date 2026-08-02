package cli

import (
	"github.com/spf13/cobra"

	mcpServer "github.com/gllm-dev/molinets/internal/adapters/handlers/mcp"
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "molinets",
		Short:         "MCP server that converts coffee grinder settings",
		Long:          `MCP server that converts coffee grinder settings between popular manual and electric grinders, through an approximate particle size in microns.`,
		Version:       mcpServer.Version,
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	cmd.PersistentFlags().String("log-level", "info", "log level: debug, info, warn or error")

	cmd.AddCommand(NewCmdHTTP())

	return cmd
}
