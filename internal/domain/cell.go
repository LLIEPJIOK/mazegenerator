package domain

type CellType int

const (
	Wall CellType = iota
	Passage
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
		return "\x1b[100m  \x1b[0m"
	case Path:
		// ANSI code for red background
		return "\x1b[41m  \x1b[0m"
	}

	return ""
}

type CellRenderData struct {
	RowID    int
	ColID    int
	Tpe      CellType
	SenderID int
	MCS      int
}

func NewCellRenderData(x, y int, tpe CellType, senderID, ms int) CellRenderData {
	return CellRenderData{
		RowID:    x,
		ColID:    y,
		Tpe:      tpe,
		SenderID: senderID,
		MCS:      ms,
	}
}
