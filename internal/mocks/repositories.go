package mocks

import (
	"context"

	"github.com/gllm-dev/molinets/internal/core/domain/brewmethod"
	"github.com/gllm-dev/molinets/internal/core/domain/grinder"
	"github.com/gllm-dev/molinets/internal/core/ports/repositories"
)

type GrindersRepositoryMock struct {
	ListFunc   func(ctx context.Context, filter *repositories.GrinderFilter) ([]grinder.Grinder, error)
	GetFunc    func(ctx context.Context, id string) (*grinder.Grinder, error)
	SearchFunc func(ctx context.Context, query string) ([]grinder.Grinder, error)
}

var _ repositories.GrindersRepository = (*GrindersRepositoryMock)(nil)

func (m *GrindersRepositoryMock) List(ctx context.Context, filter *repositories.GrinderFilter) ([]grinder.Grinder, error) {
	if m.ListFunc == nil {
		panic("GrindersRepositoryMock.List not set")
	}
	return m.ListFunc(ctx, filter)
}

func (m *GrindersRepositoryMock) Get(ctx context.Context, id string) (*grinder.Grinder, error) {
	if m.GetFunc == nil {
		panic("GrindersRepositoryMock.Get not set")
	}
	return m.GetFunc(ctx, id)
}

func (m *GrindersRepositoryMock) Search(ctx context.Context, query string) ([]grinder.Grinder, error) {
	if m.SearchFunc == nil {
		panic("GrindersRepositoryMock.Search not set")
	}
	return m.SearchFunc(ctx, query)
}

type MethodsRepositoryMock struct {
	ListFunc func(ctx context.Context) ([]brewmethod.Method, error)
	FindFunc func(ctx context.Context, query string) (*brewmethod.Method, error)
}

var _ repositories.MethodsRepository = (*MethodsRepositoryMock)(nil)

func (m *MethodsRepositoryMock) List(ctx context.Context) ([]brewmethod.Method, error) {
	if m.ListFunc == nil {
		panic("MethodsRepositoryMock.List not set")
	}
	return m.ListFunc(ctx)
}

func (m *MethodsRepositoryMock) Find(ctx context.Context, query string) (*brewmethod.Method, error) {
	if m.FindFunc == nil {
		panic("MethodsRepositoryMock.Find not set")
	}
	return m.FindFunc(ctx, query)
}
