package domain

type CellType int

const (
	Wall CellType = iota
	Passage
	Guessing
)

func (c CellType) String() string {
	switch c {
	case Wall:
		// ANSI код для черного фона
		return "\x1b[40m  \x1b[0m"
	case Passage:
		// ANSI код для белого фона
		return "\x1b[47m  \x1b[0m"
	case Guessing:
		// ANSI код для серого фона
		return "\x1b[100m  \x1b[0m"
	}

	return ""
}

type CellRenderData struct {
	X        int
	Y        int
	Tpe      CellType
	SenderID int
	MCS      int
}

func NewCellRenderData(x, y int, tpe CellType, senderID, ms int) CellRenderData {
	return CellRenderData{
		X:        x,
		Y:        y,
		Tpe:      tpe,
		SenderID: senderID,
		MCS:      ms,
	}
}
