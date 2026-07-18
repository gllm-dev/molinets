package grindersrepo

import (
	"context"
	"errors"
	"testing"

	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
	"github.com/gllm-dev/molinets/internal/core/domain/grinder"
	"github.com/gllm-dev/molinets/internal/core/ports/repositories"
)

var ctx = context.Background()

func newRepo(t *testing.T) *Repository {
	t.Helper()
	r, err := New()
	if err != nil {
		t.Fatalf("embedded dataset does not load: %v", err)
	}
	return r
}

func TestDatasetLoads(t *testing.T) {
	r := newRepo(t)
	all, _ := r.List(ctx, nil)
	if len(all) != 34 {
		t.Fatalf("got %d grinders, want 34", len(all))
	}
	manual, _ := r.List(ctx, &repositories.GrinderFilter{Type: grinder.TypeManual})
	electric, _ := r.List(ctx, &repositories.GrinderFilter{Type: grinder.TypeElectric})
	if len(manual)+len(electric) != len(all) || len(manual) == 0 || len(electric) == 0 {
		t.Fatalf("type split %d + %d != %d", len(manual), len(electric), len(all))
	}
}

func TestRejectsBadDatasets(t *testing.T) {
	cases := map[string]string{
		"not json":     `{`,
		"duplicate id": `[{"id":"a","name":"A","type":"manual","notation":{"kind":"steps"},"calibration":[{"um_per_unit":10}],"range_um":[0,900]},{"id":"a","name":"B","type":"manual","notation":{"kind":"steps"},"calibration":[{"um_per_unit":10}],"range_um":[0,900]}]`,
		"bad notation": `[{"id":"a","name":"A","type":"manual","notation":{"kind":"decimal"},"calibration":[{"um_per_unit":10}],"range_um":[0,900]}]`,
		"no slope":     `[{"id":"a","name":"A","type":"manual","notation":{"kind":"steps"},"calibration":[{"um_per_unit":0}],"range_um":[0,900]}]`,
	}
	for name, raw := range cases {
		if _, err := NewFromJSON([]byte(raw)); !errors.Is(err, domainErrors.ErrInvalidGrinder) {
			t.Errorf("%s: expected ErrInvalidGrinder, got %v", name, err)
		}
	}
}

func TestGet(t *testing.T) {
	r := newRepo(t)
	if g, err := r.Get(ctx, "kingrinder-k6"); err != nil || g.Name != "Kingrinder K4 / K6" {
		t.Fatalf("Get: %v %v", g, err)
	}
	if _, err := r.Get(ctx, "nope"); !errors.Is(err, domainErrors.ErrGrinderNotFound) {
		t.Fatalf("Get unknown: %v", err)
	}
}

func TestSearch(t *testing.T) {
	r := newRepo(t)
	exact := map[string]string{
		"Comandante C40 MK4": "comandante-c40",
		"k6":                 "kingrinder-k6",
		"Chestnut C3S Pro":   "timemore-c3",
		"c3":                 "timemore-c3",
		"encore":             "baratza-encore",
		"red clix":           "comandante-c40-red-clix",
		"ode":                "fellow-ode-gen-2",
		"Virtuoso+":          "baratza-virtuoso-plus",
		"Mahlkönig X54":      "mahlkonig-x54",
		"mahlkoenig x54":     "mahlkonig-x54",
		"niche zero":         "niche-zero",
	}
	for q, want := range exact {
		got, err := r.Search(ctx, q)
		if err != nil || len(got) != 1 || got[0].ID != want {
			t.Errorf("Search(%q) = %v, %v; want %s", q, ids(got), err, want)
		}
	}
	for _, q := range []string{"hario", "timemore", "1zpresso", "baratza", "fellow", "wilfa"} {
		if got, _ := r.Search(ctx, q); len(got) < 2 {
			t.Errorf("Search(%q) should be ambiguous, got %v", q, ids(got))
		}
	}
	for _, q := range []string{"bialetti", "", "  "} {
		if got, err := r.Search(ctx, q); err != nil || got == nil || len(got) != 0 {
			t.Errorf("Search(%q) = %#v, %v; want empty", q, got, err)
		}
	}
}

func ids(gs []grinder.Grinder) []string {
	out := make([]string, len(gs))
	for i := range gs {
		out[i] = gs[i].ID
	}
	return out
}
