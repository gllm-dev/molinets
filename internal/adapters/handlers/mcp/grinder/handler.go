package grinder

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpErrors "github.com/gllm-dev/molinets/internal/adapters/handlers/mcp/errors"
	grinderApp "github.com/gllm-dev/molinets/internal/applications/grinder"
	"github.com/gllm-dev/molinets/internal/core/domain/grinder"
	"github.com/gllm-dev/molinets/pkg/logger"
)

type Handler struct {
	app    *grinderApp.Application
	parser *parser
	logger *slog.Logger
}

func New(app *grinderApp.Application) *Handler {
	return &Handler{
		app:    app,
		parser: newParser(),
		logger: logger.New("grinder_handler"),
	}
}

func (h *Handler) Register(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_grinders",
		Title:       "List grinders",
		Description: "List the supported coffee grinders (manual and electric) with their ids, aliases, setting notation, particle-size range and effective microns per step.",
		Annotations: readOnly(),
	}, h.ListGrinders)
}

func (h *Handler) ListGrinders(ctx context.Context, _ *mcp.CallToolRequest, req ListGrindersRequest) (*mcp.CallToolResult, ListGrindersResponse, error) {
	gs, err := h.app.List(ctx, grinder.Type(req.Type))
	if err != nil {
		h.logger.WarnContext(ctx, "error listing grinders", logger.Error(err))
		return nil, ListGrindersResponse{}, mcpErrors.ToolError(err)
	}
	return nil, ListGrindersResponse{Grinders: h.parser.toAPIGrinders(gs)}, nil
}

func readOnly() *mcp.ToolAnnotations {
	closed := false
	return &mcp.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true, OpenWorldHint: &closed}
}
