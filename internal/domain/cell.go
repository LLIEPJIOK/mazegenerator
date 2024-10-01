package domain

import (
	"io"

	"github.com/fatih/color"
)

type CellType int

const (
	Wall CellType = iota
	Passage
	Guessing
)

func (c CellType) Print(writer io.Writer) {
	var clr *color.Color

	switch c {
	case Wall:
		clr = color.New(color.BgBlack)
	case Passage:
		clr = color.New(color.BgWhite)
	case Guessing:
		clr = color.New(color.BgHiBlack)
	}

	clr.Fprint(writer, "  ")
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
