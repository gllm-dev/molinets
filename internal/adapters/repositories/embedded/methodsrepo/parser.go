package methodsrepo

import (
	"fmt"

	"github.com/gllm-dev/molinets/internal/core/domain/brewmethod"
	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
)

type parser struct{}

func newParser() *parser { return &parser{} }

func (p *parser) toDomain(d methodDTO) (brewmethod.Method, error) {
	if d.ID == "" || d.Name == "" || d.MinUM <= 0 || d.MaxUM <= d.MinUM {
		return brewmethod.Method{}, fmt.Errorf("%w: invalid brew method %q", domainErrors.ErrMethodNotFound, d.ID)
	}
	return brewmethod.Method{ID: d.ID, Name: d.Name, Aliases: d.Aliases, MinUM: d.MinUM, MaxUM: d.MaxUM}, nil
}
