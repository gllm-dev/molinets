package repositories

import (
	"context"

	"github.com/gllm-dev/molinets/internal/core/domain/grinder"
)

type GrinderFilter struct {
	Type grinder.Type
}

type GrindersRepository interface {
	List(ctx context.Context, filter *GrinderFilter) ([]grinder.Grinder, error)
	Get(ctx context.Context, id string) (*grinder.Grinder, error)
	Search(ctx context.Context, query string) ([]grinder.Grinder, error)
}
