package mcp

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/gllm-dev/molinets/internal/adapters/handlers/mcp/conversion"
	"github.com/gllm-dev/molinets/internal/adapters/handlers/mcp/grinder"
)

const instructions = "Converts coffee grinder settings between popular manual and electric grinders " +
	"(Comandante, 1Zpresso, Timemore, Kingrinder, Baratza, Fellow, Niche, DF64 and more) through an " +
	"approximate particle size in microns. Use convert_grind_setting when someone gives a setting on one " +
	"grinder and wants it on another, and grind_setting_for_method for a starting point for a brew method. " +
	"Results are starting points to dial in by taste."

var Version = "dev"

func NewServer(grinders *grinder.Handler, conversions *conversion.Handler) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{Name: "molinets", Title: "Molinets", Version: Version},
		&mcp.ServerOptions{Instructions: instructions})
	grinders.Register(s)
	conversions.Register(s)
	return s
}
