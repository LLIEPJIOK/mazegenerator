package painter

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
)

type Painter struct {
	out       io.Writer
	height    int
	width     int
	paintMaze PaintingMaze
}

func New(out io.Writer, height, width int) *Painter {
	return &Painter{
		out:       out,
		height:    height,
		width:     width,
		paintMaze: newPaintingMaze(height, width),
	}
}

func clearScreen() {
	fmt.Print("\033[2J")
}

func (p *Painter) moveCursor(row, col int) {
	fmt.Fprintf(p.out, "\033[%d;%dH", row+2, 2*(col+1)+1)
}

func (p *Painter) paint(row, col int, cellType domain.CellType) {
	p.moveCursor(row, col)
	fmt.Fprint(p.out, cellType)
}

func (p *Painter) PaintGeneration(
	ctx context.Context,
	cellChan <-chan domain.PaintingData,
) {
	defer p.moveCursor(p.height+1, -1)

	clearScreen()

	for {
		select {
		case cellData, ok := <-cellChan:
			if !ok {
				return
			}

			p.paintMaze.AddCellType(cellData)
			p.paint(
				cellData.Row,
				cellData.Col,
				p.paintMaze.GetCellType(cellData.Row, cellData.Col),
			)
			time.Sleep(cellData.Delay)
		case <-ctx.Done():
			return
		}
	}
}

func (p *Painter) PaintPath(path []domain.Coord, delay time.Duration) {
	defer p.moveCursor(p.height+1, -1)

	for _, v := range path {
		p.paint(v.Row, v.Col, domain.Path)

		time.Sleep(delay)
	}
}
