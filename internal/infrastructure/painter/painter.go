package painter

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
)

type Painter struct {
	height   int
	width    int
	out      io.Writer
	drawMaze domain.DrawingMaze
}

func New(height, width int, out io.Writer) *Painter {
	return &Painter{
		height:   height,
		width:    width,
		out:      out,
		drawMaze: domain.NewDrawingMaze(height, width),
	}
}

func clearScreen() {
	fmt.Print("\033[2J")
}

func (p *Painter) moveCursor(rowID, colID int) {
	fmt.Fprintf(p.out, "\033[%d;%dH", rowID+2, 2*(colID+1)+1)
}

func (p *Painter) paint(rowID, colID int, cellType domain.CellType) {
	p.moveCursor(rowID, colID)
	fmt.Fprint(p.out, cellType)
}

func (p *Painter) PaintGeneration(ctx context.Context, cellChan <-chan domain.CellRenderData) {
	clearScreen()

	defer func() {
		p.moveCursor(p.height+1, -1)
	}()

	for {
		select {
		case cellData, ok := <-cellChan:
			if !ok {
				return
			}

			p.drawMaze.AddCellType(cellData)
			p.paint(cellData.RowID, cellData.ColID, p.drawMaze.GetCellType(cellData.RowID, cellData.ColID))
			time.Sleep(time.Duration(cellData.MCS) * time.Microsecond)
		case <-ctx.Done():
			return
		}
	}
}

func (p *Painter) PaintPath(path []domain.Coord, delay time.Duration) {
	defer func() {
		p.moveCursor(p.height+1, -1)
	}()

	for _, v := range path {
		p.paint(v.RowID, v.ColID, domain.Path)

		time.Sleep(delay)
	}
}
