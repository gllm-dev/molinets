package conversion

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	mcpErrors "github.com/gllm-dev/molinets/internal/adapters/handlers/mcp/errors"
	conversionApp "github.com/gllm-dev/molinets/internal/applications/conversion"
	"github.com/gllm-dev/molinets/pkg/logger"
)

type Handler struct {
	app    *conversionApp.Application
	parser *parser
	logger *slog.Logger
}

func New(app *conversionApp.Application) *Handler {
	return &Handler{
		app:    app,
		parser: newParser(),
		logger: logger.New("conversion_handler"),
	}
}

func (h *Handler) Register(s *mcp.Server) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "convert_grind_setting",
		Title:       "Convert grind setting",
		Description: "Convert a grind setting from one coffee grinder to equivalent settings on other grinders, via an approximate particle size in microns. Also returns the brew methods that suit that size.",
		Annotations: readOnly(),
	}, h.Convert)
	mcp.AddTool(s, &mcp.Tool{
		Name:        "grind_setting_for_method",
		Title:       "Grind setting for a brew method",
		Description: "Give the starting setting range on one grinder for a brew method (espresso, V60, AeroPress, French press, Chemex…).",
		Annotations: readOnly(),
	}, h.ForMethod)
	mcp.AddTool(s, &mcp.Tool{
		Name:        "microns_to_setting",
		Title:       "Microns to setting",
		Description: "Give the setting on one grinder that produces an approximate particle size in microns.",
		Annotations: readOnly(),
	}, h.FromMicrons)
}

func (h *Handler) Convert(ctx context.Context, _ *mcp.CallToolRequest, req ConvertRequest) (*mcp.CallToolResult, ConvertResponse, error) {
	res, err := h.app.Convert(ctx, conversionApp.ConvertCommand{
		FromGrinder: req.FromGrinder,
		Setting:     req.Setting,
		ToGrinders:  req.ToGrinders,
	})
	if err != nil {
		h.logger.InfoContext(ctx, "convert rejected", logger.Error(err))
		return nil, ConvertResponse{}, mcpErrors.ToolError(err)
	}
	return nil, h.parser.toAPIConvertResponse(res), nil
}

func (h *Handler) ForMethod(ctx context.Context, _ *mcp.CallToolRequest, req MethodRequest) (*mcp.CallToolResult, MethodResponse, error) {
	res, err := h.app.ForMethod(ctx, req.Grinder, req.Method)
	if err != nil {
		h.logger.InfoContext(ctx, "method lookup rejected", logger.Error(err))
		return nil, MethodResponse{}, mcpErrors.ToolError(err)
	}
	return nil, h.parser.toAPIMethodResponse(res), nil
}

func (h *Handler) FromMicrons(ctx context.Context, _ *mcp.CallToolRequest, req MicronsRequest) (*mcp.CallToolResult, Equivalent, error) {
	res, err := h.app.FromMicrons(ctx, req.Grinder, req.Microns)
	if err != nil {
		h.logger.InfoContext(ctx, "microns lookup rejected", logger.Error(err))
		return nil, Equivalent{}, mcpErrors.ToolError(err)
	}
	return nil, h.parser.toAPIEquivalent(res), nil
}

func readOnly() *mcp.ToolAnnotations {
	closed := false
	return &mcp.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true, OpenWorldHint: &closed}
}
