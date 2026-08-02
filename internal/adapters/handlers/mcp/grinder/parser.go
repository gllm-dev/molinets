package grinder

import (
	"math"

	"github.com/gllm-dev/molinets/internal/core/domain/grinder"
)

type parser struct{}

func newParser() *parser { return &parser{} }

func (p *parser) toAPIGrinder(g *grinder.Grinder) Grinder {
	aliases := g.Aliases
	if aliases == nil {
		aliases = []string{}
	}
	return Grinder{
		ID:                    g.ID,
		Name:                  g.Name,
		Type:                  string(g.Type),
		Aliases:               aliases,
		SettingNotation:       g.NotationHint,
		RangeUM:               [2]float64{g.FinestUM(), g.CoarsestUM()},
		EffectiveUMPerStep:    math.Round(g.UMPerStep(g.MinSetting())*10) / 10,
		ManufacturerUMPerStep: g.ManufacturerUMPerStep,
	}
}

func (p *parser) toAPIGrinders(gs []grinder.Grinder) []Grinder {
	out := make([]Grinder, len(gs))
	for i := range gs {
		out[i] = p.toAPIGrinder(&gs[i])
	}
	return out
}
