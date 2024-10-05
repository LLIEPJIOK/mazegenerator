package ui

type ErrNoInputLines struct{}

func (e ErrNoInputLines) Error() string {
	return "no more input lines"
}
