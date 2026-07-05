package grinder

import (
	"errors"
	"testing"

	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
)

func comandante(t *testing.T) *Grinder {
	t.Helper()
	c, err := NewCalibration([]Segment{{OffsetUM: -9.3, UMPerUnit: 27.616}})
	if err != nil {
		t.Fatal(err)
	}
	g := &Grinder{ID: "comandante-c40", Name: "Comandante C40 MK4", Type: TypeManual,
		NotationHint: "clicks from zero", Notation: Notation{Kind: NotationSteps},
		Calibration: c, Range: Range{MinUM: 0, MaxUM: 1090}}
	if err := g.Validate(); err != nil {
		t.Fatal(err)
	}
	return g
}

func TestParseSettingRange(t *testing.T) {
	g := comandante(t)
	for _, s := range []string{"0", "1", "26", "40"} {
		if _, err := g.ParseSetting(s); err != nil {
			t.Errorf("%s should be valid: %v", s, err)
		}
	}
	for _, s := range []string{"42", "60", "200"} {
		if _, err := g.ParseSetting(s); !errors.Is(err, domainErrors.ErrSettingOutOfRange) {
			t.Errorf("%s: expected ErrSettingOutOfRange, got %v", s, err)
		}
	}
	if _, err := g.ParseSetting("26.5"); !errors.Is(err, domainErrors.ErrInvalidSetting) {
		t.Errorf("26.5: expected ErrInvalidSetting, got %v", err)
	}
}

func TestNoGrinderClaimsSubTurkishSizes(t *testing.T) {
	g := comandante(t)
	if g.FinestUM() != MinFinestUM {
		t.Fatalf("finest = %v, want %v", g.FinestUM(), MinFinestUM)
	}
	if g.CanGrind(20) {
		t.Error("should not claim 20 µm")
	}
	if !g.CanGrind(60) || !g.CanGrind(1090) || g.CanGrind(1200) {
		t.Error("range edges are wrong")
	}
}

func TestBurrContact(t *testing.T) {
	if !AtBurrContact(0) || !AtBurrContact(39) || AtBurrContact(40) {
		t.Error("burr contact threshold is wrong")
	}
}

func TestValidate(t *testing.T) {
	g := comandante(t)
	broken := []func(*Grinder){
		func(g *Grinder) { g.ID = "" },
		func(g *Grinder) { g.Type = "hybrid" },
		func(g *Grinder) { g.Range = Range{MinUM: 500, MaxUM: 400} },
		func(g *Grinder) { g.Notation = Notation{Kind: NotationDecimal} },
	}
	for i, breakIt := range broken {
		c := *g
		breakIt(&c)
		if err := c.Validate(); !errors.Is(err, domainErrors.ErrInvalidGrinder) {
			t.Errorf("case %d: expected ErrInvalidGrinder, got %v", i, err)
		}
	}
}
