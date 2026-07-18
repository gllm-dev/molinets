package repositories

import (
	"context"

	"github.com/gllm-dev/molinets/internal/core/domain/brewmethod"
)

type MethodsRepository interface {
	List(ctx context.Context) ([]brewmethod.Method, error)
	Find(ctx context.Context, query string) (*brewmethod.Method, error)
}
