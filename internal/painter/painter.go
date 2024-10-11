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
	data      domain.MazeData
	paintMaze PaintingMaze
}

func New(out io.Writer, data domain.MazeData) *Painter {
	return &Painter{
		out:       out,
		data:      data,
		paintMaze: newPaintingMaze(data.Height, data.Width),
	}
}

func clearScreen() {
	fmt.Print("\033[2J")
}

func (p *Painter) moveCursor(row, col int) {
	fmt.Fprintf(p.out, "\033[%d;%dH", row+2, 2*(col+1)+1)
}

func (p *Painter) paintStartEndString(row, col int, str string) {
	p.moveCursor(row, col)

	// ANSI code for red letters
	fmt.Fprintf(p.out, "\033[31m%s\033[0m", str)
}

func (p *Painter) paintStartEnd() {
	switch {
	case p.data.Start.Row == 0:
		p.paintStartEndString(-1, p.data.Start.Col, "vv")

	case p.data.Start.Row == p.data.Height-1:
		p.paintStartEndString(p.data.Height, p.data.Start.Col, "^^")

	case p.data.Start.Col == 0:
		p.paintStartEndString(p.data.Start.Row, -1, ">")

	case p.data.Start.Col == p.data.Width-1:
		p.paintStartEndString(p.data.Start.Row, p.data.Width, "<")
	}

	switch {
	case p.data.End.Row == 0:
		p.paintStartEndString(-1, p.data.End.Col, "^^")

	case p.data.End.Row == p.data.Height-1:
		p.paintStartEndString(p.data.Height, p.data.End.Col, "vv")

	case p.data.End.Col == 0:
		p.paintStartEndString(p.data.End.Row, -1, "<")

	case p.data.End.Col == p.data.Width-1:
		p.paintStartEndString(p.data.End.Row, p.data.Width, ">")
	}
}

func (p *Painter) paint(row, col int, cellType domain.CellType) {
	p.moveCursor(row, col)
	fmt.Fprint(p.out, cellType)
}

func (p *Painter) PaintGeneration(
	ctx context.Context,
	cellChan <-chan domain.CellPaintingData,
) {
	defer p.moveCursor(p.data.Height+1, -1)

	clearScreen()
	p.paintStartEnd()

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

func (p *Painter) PaintPath(pathChan <-chan []domain.Coord, delay time.Duration) {
	defer p.moveCursor(p.data.Height+1, -1)

	var prevPath []domain.Coord

	for path := range pathChan {
		for _, v := range prevPath {
			p.paint(
				v.Row,
				v.Col,
				p.paintMaze.GetCellType(v.Row, v.Col),
			)
		}

		for _, v := range path {
			p.paint(v.Row, v.Col, domain.Path)
		}

		prevPath = path

		time.Sleep(delay)
	}
}
