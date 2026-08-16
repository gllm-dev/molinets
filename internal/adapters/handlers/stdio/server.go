package stdio

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	mcp *mcp.Server
}

func NewServer(mcpServer *mcp.Server) *Server {
	return &Server{mcp: mcpServer}
}

func (s *Server) Run(ctx context.Context) error {
	err := s.mcp.Run(ctx, &mcp.StdioTransport{})
	if err == nil || errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) || strings.Contains(err.Error(), "EOF") {
		return nil
	}
	return err
}
