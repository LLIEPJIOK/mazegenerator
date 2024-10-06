package painter

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
)

type Painter struct {
	out      io.Writer
	drawMaze domain.DrawingMaze
}

func New(out io.Writer) *Painter {
	return &Painter{
		out: out,
	}
}

func ClearScreen() {
	fmt.Print("\033[2J")
}

func (p *Painter) MoveCursor(row, col int) {
	fmt.Fprintf(p.out, "\033[%d;%dH", row+2, 2*(col+1)+1)
}

func (p *Painter) paint(row, col int, cellType domain.CellType) {
	p.MoveCursor(row, col)
	fmt.Fprint(p.out, cellType)
}

func (p *Painter) PaintGeneration(
	ctx context.Context,
	mazeHeight, mazeWidth int,
	cellChan <-chan domain.CellRenderData,
) {
	ClearScreen()

	p.drawMaze = domain.NewDrawingMaze(mazeHeight, mazeWidth)

	for {
		select {
		case cellData, ok := <-cellChan:
			if !ok {
				return
			}

			p.drawMaze.AddCellType(cellData)
			p.paint(
				cellData.Row,
				cellData.Col,
				p.drawMaze.GetCellType(cellData.Row, cellData.Col),
			)
			time.Sleep(cellData.Delay)
		case <-ctx.Done():
			return
		}
	}
}

func (p *Painter) PaintPath(path []domain.Coord, delay time.Duration) {
	for _, v := range path {
		p.paint(v.Row, v.Col, domain.Path)

		time.Sleep(delay)
	}
}
