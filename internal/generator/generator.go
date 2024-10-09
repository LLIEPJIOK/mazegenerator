package generator

import (
	"fmt"
	"sync"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
	"golang.org/x/sync/errgroup"
)

const (
	drawingDelay = 3 * time.Millisecond
	clearDelay   = 100 * time.Microsecond
	mergeDelay   = 30 * time.Microsecond
)

type cell struct {
	Row   int
	Col   int
	Tpe   domain.CellType
	Delay time.Duration
}

func newCell(row, col int, tpe domain.CellType, delay time.Duration) cell {
	return cell{
		Row:   row,
		Col:   col,
		Tpe:   tpe,
		Delay: delay,
	}
}

type Algorithm interface {
	createMazeCellsFromCoord(
		height, width int,
		start domain.Coord,
		drawingChan chan<- cell,
	) ([][]domain.CellType, error)
}

type Generator struct {
	algo   Algorithm
	dirRow []int
	dirCol []int
}

func New(algo Algorithm) *Generator {
	return &Generator{
		algo:   algo,
		dirRow: []int{-1, 1, 0, 0},
		dirCol: []int{0, 0, -1, 1},
	}
}

func (g *Generator) clearDeadEnd(
	cells [][]domain.CellType,
	height, width int,
	start domain.Coord,
	drawingChan chan<- cell,
) [][]domain.CellType {
	newCells := make([][]domain.CellType, height)

	for i := range height {
		newCells[i] = make([]domain.CellType, width)

		for j := range width {
			cntPassages := 0

			for k := range g.dirRow {
				newRowID, newColID := i+g.dirRow[k], j+g.dirCol[k]
				if min(newRowID, newColID) < 0 || newRowID >= height || newColID >= width {
					continue
				}

				if cells[newRowID][newColID] == domain.Passage {
					cntPassages++
				}
			}

			if cntPassages == 1 && (i != start.Row || j != start.Col) {
				newCells[i][j] = domain.Wall
				drawingChan <- newCell(i, j, domain.Wall, clearDelay)
			} else {
				newCells[i][j] = cells[i][j]
			}
		}
	}

	return newCells
}

const clearDeadEndNumber = 10

func (g *Generator) generateMazeCellsFromCoord(
	height, width int,
	start domain.Coord,
	drawingChan chan<- cell,
) ([][]domain.CellType, error) {
	cells, err := g.algo.createMazeCellsFromCoord(height, width, start, drawingChan)
	if err != nil {
		return nil, fmt.Errorf(
			"algorithm.createMazeFromCoord(%d, %d, %#v): %w",
			height,
			width,
			start,
			err,
		)
	}

	// minimum is needed so as not to remove all the cells if the sizes are small
	for range min(clearDeadEndNumber, min(height, width)-1) {
		cells = g.clearDeadEnd(cells, height, width, start, drawingChan)
	}

	return cells, nil
}

func mergeMazes(first, second domain.Maze, drawingChan chan<- domain.PaintingData,
	processID int,
) domain.Maze {
	height, width := first.Data.Height, first.Data.Width
	mergedCells := make([][]domain.CellType, height)

	for i := range height {
		mergedCells[i] = make([]domain.CellType, width)

		for j := range width {
			if first.Cells[i][j] == domain.Wall {
				mergedCells[i][j] = second.Cells[i][j]
			} else {
				mergedCells[i][j] = first.Cells[i][j]
			}

			drawingChan <- domain.NewPaintingData(i, j, mergedCells[i][j], processID, mergeDelay)
		}
	}

	return domain.NewMaze(first.Data, mergedCells)
}

func cellToPaintingData(c cell, id int) domain.PaintingData {
	return domain.NewPaintingData(c.Row, c.Col, c.Tpe, id, c.Delay)
}

func mergeChannels(out chan<- domain.PaintingData, in ...chan cell) {
	wg := &sync.WaitGroup{}

	for i, ch := range in {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			for data := range ch {
				out <- cellToPaintingData(data, id)
			}
		}(i + 1)
	}

	wg.Wait()
}

func (g *Generator) GenerateMaze(
	data domain.MazeData,
	paintingChan chan<- domain.PaintingData,
) (domain.Maze, error) {
	eg := &errgroup.Group{}

	var startMaze, endMaze domain.Maze

	ch1, ch2 := make(chan cell), make(chan cell)

	eg.Go(func() error {
		defer close(ch1)

		cells, err := g.generateMazeCellsFromCoord(data.Height, data.Width, data.Start, ch1)
		if err != nil {
			return fmt.Errorf("generating maze from coord: %w", err)
		}

		startMaze = domain.NewMaze(data, cells)

		return nil
	})

	eg.Go(func() error {
		defer close(ch2)

		cells, err := g.generateMazeCellsFromCoord(data.Height, data.Width, data.End, ch2)
		if err != nil {
			return fmt.Errorf("generating maze from coord: %w", err)
		}

		endMaze = domain.NewMaze(data, cells)

		return nil
	})

	eg.Go(func() error {
		mergeChannels(paintingChan, ch1, ch2)
		return nil
	})

	if err := eg.Wait(); err != nil {
		return domain.Maze{}, fmt.Errorf("errgroup: %w", err)
	}

	maze := mergeMazes(startMaze, endMaze, paintingChan, 0)

	return maze, nil
}
