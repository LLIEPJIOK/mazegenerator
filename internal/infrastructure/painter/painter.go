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

func (p *Painter) MoveCursor(rowID, colID int) {
	fmt.Fprintf(p.out, "\033[%d;%dH", rowID+2, 2*(colID+1)+1)
}

func (p *Painter) paint(rowID, colID int, cellType domain.CellType) {
	p.MoveCursor(rowID, colID)
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
				cellData.RowID,
				cellData.ColID,
				p.drawMaze.GetCellType(cellData.RowID, cellData.ColID),
			)
			time.Sleep(time.Duration(cellData.MCS) * time.Microsecond)
		case <-ctx.Done():
			return
		}
	}
}

func (p *Painter) PaintPath(path []domain.Coord, delay time.Duration) {
	for _, v := range path {
		p.paint(v.RowID, v.ColID, domain.Path)

		time.Sleep(delay)
	}
}
