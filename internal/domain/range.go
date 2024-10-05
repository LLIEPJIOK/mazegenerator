package domain

import "fmt"

type RangePoint struct {
	Value   int
	IsValid bool
}

func NewRangePoint(val int, isValid bool) RangePoint {
	return RangePoint{
		Value:   val,
		IsValid: isValid,
	}
}

type Range struct {
	mn RangePoint
	mx RangePoint
}

func NewRange(mn, mx RangePoint) (Range, error) {
	if mn.IsValid && mx.IsValid && mn.Value > mx.Value {
		return Range{}, NewErrInvalidRange(mn.Value, mx.Value)
	}

	return Range{
		mn: mn,
		mx: mx,
	}, nil
}

func (r Range) String() string {
	switch {
	case !r.mn.IsValid && !r.mx.IsValid:
		return "(-inf, inf)"
	case !r.mn.IsValid:
		return fmt.Sprintf("(-inf, %d]", r.mx.Value)
	case !r.mx.IsValid:
		return fmt.Sprintf("[%d, inf)", r.mn.Value)
	default:
		return fmt.Sprintf("[%d, %d]", r.mn.Value, r.mx.Value)
	}
}

func (r *Range) Contains(point int) bool {
	switch {
	case !r.mn.IsValid && !r.mx.IsValid:
		return true
	case !r.mn.IsValid:
		return point <= r.mx.Value
	case !r.mx.IsValid:
		return r.mn.Value <= point
	default:
		return r.mn.Value <= point && point <= r.mx.Value
	}
}
