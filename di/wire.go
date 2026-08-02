//go:build wireinject
// +build wireinject

package di

import (
	"context"

	"github.com/google/wire"

	httpServer "github.com/gllm-dev/molinets/internal/adapters/handlers/http"
	"github.com/gllm-dev/molinets/internal/adapters/handlers/http/healthzhdl"
	mcpServer "github.com/gllm-dev/molinets/internal/adapters/handlers/mcp"
	conversionHandler "github.com/gllm-dev/molinets/internal/adapters/handlers/mcp/conversion"
	grinderHandler "github.com/gllm-dev/molinets/internal/adapters/handlers/mcp/grinder"
	"github.com/gllm-dev/molinets/internal/adapters/repositories/embedded/grindersrepo"
	"github.com/gllm-dev/molinets/internal/adapters/repositories/embedded/methodsrepo"
	conversionApp "github.com/gllm-dev/molinets/internal/applications/conversion"
	grinderApp "github.com/gllm-dev/molinets/internal/applications/grinder"
	"github.com/gllm-dev/molinets/internal/core/ports/repositories"
)

var repositoriesSet = wire.NewSet(
	grindersrepo.New,
	wire.Bind(new(repositories.GrindersRepository), new(*grindersrepo.Repository)),
	methodsrepo.New,
	wire.Bind(new(repositories.MethodsRepository), new(*methodsrepo.Repository)),
)

var applicationsSet = wire.NewSet(
	grinderApp.New,
	conversionApp.New,
)

var mcpSet = wire.NewSet(
	grinderHandler.New,
	conversionHandler.New,
	mcpServer.NewServer,
)

func InitializeHTTPServer(ctx context.Context) (*httpServer.Server, error) {
	wire.Build(
		repositoriesSet,
		applicationsSet,
		mcpSet,
		healthzhdl.New,
		httpServer.LoadServerConfig,
		httpServer.NewServer,
	)
	return nil, nil
}
