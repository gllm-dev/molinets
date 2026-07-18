package methodsrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gllm-dev/molinets/internal/adapters/repositories/embedded"
	"github.com/gllm-dev/molinets/internal/core/domain/brewmethod"
	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
	"github.com/gllm-dev/molinets/internal/core/ports/repositories"
	"github.com/gllm-dev/molinets/pkg/textnorm"
)

type Repository struct {
	methods []brewmethod.Method
}

var _ repositories.MethodsRepository = (*Repository)(nil)

func New() (*Repository, error) {
	var dtos []methodDTO
	if err := json.Unmarshal(embedded.MethodsJSON, &dtos); err != nil {
		return nil, fmt.Errorf("brew methods: %w", err)
	}
	p := newParser()
	r := &Repository{methods: make([]brewmethod.Method, 0, len(dtos))}
	for _, d := range dtos {
		m, err := p.toDomain(d)
		if err != nil {
			return nil, err
		}
		r.methods = append(r.methods, m)
	}
	return r, nil
}

func (r *Repository) List(context.Context) ([]brewmethod.Method, error) {
	return append([]brewmethod.Method(nil), r.methods...), nil
}

func (r *Repository) Find(_ context.Context, query string) (*brewmethod.Method, error) {
	q := textnorm.Key(query)
	if q != "" {
		for i := range r.methods {
			m := &r.methods[i]
			if textnorm.Key(m.ID) == q || textnorm.Key(m.Name) == q {
				return m, nil
			}
			for _, a := range m.Aliases {
				if textnorm.Key(a) == q {
					return m, nil
				}
			}
		}
		for i := range r.methods {
			if strings.HasPrefix(textnorm.Key(r.methods[i].Name), q) {
				return &r.methods[i], nil
			}
		}
	}
	ids := make([]string, len(r.methods))
	for i, m := range r.methods {
		ids[i] = m.ID
	}
	return nil, fmt.Errorf("%w: %q (use one of: %s)", domainErrors.ErrMethodNotFound, query, strings.Join(ids, ", "))
}
