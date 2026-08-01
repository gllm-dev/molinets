package grinder

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
	"github.com/gllm-dev/molinets/internal/core/domain/grinder"
	"github.com/gllm-dev/molinets/internal/core/ports/repositories"
	"github.com/gllm-dev/molinets/pkg/logger"
)

type Application struct {
	grinders repositories.GrindersRepository
	logger   *slog.Logger
}

func New(grinders repositories.GrindersRepository) *Application {
	return &Application{
		grinders: grinders,
		logger:   logger.New("grinder_application"),
	}
}

func (a *Application) List(ctx context.Context, t grinder.Type) ([]grinder.Grinder, error) {
	var filter *repositories.GrinderFilter
	if t != "" {
		if t != grinder.TypeManual && t != grinder.TypeElectric {
			return nil, fmt.Errorf("%w: type must be %q or %q", domainErrors.ErrInvalidGrinder, grinder.TypeManual, grinder.TypeElectric)
		}
		filter = &repositories.GrinderFilter{Type: t}
	}
	return a.grinders.List(ctx, filter)
}

func (a *Application) Resolve(ctx context.Context, query string) (*grinder.Grinder, error) {
	if strings.TrimSpace(query) == "" {
		return nil, domainErrors.ErrEmptyGrinderQuery
	}
	matches, err := a.grinders.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("%w: %q", domainErrors.ErrGrinderNotFound, query)
	case 1:
		return &matches[0], nil
	default:
		ids := make([]string, len(matches))
		for i := range matches {
			ids[i] = matches[i].ID
		}
		sort.Strings(ids)
		a.logger.DebugContext(ctx, "ambiguous grinder query", slog.String("query", query), slog.Int("matches", len(ids)))
		return nil, fmt.Errorf("%w: %q could be %s", domainErrors.ErrAmbiguousGrinder, query, strings.Join(ids, ", "))
	}
}

func (a *Application) ResolveAll(ctx context.Context, queries []string) ([]grinder.Grinder, error) {
	seen := make(map[string]bool, len(queries))
	out := make([]grinder.Grinder, 0, len(queries))
	for _, q := range queries {
		g, err := a.Resolve(ctx, q)
		if err != nil {
			return nil, err
		}
		if seen[g.ID] {
			continue
		}
		seen[g.ID] = true
		out = append(out, *g)
	}
	return out, nil
}
