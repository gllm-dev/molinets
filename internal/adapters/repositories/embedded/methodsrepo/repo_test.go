package methodsrepo

import (
	"context"
	"errors"
	"strings"
	"testing"

	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
)

func TestFind(t *testing.T) {
	r, err := New()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	cases := map[string]string{
		"espresso": "espresso", "V60": "v60", "french press": "french_press", "French Press": "french_press",
		"chemex": "pour_over", "Origami": "pour_over", "Hario Switch": "steep_release", "clever": "steep_release",
		"cafetera italiana": "moka", "premsa francesa": "french_press", "Moccamaster": "filter_machine",
		"cold brew": "cold_brew", "aero": "aeropress",
	}
	for q, want := range cases {
		m, err := r.Find(ctx, q)
		if err != nil || m.ID != want {
			t.Errorf("Find(%q) = %v, %v; want %s", q, m, err, want)
		}
	}
	for _, q := range []string{"", "turkish delight", "syphon"} {
		_, err := r.Find(ctx, q)
		if !errors.Is(err, domainErrors.ErrMethodNotFound) || !strings.Contains(err.Error(), "use one of") {
			t.Errorf("Find(%q): %v", q, err)
		}
	}
	all, _ := r.List(ctx)
	if len(all) != 11 {
		t.Fatalf("got %d methods, want 11", len(all))
	}
}
