package grinder

import (
	"context"
	"errors"
	"strings"
	"testing"

	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
	"github.com/gllm-dev/molinets/internal/core/domain/grinder"
	"github.com/gllm-dev/molinets/internal/core/ports/repositories"
	"github.com/gllm-dev/molinets/internal/mocks"
)

var ctx = context.Background()

func g(id string) grinder.Grinder { return grinder.Grinder{ID: id, Name: id} }

func TestResolve(t *testing.T) {
	repo := &mocks.GrindersRepositoryMock{SearchFunc: func(_ context.Context, q string) ([]grinder.Grinder, error) {
		switch q {
		case "k6":
			return []grinder.Grinder{g("kingrinder-k6")}, nil
		case "hario":
			return []grinder.Grinder{g("hario-skerton-pro"), g("hario-mini-mill-plus")}, nil
		default:
			return []grinder.Grinder{}, nil
		}
	}}
	app := New(repo)

	if got, err := app.Resolve(ctx, "k6"); err != nil || got.ID != "kingrinder-k6" {
		t.Fatalf("k6: %v %v", got, err)
	}
	_, err := app.Resolve(ctx, "hario")
	if !errors.Is(err, domainErrors.ErrAmbiguousGrinder) || !strings.Contains(err.Error(), "hario-mini-mill-plus, hario-skerton-pro") {
		t.Fatalf("hario: %v", err)
	}
	if _, err := app.Resolve(ctx, "bialetti"); !errors.Is(err, domainErrors.ErrGrinderNotFound) {
		t.Fatalf("bialetti: %v", err)
	}
}

func TestResolveEmptyNeverHitsTheRepository(t *testing.T) {
	app := New(&mocks.GrindersRepositoryMock{})
	if _, err := app.Resolve(ctx, "  "); !errors.Is(err, domainErrors.ErrEmptyGrinderQuery) {
		t.Fatalf("got %v", err)
	}
}

func TestResolveAllDropsRepeats(t *testing.T) {
	repo := &mocks.GrindersRepositoryMock{SearchFunc: func(context.Context, string) ([]grinder.Grinder, error) {
		return []grinder.Grinder{g("comandante-c40")}, nil
	}}
	got, err := New(repo).ResolveAll(ctx, []string{"comandante", "c40", "Comandante C40 MK4"})
	if err != nil || len(got) != 1 {
		t.Fatalf("got %d grinders, %v; want 1", len(got), err)
	}
}

func TestListValidatesType(t *testing.T) {
	var gotFilter *repositories.GrinderFilter
	repo := &mocks.GrindersRepositoryMock{ListFunc: func(_ context.Context, f *repositories.GrinderFilter) ([]grinder.Grinder, error) {
		gotFilter = f
		return nil, nil
	}}
	app := New(repo)
	if _, err := app.List(ctx, "hybrid"); !errors.Is(err, domainErrors.ErrInvalidGrinder) {
		t.Fatalf("hybrid: %v", err)
	}
	if _, _ = app.List(ctx, ""); gotFilter != nil {
		t.Fatal("no type should mean no filter")
	}
	if _, _ = app.List(ctx, grinder.TypeManual); gotFilter == nil || gotFilter.Type != grinder.TypeManual {
		t.Fatal("manual filter not passed")
	}
}
