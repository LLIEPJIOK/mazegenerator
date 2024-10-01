package painter

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
)

type drawingMaze struct {
	cells [][]map[int]domain.CellType
}

func (dm *drawingMaze) getCellType(x, y int) domain.CellType {
	switch len(dm.cells[x][y]) {
	case 0:
		return domain.Wall
	case 1:
		for _, v := range dm.cells[x][y] {
			return v
		}
	}

	return domain.Guessing
}

func (dm *drawingMaze) addCellType(cellData domain.CellRenderData) {
	if cellData.Tpe == domain.Wall {
		delete(dm.cells[cellData.X][cellData.Y], cellData.SenderID)
	} else {
		if cellData.SenderID == 0 {
			clear(dm.cells[cellData.X][cellData.Y])
		}

		dm.cells[cellData.X][cellData.Y][cellData.SenderID] = cellData.Tpe
	}
}

type Painter struct {
	height   int
	width    int
	out      io.Writer
	drawMaze drawingMaze
}

func New(height, width int, out io.Writer) *Painter {
	cells := make([][]map[int]domain.CellType, height)

	for i := range height {
		cells[i] = make([]map[int]domain.CellType, width)
		for j := range width {
			cells[i][j] = make(map[int]domain.CellType)
		}
	}

	return &Painter{
		height: height,
		width:  width,
		out:    out,
		drawMaze: drawingMaze{
			cells: cells,
		},
	}
}

func clearScreen() {
	fmt.Print("\033[2J")
}

func (p *Painter) moveCursor(x, y int) {
	fmt.Fprintf(p.out, "\033[%d;%dH", x+2, 2*(y+1)+1)
}

func (p *Painter) paint(x, y int, cellType domain.CellType) {
	p.moveCursor(x, y)
	cellType.Print(p.out)
}

func (p *Painter) Paint(ctx context.Context, in <-chan domain.CellRenderData) {
	clearScreen()

	defer func() {
		p.moveCursor(p.height-1, p.width+2)
	}()

	for {
		select {
		case cellData, ok := <-in:
			if !ok {
				return
			}

			p.drawMaze.addCellType(cellData)
			p.paint(cellData.X, cellData.Y, p.drawMaze.getCellType(cellData.X, cellData.Y))
			time.Sleep(time.Duration(cellData.MCS) * time.Microsecond)
		case <-ctx.Done():
			return
		}
	}
}
