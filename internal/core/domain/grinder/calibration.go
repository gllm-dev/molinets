package grinder

import (
	"fmt"
	"math"
	"sort"

	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
)

type Segment struct {
	From      float64
	OffsetUM  float64
	UMPerUnit float64
}

type Calibration struct {
	Segments []Segment
}

func NewCalibration(segments []Segment) (Calibration, error) {
	if len(segments) == 0 {
		return Calibration{}, fmt.Errorf("%w: calibration has no segments", domainErrors.ErrInvalidGrinder)
	}
	s := append([]Segment(nil), segments...)
	sort.Slice(s, func(i, j int) bool { return s[i].From < s[j].From })
	for i, seg := range s {
		if seg.UMPerUnit <= 0 {
			return Calibration{}, fmt.Errorf("%w: calibration slope must be positive", domainErrors.ErrInvalidGrinder)
		}
		if i > 0 {
			prev := s[i-1]
			if end := prev.OffsetUM + prev.UMPerUnit*(seg.From-prev.From); seg.OffsetUM < end-1e-9 {
				return Calibration{}, fmt.Errorf("%w: calibration must not go backwards at %v", domainErrors.ErrInvalidGrinder, seg.From)
			}
		}
	}
	return Calibration{Segments: s}, nil
}

func (c Calibration) segmentForSetting(x float64) Segment {
	seg := c.Segments[0]
	for _, s := range c.Segments[1:] {
		if x >= s.From {
			seg = s
		}
	}
	return seg
}

func (c Calibration) segmentForMicrons(um float64) Segment {
	seg := c.Segments[0]
	for _, s := range c.Segments[1:] {
		if um >= s.OffsetUM {
			seg = s
		}
	}
	return seg
}

func (c Calibration) Microns(x float64) float64 {
	s := c.segmentForSetting(x)
	return math.Max(0, s.OffsetUM+s.UMPerUnit*(x-s.From))
}

func (c Calibration) Setting(um float64) float64 {
	s := c.segmentForMicrons(um)
	return s.From + (um-s.OffsetUM)/s.UMPerUnit
}

func (c Calibration) SlopeAt(x float64) float64 {
	return c.segmentForSetting(x).UMPerUnit
}
