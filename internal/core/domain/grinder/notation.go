package grinder

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
)

type NotationKind string

const (
	NotationSteps               NotationKind = "steps"
	NotationDecimal             NotationKind = "decimal"
	NotationRotationNumberClick NotationKind = "rotation-number-click"
	NotationRotationClicks      NotationKind = "rotation-clicks"
	NotationNumberClick         NotationKind = "number-click"
	NotationThirds              NotationKind = "thirds"
)

type Notation struct {
	Kind              NotationKind
	ClicksPerRotation int
	ClicksPerNumber   int
	Step              float64
}

func NewNotation(n Notation) (Notation, error) {
	bad := func(why string) (Notation, error) {
		return Notation{}, fmt.Errorf("%w: %s notation %s", domainErrors.ErrInvalidGrinder, n.Kind, why)
	}
	switch n.Kind {
	case NotationSteps, NotationThirds:
	case NotationDecimal:
		if n.Step <= 0 || n.Step >= 1 {
			return bad("needs a step between 0 and 1")
		}
	case NotationRotationNumberClick:
		if n.ClicksPerRotation <= 0 || n.ClicksPerNumber <= 0 || n.ClicksPerRotation%n.ClicksPerNumber != 0 {
			return bad("needs clicks per rotation divisible by clicks per number")
		}
	case NotationRotationClicks:
		if n.ClicksPerRotation <= 0 || n.ClicksPerRotation > 100 {
			return bad("needs between 1 and 100 clicks per rotation")
		}
	case NotationNumberClick:
		if n.ClicksPerNumber <= 0 {
			return bad("needs clicks per number")
		}
	default:
		return bad("is unknown")
	}
	return n, nil
}

func (n Notation) StepSize() float64 {
	switch n.Kind {
	case NotationDecimal:
		return n.Step
	case NotationThirds:
		return 1.0 / 3
	default:
		return 1
	}
}

func (n Notation) Snap(x float64) float64 {
	s := n.StepSize()
	return math.Round(x/s) * s
}

func (n Notation) Format(x float64) string {
	x = n.Snap(x)
	xi := int(math.Round(x))
	switch n.Kind {
	case NotationRotationNumberClick:
		r, rem := xi/n.ClicksPerRotation, xi%n.ClicksPerRotation
		return fmt.Sprintf("%d.%d.%d", r, rem/n.ClicksPerNumber, rem%n.ClicksPerNumber)
	case NotationRotationClicks:
		return fmt.Sprintf("%d.%02d", xi/n.ClicksPerRotation, xi%n.ClicksPerRotation)
	case NotationNumberClick:
		return fmt.Sprintf("%d.%d", xi/n.ClicksPerNumber, xi%n.ClicksPerNumber)
	case NotationThirds:
		t := int(math.Round(x * 3))
		if t%3 == 0 {
			return strconv.Itoa(t / 3)
		}
		return fmt.Sprintf("%d.%d", t/3, t%3)
	case NotationDecimal:
		decimals := len(strconv.FormatFloat(n.Step, 'f', -1, 64)) - 2
		out := strconv.FormatFloat(x, 'f', decimals, 64)
		if strings.Contains(out, ".") {
			out = strings.TrimRight(strings.TrimRight(out, "0"), ".")
		}
		return out
	default:
		return strconv.Itoa(xi)
	}
}

func (n Notation) Parse(raw string) (float64, error) {
	s := strings.TrimSpace(strings.ReplaceAll(raw, ",", "."))
	if s == "" {
		return 0, domainErrors.ErrEmptySetting
	}
	bad := func(why string) (float64, error) {
		return 0, fmt.Errorf("%w %q: %s", domainErrors.ErrInvalidSetting, raw, why)
	}
	parts := strings.Split(s, ".")
	nums := make([]int, len(parts))
	for i, p := range parts {
		if p == "" {
			return bad("empty part")
		}
		for _, r := range p {
			if r < '0' || r > '9' {
				return bad("only digits and dots are allowed")
			}
		}
		v, err := strconv.Atoi(p)
		if err != nil {
			return bad(err.Error())
		}
		nums[i] = v
	}
	at := func(i int) int {
		if i < len(nums) {
			return nums[i]
		}
		return 0
	}

	switch n.Kind {
	case NotationRotationNumberClick:
		if len(parts) > 3 {
			return bad("use rotation.number.click, e.g. 2.8.4")
		}
		numbers := n.ClicksPerRotation / n.ClicksPerNumber
		if at(1) >= numbers {
			return bad(fmt.Sprintf("the number must be 0–%d", numbers-1))
		}
		if at(2) >= n.ClicksPerNumber {
			return bad(fmt.Sprintf("the click must be 0–%d", n.ClicksPerNumber-1))
		}
		return float64(at(0)*n.ClicksPerRotation + at(1)*n.ClicksPerNumber + at(2)), nil

	case NotationRotationClicks:
		if len(parts) > 2 {
			return bad("use rotation.clicks, e.g. 1.22")
		}
		clicks := 0
		if len(parts) == 2 {
			c := parts[1]
			if len(c) == 1 {
				c += "0"
			}
			clicks, _ = strconv.Atoi(c)
		}
		if clicks >= n.ClicksPerRotation {
			return bad(fmt.Sprintf("clicks must be 0–%d", n.ClicksPerRotation-1))
		}
		return float64(at(0)*n.ClicksPerRotation + clicks), nil

	case NotationNumberClick:
		if len(parts) > 2 {
			return bad("use number.click, e.g. 4.5")
		}
		if at(1) >= n.ClicksPerNumber {
			return bad(fmt.Sprintf("the click must be 0–%d", n.ClicksPerNumber-1))
		}
		return float64(at(0)*n.ClicksPerNumber + at(1)), nil

	case NotationThirds:
		if len(parts) > 2 || at(1) > 2 {
			return bad("use number.thirds, e.g. 4.2 for 4⅔")
		}
		return float64(at(0)) + float64(at(1))/3, nil

	case NotationDecimal:
		if len(parts) > 2 {
			return bad("use a single decimal number")
		}
		v, _ := strconv.ParseFloat(s, 64)
		return v, nil

	default:
		if len(parts) > 1 {
			return bad("this grinder only has whole steps")
		}
		return float64(at(0)), nil
	}
}
