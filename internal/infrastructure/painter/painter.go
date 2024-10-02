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

func (p *Painter) moveCursor(x, y int) {
	fmt.Fprintf(p.out, "\033[%d;%dH", x+2, 2*(y+1)+1)
}

func (p *Painter) paint(x, y int, cellType domain.CellType) {
	p.moveCursor(x, y)
	fmt.Fprint(p.out, cellType)
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

			p.drawMaze.AddCellType(cellData)
			p.paint(cellData.X, cellData.Y, p.drawMaze.GetCellType(cellData.X, cellData.Y))
			time.Sleep(time.Duration(cellData.MCS) * time.Microsecond)
		case <-ctx.Done():
			return
		}
	}
}
