package grinder

import (
	"errors"
	"strings"
	"testing"

	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
)

var (
	jmax   = Notation{Kind: NotationRotationNumberClick, ClicksPerRotation: 90, ClicksPerNumber: 10}
	q2     = Notation{Kind: NotationRotationNumberClick, ClicksPerRotation: 30, ClicksPerNumber: 3}
	k6     = Notation{Kind: NotationRotationClicks, ClicksPerRotation: 60}
	zp6    = Notation{Kind: NotationNumberClick, ClicksPerNumber: 10}
	c5     = Notation{Kind: NotationNumberClick, ClicksPerNumber: 3}
	ode    = Notation{Kind: NotationThirds}
	clicks = Notation{Kind: NotationSteps}
	opus   = Notation{Kind: NotationDecimal, Step: 0.25}
	x54    = Notation{Kind: NotationDecimal, Step: 0.5}
	sculpt = Notation{Kind: NotationDecimal, Step: 0.1}
)

func TestNotationRoundTrip(t *testing.T) {
	cases := []struct {
		n Notation
		s string
	}{
		{jmax, "2.8.4"}, {jmax, "0.0.0"}, {q2, "1.2.2"}, {k6, "1.22"}, {k6, "0.05"}, {zp6, "4.5"},
		{c5, "11.2"}, {ode, "4.2"}, {ode, "6"}, {clicks, "26"}, {opus, "5.25"}, {opus, "6"},
		{x54, "12.5"}, {sculpt, "3.3"}, {sculpt, "9.7"},
	}
	for _, c := range cases {
		x, err := c.n.Parse(c.s)
		if err != nil {
			t.Fatalf("%s %q: %v", c.n.Kind, c.s, err)
		}
		if got := c.n.Format(x); got != c.s {
			t.Errorf("%s: %q -> %v -> %q", c.n.Kind, c.s, x, got)
		}
	}
}

func TestNotationParseValues(t *testing.T) {
	cases := []struct {
		n    Notation
		s    string
		want float64
	}{
		{clicks, " 26 ", 26},
		{ode, "4,2", 4 + 2.0/3},
		{k6, "1.5", 110},
		{k6, "1.05", 65},
		{k6, "2", 120},
		{jmax, "2.8", 260},
		{jmax, "3", 270},
		{zp6, "4", 40},
	}
	for _, c := range cases {
		got, err := c.n.Parse(c.s)
		if err != nil || got != c.want {
			t.Errorf("%s %q = %v, %v; want %v", c.n.Kind, c.s, got, err, c.want)
		}
	}
}

func TestNotationRejectsMalformed(t *testing.T) {
	cases := []struct {
		n       Notation
		s, want string
	}{
		{clicks, "abc", "only digits"},
		{clicks, "-3", "only digits"},
		{clicks, "1e3", "only digits"},
		{clicks, "26.5", "whole steps"},
		{clicks, "26 clicks", "only digits"},
		{jmax, "1.9.0", "number must be 0–8"},
		{jmax, "1.2.10", "click must be 0–9"},
		{jmax, "1.2.3.4", "rotation.number.click"},
		{q2, "1.2.3", "click must be 0–2"},
		{k6, "1.60", "clicks must be 0–59"},
		{k6, "1.2.3", "rotation.clicks"},
		{zp6, "4.10", "click must be 0–9"},
		{c5, "7.3", "click must be 0–2"},
		{ode, "4.3", "thirds"},
		{ode, "4..2", "empty part"},
		{opus, "1.2.3", "single decimal"},
	}
	for _, c := range cases {
		_, err := c.n.Parse(c.s)
		if !errors.Is(err, domainErrors.ErrInvalidSetting) || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s %q: got %v, want ErrInvalidSetting mentioning %q", c.n.Kind, c.s, err, c.want)
		}
	}
	if _, err := clicks.Parse("   "); !errors.Is(err, domainErrors.ErrEmptySetting) {
		t.Errorf("blank setting: %v", err)
	}
}

func TestDecimalFormatHasNoFloatNoise(t *testing.T) {
	for _, n := range []Notation{opus, x54, sculpt} {
		for x := 0.0; x < 20; x += 0.07 {
			s := n.Format(x)
			if strings.Contains(s, "0000") || strings.Contains(s, "9999") {
				t.Fatalf("%v formats %v as %q", n.Step, x, s)
			}
		}
	}
}

func TestNewNotationValidates(t *testing.T) {
	bad := []Notation{
		{Kind: "nope"},
		{Kind: NotationDecimal},
		{Kind: NotationDecimal, Step: 1},
		{Kind: NotationRotationNumberClick, ClicksPerRotation: 90, ClicksPerNumber: 7},
		{Kind: NotationRotationClicks},
		{Kind: NotationNumberClick},
	}
	for _, n := range bad {
		if _, err := NewNotation(n); !errors.Is(err, domainErrors.ErrInvalidGrinder) {
			t.Errorf("%+v: expected ErrInvalidGrinder, got %v", n, err)
		}
	}
	for _, n := range []Notation{jmax, q2, k6, zp6, c5, ode, clicks, opus, x54, sculpt} {
		if _, err := NewNotation(n); err != nil {
			t.Errorf("%+v: %v", n, err)
		}
	}
}
