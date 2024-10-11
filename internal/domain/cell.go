package domain

import "time"

type CellType int

const (
	Wall CellType = iota
	Passage
	Money
	Sand
	River
	Guessing
	Path
)

func (c CellType) String() string {
	switch c {
	case Wall:
		// ANSI code for black background
		return "\x1b[40m  \x1b[0m"
	case Passage:
		// ANSI code for white background
		return "\x1b[47m  \x1b[0m"
	case Guessing:
		// ANSI code for gray background
		return "\x1b[100m??\x1b[0m"
	case Money:
		// ANSI code for yellow background and green symbols
		return "\033[0;103m\033[32m$$\033[0m"
	case Sand:
		// ANSI green background
		return "\033[43m▒▒\033[0m"
	case River:
		// ANSI blue background
		return "\033[44m~~\033[0m"
	case Path:
		// ANSI code for red background
		return "\x1b[41m  \x1b[0m"
	}

	return ""
}

func (c CellType) IsTraversable() bool {
	return c != Wall && c != Guessing && c != Path
}

func (c CellType) Cost() int {
	switch c {
	case Passage:
		return 3

	case Money:
		return 1

	case Sand:
		return 5

	case River:
		return 10

	case Wall, Guessing, Path:
		return 0
	}

	return 0
}

type PaintingData struct {
	Row      int
	Col      int
	Tpe      CellType
	SenderID int
	Delay    time.Duration
}

func NewPaintingData(row, col int, tpe CellType, senderID int, delay time.Duration) PaintingData {
	return PaintingData{
		Row:      row,
		Col:      col,
		Tpe:      tpe,
		SenderID: senderID,
		Delay:    delay,
	}
}
