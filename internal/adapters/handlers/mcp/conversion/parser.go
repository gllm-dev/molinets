package conversion

import (
	"fmt"
	"math"

	"github.com/gllm-dev/molinets/internal/core/domain/conversion"
)

const (
	startingPointNote = "Starting point to dial in by taste, not an exact value: zero points differ by 1–2 clicks between units of the same model. Count every setting from burr contact (zero)."
	burrContactNote   = "This setting is at or next to burr contact: too fine to brew with, and grinding there can damage the burrs. "
	partialRangeNote  = "This grinder covers only part of the method's range. "
)

type parser struct{}

func newParser() *parser { return &parser{} }

func round5(um float64) int { return int(math.Round(um/5) * 5) }

func (p *parser) toAPIEquivalent(e *conversion.Equivalent) Equivalent {
	out := Equivalent{
		Grinder:         e.Grinder.Name,
		ID:              e.Grinder.ID,
		InRange:         e.Reachable,
		SettingNotation: e.Grinder.NotationHint,
	}
	switch {
	case e.Reachable:
		out.Setting = e.Setting
		out.ExactValue = math.Round(e.Exact*100) / 100
	case e.TooFine:
		out.Note = fmt.Sprintf("cannot grind this fine (finest ≈ %.0f µm)", e.Grinder.FinestUM())
	default:
		out.Note = fmt.Sprintf("cannot grind this coarse (coarsest ≈ %.0f µm)", e.Grinder.CoarsestUM())
	}
	return out
}

func (p *parser) toAPIConvertResponse(r *conversion.Result) ConvertResponse {
	note := startingPointNote
	if r.AtBurrContact {
		note = burrContactNote + note
	}
	out := ConvertResponse{
		From:            r.From.Name,
		Setting:         r.Setting,
		ApproxMicrons:   round5(r.Microns),
		MatchingMethods: r.MatchingMethods,
		Equivalents:     make([]Equivalent, len(r.Equivalents)),
		Note:            note,
	}
	for i := range r.Equivalents {
		out.Equivalents[i] = p.toAPIEquivalent(&r.Equivalents[i])
	}
	return out
}

func (p *parser) toAPIMethodResponse(r *conversion.MethodRange) MethodResponse {
	note := startingPointNote
	if r.Partial {
		note = partialRangeNote + note
	}
	return MethodResponse{
		Grinder:         r.Grinder.Name,
		Method:          r.Method,
		MicronRange:     fmt.Sprintf("%.0f–%.0f µm", r.MinUM, r.MaxUM),
		SettingFrom:     r.From,
		SettingMiddle:   r.Middle,
		SettingTo:       r.To,
		SettingNotation: r.Grinder.NotationHint,
		Note:            note,
	}
}
