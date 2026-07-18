package grindersrepo

import (
	"encoding/json"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/gllm-dev/molinets/internal/core/domain/grinder"
)

var methodRanges = map[string][2]float64{
	"espresso":     {180, 380},
	"v60":          {400, 700},
	"pour_over":    {410, 930},
	"aeropress":    {320, 960},
	"french_press": {690, 1300},
}

func TestMatchesReferenceTables(t *testing.T) {
	raw, err := os.ReadFile("../testdata/hcg_reference.json")
	if err != nil {
		t.Fatal(err)
	}
	var rows []struct{ Grinder, Method, From, To string }
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatal(err)
	}
	r := newRepo(t)
	tolerance := map[string]float64{"niche-zero": 2.5}
	checked := 0
	for _, row := range rows {
		g, err := r.Get(ctx, row.Grinder)
		if err != nil {
			t.Fatal(err)
		}
		mr := methodRanges[row.Method]
		for _, p := range []struct {
			s  string
			um float64
		}{{row.From, mr[0]}, {row.To, mr[1]}} {
			if p.um <= g.Range.MinUM || p.um >= g.Range.MaxUM {
				continue
			}
			x, err := g.Notation.Parse(p.s)
			if err != nil {
				t.Errorf("%s %s: reference value %q does not parse: %v", g.ID, row.Method, p.s, err)
				continue
			}
			off := math.Abs(g.FromMicrons(p.um)-x) / g.Notation.StepSize()
			tol := 1.5
			if v, ok := tolerance[g.ID]; ok {
				tol = v
			}
			if off > tol {
				t.Errorf("%s %s: reference %s at %.0f µm, model %s (%.1f steps off)",
					g.ID, row.Method, p.s, p.um, g.FormatSetting(g.FromMicrons(p.um)), off)
			}
			checked++
		}
	}
	if checked < 250 {
		t.Fatalf("only %d reference points checked", checked)
	}
	t.Logf("%d reference points within tolerance", checked)
}

func TestMatchesBeean(t *testing.T) {
	r := newRepo(t)
	c40, _ := r.Get(ctx, "comandante-c40")
	c3, _ := r.Get(ctx, "timemore-c3")
	if x := c3.FromMicrons(c40.ToMicrons(26)); x < 17.5 || x > 19.5 {
		t.Fatalf("C40 26 -> C3S Pro %.2f, want 18–19", x)
	}
}

func each(t *testing.T, fn func(g *grinder.Grinder)) {
	all, _ := newRepo(t).List(ctx, nil)
	for i := range all {
		fn(&all[i])
	}
}

func TestEveryGrinderFormatsCleanly(t *testing.T) {
	each(t, func(g *grinder.Grinder) {
		for um := g.FinestUM(); um <= g.CoarsestUM(); um += 7 {
			s := g.FormatSetting(g.FromMicrons(um))
			if strings.ContainsAny(s, "-NI") || strings.Contains(s, "0000") || strings.Contains(s, "9999") {
				t.Fatalf("%s at %v µm formats as %q", g.ID, um, s)
			}
			x, err := g.Notation.Parse(s)
			if err != nil {
				t.Fatalf("%s: formatted %q does not parse back: %v", g.ID, s, err)
			}
			if g.FormatSetting(x) != s {
				t.Fatalf("%s: %q is not stable", g.ID, s)
			}
		}
	})
}

func TestEveryGrinderIsMonotonic(t *testing.T) {
	each(t, func(g *grinder.Grinder) {
		prev := -1.0
		for x := g.MinSetting(); x <= g.MaxSetting(); x += g.Notation.StepSize() {
			um := g.ToMicrons(x)
			if um < prev || (um == prev && um > 0) {
				t.Fatalf("%s is not monotonic at %v", g.ID, x)
			}
			prev = um
		}
	})
}

func TestNoGrinderGoesNegative(t *testing.T) {
	each(t, func(g *grinder.Grinder) {
		if um := g.ToMicrons(0); um < 0 {
			t.Errorf("%s at 0 = %v µm", g.ID, um)
		}
	})
}

func TestPairwiseRoundTrip(t *testing.T) {
	all, _ := newRepo(t).List(ctx, nil)
	for i := range all {
		a := &all[i]
		for j := range all {
			b := &all[j]
			for um := 250.0; um <= 1000; um += 50 {
				if !a.CanGrind(um) || !b.CanGrind(um) {
					continue
				}
				xa := a.Notation.Snap(a.FromMicrons(um))
				back := a.FromMicrons(b.ToMicrons(b.FromMicrons(a.ToMicrons(xa))))
				if math.Abs(back-xa) > 1e-6 {
					t.Fatalf("%s -> %s -> %s drifts by %v", a.ID, b.ID, a.ID, back-xa)
				}
			}
		}
	}
}
