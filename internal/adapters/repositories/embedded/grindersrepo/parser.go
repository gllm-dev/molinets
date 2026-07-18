package grindersrepo

import "github.com/gllm-dev/molinets/internal/core/domain/grinder"

type parser struct{}

func newParser() *parser { return &parser{} }

func (p *parser) toDomain(d grinderDTO) (grinder.Grinder, error) {
	segments := make([]grinder.Segment, len(d.Calibration))
	for i, s := range d.Calibration {
		segments[i] = grinder.Segment{From: s.From, OffsetUM: s.OffsetUM, UMPerUnit: s.UMPerUnit}
	}
	cal, err := grinder.NewCalibration(segments)
	if err != nil {
		return grinder.Grinder{}, err
	}
	g := grinder.Grinder{
		ID:           d.ID,
		Name:         d.Name,
		Aliases:      d.Aliases,
		Type:         grinder.Type(d.Type),
		NotationHint: d.NotationHint,
		Notation: grinder.Notation{
			Kind:              grinder.NotationKind(d.Notation.Kind),
			ClicksPerRotation: d.Notation.ClicksPerRotation,
			ClicksPerNumber:   d.Notation.ClicksPerNumber,
			Step:              d.Notation.Step,
		},
		Calibration:           cal,
		Range:                 grinder.Range{MinUM: d.RangeUM[0], MaxUM: d.RangeUM[1]},
		ManufacturerUMPerStep: d.ManufacturerUMPerStep,
		FitRMSEUM:             d.FitRMSEUM,
		Reference:             d.Reference,
	}
	if err := g.Validate(); err != nil {
		return grinder.Grinder{}, err
	}
	return g, nil
}
