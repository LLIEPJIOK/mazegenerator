package presentation

import "fmt"

type rangePoint struct {
	value   int
	isValid bool
}

func newRangePoint(val int, isValid bool) rangePoint {
	return rangePoint{
		value:   val,
		isValid: isValid,
	}
}

type rangeNumber struct {
	mn rangePoint
	mx rangePoint
}

func newRange(mn, mx rangePoint) (rangeNumber, error) {
	if mn.isValid && mx.isValid && mn.value > mx.value {
		return rangeNumber{}, NewErrInvalidRange(mn.value, mx.value)
	}

	return rangeNumber{
		mn: mn,
		mx: mx,
	}, nil
}

func (r rangeNumber) String() string {
	switch {
	case !r.mn.isValid && !r.mx.isValid:
		return "(-inf, inf)"
	case !r.mn.isValid:
		return fmt.Sprintf("(-inf, %d]", r.mx.value)
	case !r.mx.isValid:
		return fmt.Sprintf("[%d, inf)", r.mn.value)
	default:
		return fmt.Sprintf("[%d, %d]", r.mn.value, r.mx.value)
	}
}

func (r *rangeNumber) Contains(point int) bool {
	switch {
	case !r.mn.isValid && !r.mx.isValid:
		return true
	case !r.mn.isValid:
		return point <= r.mx.value
	case !r.mx.isValid:
		return r.mn.value <= point
	default:
		return r.mn.value <= point && point <= r.mx.value
	}
}
