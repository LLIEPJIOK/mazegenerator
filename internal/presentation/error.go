package presentation

import "fmt"

type ErrNoInputLines struct{}

func (e ErrNoInputLines) Error() string {
	return "no more input lines"
}

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
