package domain

import "fmt"

type ErrInvalidRange struct {
	mn int
	mx int
}

func NewErrInvalidRange(mn, mx int) ErrInvalidRange {
	return ErrInvalidRange{
		mn: mn,
		mx: mx,
	}
}

func (e ErrInvalidRange) Error() string {
	return fmt.Sprintf("range [%d, %d] is invalid", e.mn, e.mx)
}
