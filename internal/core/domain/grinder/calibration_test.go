package grinder

import (
	"errors"
	"math"
	"testing"

	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
)

func encoreESP(t *testing.T) Calibration {
	t.Helper()
	c, err := NewCalibration([]Segment{{From: 20, OffsetUM: 460, UMPerUnit: 46}, {From: 0, OffsetUM: 230, UMPerUnit: 11.5}})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestPiecewiseCalibration(t *testing.T) {
	c := encoreESP(t)
	if got := c.Microns(20); got != 460 {
		t.Fatalf("step 20 = %v µm, want 460", got)
	}
	if jump := c.Microns(21) - c.Microns(20); jump != 46 {
		t.Fatalf("macro step = %v µm, want 46", jump)
	}
	if micro := c.Microns(5) - c.Microns(4); micro != 11.5 {
		t.Fatalf("micro step = %v µm, want 11.5", micro)
	}
	for _, x := range []float64{0, 5, 19.5, 20, 21, 33, 40} {
		if back := c.Setting(c.Microns(x)); math.Abs(back-x) > 1e-9 {
			t.Errorf("round trip %v -> %v", x, back)
		}
	}
}

func TestMicronsNeverNegative(t *testing.T) {
	c, _ := NewCalibration([]Segment{{OffsetUM: -10, UMPerUnit: 27.6}})
	if got := c.Microns(0); got != 0 {
		t.Fatalf("burr contact = %v µm, want 0", got)
	}
}

func TestNewCalibrationValidates(t *testing.T) {
	cases := [][]Segment{
		nil,
		{{UMPerUnit: 0}},
		{{UMPerUnit: -1}},
		{{From: 0, OffsetUM: 230, UMPerUnit: 11.5}, {From: 20, OffsetUM: 300, UMPerUnit: 46}},
	}
	for _, s := range cases {
		if _, err := NewCalibration(s); !errors.Is(err, domainErrors.ErrInvalidGrinder) {
			t.Errorf("%+v: expected ErrInvalidGrinder, got %v", s, err)
		}
	}
}
