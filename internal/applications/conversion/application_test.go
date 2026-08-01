package conversion

import (
	"context"
	"errors"
	"testing"

	grinderApp "github.com/gllm-dev/molinets/internal/applications/grinder"
	"github.com/gllm-dev/molinets/internal/core/domain/brewmethod"
	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
	"github.com/gllm-dev/molinets/internal/core/domain/grinder"
	"github.com/gllm-dev/molinets/internal/core/ports/repositories"
	"github.com/gllm-dev/molinets/internal/mocks"
)

var ctx = context.Background()

func fixtures(t *testing.T) map[string]grinder.Grinder {
	t.Helper()
	mk := func(id string, offset, slope, lo, hi float64) grinder.Grinder {
		c, err := grinder.NewCalibration([]grinder.Segment{{OffsetUM: offset, UMPerUnit: slope}})
		if err != nil {
			t.Fatal(err)
		}
		return grinder.Grinder{ID: id, Name: id, Type: grinder.TypeManual, Notation: grinder.Notation{Kind: grinder.NotationSteps},
			Calibration: c, Range: grinder.Range{MinUM: lo, MaxUM: hi}}
	}
	return map[string]grinder.Grinder{
		"fine":   mk("fine", -10, 20, 0, 1000),
		"coarse": mk("coarse", 250, 50, 400, 1200),
	}
}

var methods = []brewmethod.Method{
	{ID: "espresso", Name: "Espresso", MinUM: 180, MaxUM: 380},
	{ID: "v60", Name: "V60", MinUM: 400, MaxUM: 700},
	{ID: "cold_brew", Name: "Cold brew", MinUM: 800, MaxUM: 1400},
}

func newApp(t *testing.T) *Application {
	gs := fixtures(t)
	repo := &mocks.GrindersRepositoryMock{
		SearchFunc: func(_ context.Context, q string) ([]grinder.Grinder, error) {
			if g, ok := gs[q]; ok {
				return []grinder.Grinder{g}, nil
			}
			return []grinder.Grinder{}, nil
		},
		ListFunc: func(context.Context, *repositories.GrinderFilter) ([]grinder.Grinder, error) {
			return []grinder.Grinder{gs["fine"], gs["coarse"]}, nil
		},
	}
	methodsRepo := &mocks.MethodsRepositoryMock{
		ListFunc: func(context.Context) ([]brewmethod.Method, error) { return methods, nil },
		FindFunc: func(_ context.Context, q string) (*brewmethod.Method, error) {
			for i := range methods {
				if methods[i].ID == q {
					return &methods[i], nil
				}
			}
			return nil, domainErrors.ErrMethodNotFound
		},
	}
	return New(grinderApp.New(repo), methodsRepo)
}

func TestConvertToEveryOtherGrinder(t *testing.T) {
	res, err := newApp(t).Convert(ctx, ConvertCommand{FromGrinder: "fine", Setting: "28"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Microns != 550 {
		t.Fatalf("fine 28 = %v µm, want 550", res.Microns)
	}
	if len(res.Equivalents) != 1 || res.Equivalents[0].Grinder.ID != "coarse" {
		t.Fatalf("equivalents should be every grinder but the source, got %d", len(res.Equivalents))
	}
	if e := res.Equivalents[0]; !e.Reachable || e.Setting != "6" {
		t.Fatalf("coarse = %+v, want setting 6", e)
	}
	if len(res.MatchingMethods) != 1 || res.MatchingMethods[0] != "V60" {
		t.Fatalf("methods = %v", res.MatchingMethods)
	}
}

func TestConvertMarksUnreachableTargets(t *testing.T) {
	res, err := newApp(t).Convert(ctx, ConvertCommand{FromGrinder: "fine", Setting: "10", ToGrinders: []string{"coarse"}})
	if err != nil {
		t.Fatal(err)
	}
	if e := res.Equivalents[0]; e.Reachable || !e.TooFine || e.Setting != "" {
		t.Fatalf("190 µm on coarse should be too fine, got %+v", e)
	}
}

func TestConvertAtBurrContact(t *testing.T) {
	res, err := newApp(t).Convert(ctx, ConvertCommand{FromGrinder: "fine", Setting: "0", ToGrinders: []string{"coarse"}})
	if err != nil {
		t.Fatal(err)
	}
	if !res.AtBurrContact || res.Microns != 0 || res.MatchingMethods == nil {
		t.Fatalf("got %+v", res)
	}
}

func TestConvertErrors(t *testing.T) {
	app := newApp(t)
	cases := []struct {
		cmd  ConvertCommand
		want error
	}{
		{ConvertCommand{FromGrinder: "nope", Setting: "1"}, domainErrors.ErrGrinderNotFound},
		{ConvertCommand{FromGrinder: "fine", Setting: "abc"}, domainErrors.ErrInvalidSetting},
		{ConvertCommand{FromGrinder: "fine", Setting: "90"}, domainErrors.ErrSettingOutOfRange},
		{ConvertCommand{FromGrinder: "fine", Setting: "10", ToGrinders: []string{"nope"}}, domainErrors.ErrGrinderNotFound},
	}
	for _, c := range cases {
		if _, err := app.Convert(ctx, c.cmd); !errors.Is(err, c.want) {
			t.Errorf("%+v: got %v, want %v", c.cmd, err, c.want)
		}
	}
}

func TestForMethod(t *testing.T) {
	app := newApp(t)
	r, err := app.ForMethod(ctx, "fine", "v60")
	if err != nil {
		t.Fatal(err)
	}
	if r.From != "21" || r.To != "36" || r.Partial {
		t.Fatalf("fine V60 = %+v, want 21–36", r)
	}
	r, err = app.ForMethod(ctx, "fine", "cold_brew")
	if err != nil || !r.Partial || r.MaxUM != 1000 {
		t.Fatalf("fine cold brew should be partial up to 1000 µm: %+v %v", r, err)
	}
	if _, err := app.ForMethod(ctx, "coarse", "espresso"); !errors.Is(err, domainErrors.ErrMethodUnreachable) {
		t.Fatalf("coarse espresso: %v", err)
	}
	if _, err := app.ForMethod(ctx, "fine", "chemex"); !errors.Is(err, domainErrors.ErrMethodNotFound) {
		t.Fatalf("unknown method: %v", err)
	}
}

func TestFromMicrons(t *testing.T) {
	app := newApp(t)
	for _, um := range []float64{0, -5, 2001} {
		if _, err := app.FromMicrons(ctx, "fine", um); !errors.Is(err, domainErrors.ErrInvalidMicrons) {
			t.Errorf("%v µm: %v", um, err)
		}
	}
	e, err := app.FromMicrons(ctx, "fine", 20)
	if err != nil || e.Reachable {
		t.Fatalf("20 µm should be unreachable: %+v %v", e, err)
	}
	e, err = app.FromMicrons(ctx, "coarse", 550)
	if err != nil || e.Setting != "6" {
		t.Fatalf("coarse 550 µm = %+v %v", e, err)
	}
}
