package grindersrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gllm-dev/molinets/internal/adapters/repositories/embedded"
	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
	"github.com/gllm-dev/molinets/internal/core/domain/grinder"
	"github.com/gllm-dev/molinets/internal/core/ports/repositories"
	"github.com/gllm-dev/molinets/pkg/textnorm"
)

type Repository struct {
	grinders []grinder.Grinder
	byID     map[string]int
	keys     [][]string
}

var _ repositories.GrindersRepository = (*Repository)(nil)

func New() (*Repository, error) { return NewFromJSON(embedded.GrindersJSON) }

func NewFromJSON(raw []byte) (*Repository, error) {
	var dtos []grinderDTO
	if err := json.Unmarshal(raw, &dtos); err != nil {
		return nil, fmt.Errorf("%w: %v", domainErrors.ErrInvalidGrinder, err)
	}
	p := newParser()
	r := &Repository{byID: make(map[string]int, len(dtos))}
	for _, d := range dtos {
		g, err := p.toDomain(d)
		if err != nil {
			return nil, err
		}
		if _, dup := r.byID[g.ID]; dup {
			return nil, fmt.Errorf("%w: duplicate id %s", domainErrors.ErrInvalidGrinder, g.ID)
		}
		r.byID[g.ID] = len(r.grinders)
		r.grinders = append(r.grinders, g)
		keys := []string{textnorm.Key(g.ID), textnorm.Key(g.Name)}
		for _, a := range g.Aliases {
			keys = append(keys, textnorm.Key(a))
		}
		r.keys = append(r.keys, keys)
	}
	return r, nil
}

func (r *Repository) List(_ context.Context, filter *repositories.GrinderFilter) ([]grinder.Grinder, error) {
	out := make([]grinder.Grinder, 0, len(r.grinders))
	for _, g := range r.grinders {
		if filter != nil && filter.Type != "" && g.Type != filter.Type {
			continue
		}
		out = append(out, g)
	}
	return out, nil
}

func (r *Repository) Get(_ context.Context, id string) (*grinder.Grinder, error) {
	i, ok := r.byID[id]
	if !ok {
		return nil, fmt.Errorf("%w: %q", domainErrors.ErrGrinderNotFound, id)
	}
	g := r.grinders[i]
	return &g, nil
}

func (r *Repository) Search(_ context.Context, query string) ([]grinder.Grinder, error) {
	q := textnorm.Key(query)
	if q == "" {
		return []grinder.Grinder{}, nil
	}
	for i, keys := range r.keys {
		for _, k := range keys {
			if k == q {
				return []grinder.Grinder{r.grinders[i]}, nil
			}
		}
	}
	var out []grinder.Grinder
	for i, keys := range r.keys {
		for _, k := range keys {
			if strings.Contains(k, q) {
				out = append(out, r.grinders[i])
				break
			}
		}
	}
	if out == nil {
		out = []grinder.Grinder{}
	}
	return out, nil
}
