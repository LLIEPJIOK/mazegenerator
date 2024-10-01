package domain

import (
	"fmt"
	"io"
)

type Maze struct {
	Height int
	Width  int
	Cells  [][]CellType
}

func NewMaze(height, width int, cells [][]CellType) *Maze {
	return &Maze{
		Height: height,
		Width:  width,
		Cells:  cells,
	}
}

func (m *Maze) Print(writer io.Writer) {
	for i := range m.Height {
		for j := range m.Width {
			m.Cells[i][j].Print(writer)
		}

		fmt.Fprintln(writer)
	}
}
