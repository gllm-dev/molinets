package grinder

import (
	"fmt"
	"math"

	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
)

type Type string

const (
	TypeManual   Type = "manual"
	TypeElectric Type = "electric"
)

const MinFinestUM = 40

const rangeToleranceUM = 10

type Range struct {
	MinUM float64
	MaxUM float64
}

type Grinder struct {
	ID                    string
	Name                  string
	Aliases               []string
	Type                  Type
	NotationHint          string
	Notation              Notation
	Calibration           Calibration
	Range                 Range
	ManufacturerUMPerStep string
	FitRMSEUM             *float64
	Reference             string
}

func (g *Grinder) Validate() error {
	if g.ID == "" || g.Name == "" {
		return fmt.Errorf("%w: id and name are required", domainErrors.ErrInvalidGrinder)
	}
	if g.Type != TypeManual && g.Type != TypeElectric {
		return fmt.Errorf("%w: %s has unknown type %q", domainErrors.ErrInvalidGrinder, g.ID, g.Type)
	}
	if g.Range.MinUM < 0 || g.Range.MaxUM <= g.Range.MinUM {
		return fmt.Errorf("%w: %s has an invalid range", domainErrors.ErrInvalidGrinder, g.ID)
	}
	if _, err := NewNotation(g.Notation); err != nil {
		return fmt.Errorf("%s: %w", g.ID, err)
	}
	if _, err := NewCalibration(g.Calibration.Segments); err != nil {
		return fmt.Errorf("%s: %w", g.ID, err)
	}
	return nil
}

func (g *Grinder) FinestUM() float64 { return math.Max(g.Range.MinUM, MinFinestUM) }

func (g *Grinder) CoarsestUM() float64 { return g.Range.MaxUM }

func (g *Grinder) ToMicrons(x float64) float64 { return g.Calibration.Microns(x) }

func (g *Grinder) FromMicrons(um float64) float64 { return g.Calibration.Setting(um) }

func (g *Grinder) CanGrind(um float64) bool {
	return um >= g.FinestUM()-rangeToleranceUM &&
		um <= g.CoarsestUM()+rangeToleranceUM &&
		g.FromMicrons(um) >= 0
}

func (g *Grinder) MinSetting() float64 { return math.Max(0, g.FromMicrons(g.Range.MinUM)) }
func (g *Grinder) MaxSetting() float64 { return g.FromMicrons(g.Range.MaxUM) }

func (g *Grinder) ParseSetting(raw string) (float64, error) {
	x, err := g.Notation.Parse(raw)
	if err != nil {
		return 0, fmt.Errorf("%s (%s): %w", g.Name, g.NotationHint, err)
	}
	tol := g.Notation.StepSize() + 1e-9
	if x < g.MinSetting()-tol || x > g.MaxSetting()+tol {
		return 0, fmt.Errorf("%w: %s on %s (usable %s to %s)", domainErrors.ErrSettingOutOfRange,
			g.FormatSetting(x), g.Name, g.FormatSetting(g.MinSetting()), g.FormatSetting(g.MaxSetting()))
	}
	return x, nil
}

func (g *Grinder) FormatSetting(x float64) string { return g.Notation.Format(x) }

func (g *Grinder) UMPerStep(x float64) float64 {
	return g.Calibration.SlopeAt(x) * g.Notation.StepSize()
}

func AtBurrContact(um float64) bool { return um < MinFinestUM }
