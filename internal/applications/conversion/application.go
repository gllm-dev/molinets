package conversion

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	grinderApp "github.com/gllm-dev/molinets/internal/applications/grinder"
	"github.com/gllm-dev/molinets/internal/core/domain/brewmethod"
	"github.com/gllm-dev/molinets/internal/core/domain/conversion"
	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
	"github.com/gllm-dev/molinets/internal/core/domain/grinder"
	"github.com/gllm-dev/molinets/internal/core/ports/repositories"
	"github.com/gllm-dev/molinets/pkg/logger"
)

const MaxMicrons = 2000

type Application struct {
	grinders *grinderApp.Application
	methods  repositories.MethodsRepository
	logger   *slog.Logger
}

func New(grinders *grinderApp.Application, methods repositories.MethodsRepository) *Application {
	return &Application{
		grinders: grinders,
		methods:  methods,
		logger:   logger.New("conversion_application"),
	}
}

type ConvertCommand struct {
	FromGrinder string
	Setting     string
	ToGrinders  []string
}

func (a *Application) Convert(ctx context.Context, cmd ConvertCommand) (*conversion.Result, error) {
	src, err := a.grinders.Resolve(ctx, cmd.FromGrinder)
	if err != nil {
		return nil, err
	}
	x, err := src.ParseSetting(cmd.Setting)
	if err != nil {
		return nil, err
	}
	um := src.ToMicrons(x)

	targets, err := a.targets(ctx, src, cmd.ToGrinders)
	if err != nil {
		return nil, err
	}
	methods, err := a.methods.List(ctx)
	if err != nil {
		return nil, err
	}

	res := &conversion.Result{
		From:            src,
		Setting:         src.FormatSetting(x),
		Microns:         um,
		MatchingMethods: brewmethod.Matching(methods, um),
		AtBurrContact:   grinder.AtBurrContact(um),
		Equivalents:     make([]conversion.Equivalent, 0, len(targets)),
	}
	for i := range targets {
		res.Equivalents = append(res.Equivalents, conversion.NewEquivalent(&targets[i], um))
	}
	a.logger.DebugContext(ctx, "converted setting",
		slog.String("from", src.ID), slog.String("setting", res.Setting), slog.Float64("microns", um), slog.Int("targets", len(targets)))
	return res, nil
}

func (a *Application) targets(ctx context.Context, src *grinder.Grinder, queries []string) ([]grinder.Grinder, error) {
	if len(queries) > 0 {
		return a.grinders.ResolveAll(ctx, queries)
	}
	all, err := a.grinders.List(ctx, "")
	if err != nil {
		return nil, err
	}
	out := make([]grinder.Grinder, 0, len(all))
	for _, g := range all {
		if g.ID != src.ID {
			out = append(out, g)
		}
	}
	return out, nil
}

func (a *Application) ForMethod(ctx context.Context, grinderQuery, methodQuery string) (*conversion.MethodRange, error) {
	g, err := a.grinders.Resolve(ctx, grinderQuery)
	if err != nil {
		return nil, err
	}
	m, err := a.methods.Find(ctx, methodQuery)
	if err != nil {
		return nil, err
	}
	lo, hi := math.Max(m.MinUM, g.FinestUM()), math.Min(m.MaxUM, g.CoarsestUM())
	if lo > hi {
		return nil, fmt.Errorf("%w: %s grinds %.0f–%.0f µm and %s needs %.0f–%.0f µm",
			domainErrors.ErrMethodUnreachable, g.Name, g.FinestUM(), g.CoarsestUM(), m.Name, m.MinUM, m.MaxUM)
	}
	return &conversion.MethodRange{
		Grinder: g,
		Method:  m.Name,
		MinUM:   lo,
		MaxUM:   hi,
		From:    g.FormatSetting(g.FromMicrons(lo)),
		Middle:  g.FormatSetting(g.FromMicrons((lo + hi) / 2)),
		To:      g.FormatSetting(g.FromMicrons(hi)),
		Partial: lo > m.MinUM || hi < m.MaxUM,
	}, nil
}

func (a *Application) FromMicrons(ctx context.Context, grinderQuery string, um float64) (*conversion.Equivalent, error) {
	if um <= 0 || um > MaxMicrons || math.IsNaN(um) {
		return nil, fmt.Errorf("%w: must be between 1 and %d µm, got %v", domainErrors.ErrInvalidMicrons, MaxMicrons, um)
	}
	g, err := a.grinders.Resolve(ctx, grinderQuery)
	if err != nil {
		return nil, err
	}
	e := conversion.NewEquivalent(g, um)
	return &e, nil
}
