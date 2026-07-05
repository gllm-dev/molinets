package conversion

import "github.com/gllm-dev/molinets/internal/core/domain/grinder"

type Equivalent struct {
	Grinder   *grinder.Grinder
	Reachable bool
	Setting   string
	Exact     float64
	TooFine   bool
	TooCoarse bool
}

type Result struct {
	From            *grinder.Grinder
	Setting         string
	Microns         float64
	MatchingMethods []string
	AtBurrContact   bool
	Equivalents     []Equivalent
}

type MethodRange struct {
	Grinder *grinder.Grinder
	Method  string
	MinUM   float64
	MaxUM   float64
	From    string
	Middle  string
	To      string
	Partial bool
}

func NewEquivalent(g *grinder.Grinder, um float64) Equivalent {
	e := Equivalent{Grinder: g, Reachable: g.CanGrind(um)}
	if !e.Reachable {
		e.TooFine = um < g.FinestUM()
		e.TooCoarse = !e.TooFine
		return e
	}
	x := g.FromMicrons(um)
	e.Setting = g.FormatSetting(x)
	e.Exact = x
	return e
}
