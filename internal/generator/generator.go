package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
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
	algo Algorithm
	dir  domain.Direction
}

func New(algo Algorithm) *Generator {
	return &Generator{
		algo: algo,
		dir:  domain.DefaultDirection(),
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

			for k := range g.dir.Rows {
				newRowID, newColID := i+g.dir.Rows[k], j+g.dir.Cols[k]
				if min(newRowID, newColID) < 0 || newRowID >= height || newColID >= width {
					continue
				}

				if cells[newRowID][newColID] != domain.Wall {
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

func mergeMazes(first, second domain.Maze, drawingChan chan<- domain.CellPaintingData,
	processID int,
) (domain.Maze, error) {
	height, width := first.Data.Height, first.Data.Width
	mergedCells := make([][]domain.CellType, height)

	for i := range height {
		mergedCells[i] = make([]domain.CellType, width)

		for j := range width {
			switch {
			case first.Cells[i][j] == domain.Wall:
				mergedCells[i][j] = second.Cells[i][j]
			case second.Cells[i][j] == domain.Wall:
				mergedCells[i][j] = first.Cells[i][j]
			default:
				numb, err := rand.Int(rand.Reader, big.NewInt(2))
				if err != nil {
					return domain.Maze{}, fmt.Errorf("generate random number: %w", err)
				}

				if numb.Int64() == 0 {
					mergedCells[i][j] = first.Cells[i][j]
				} else {
					mergedCells[i][j] = second.Cells[i][j]
				}
			}

			drawingChan <- domain.NewCellPaintingData(i, j, mergedCells[i][j], processID, mergeDelay)
		}
	}

	return domain.NewMaze(first.Data, mergedCells), nil
}

func cellToPaintingData(c cell, id int) domain.CellPaintingData {
	return domain.NewCellPaintingData(c.Row, c.Col, c.Tpe, id, c.Delay)
}

func mergeChannels(out chan<- domain.CellPaintingData, in ...chan cell) {
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
	paintingChan chan<- domain.CellPaintingData,
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

	maze, err := mergeMazes(startMaze, endMaze, paintingChan, 0)
	if err != nil {
		return domain.Maze{}, fmt.Errorf("merge mazes: %w", err)
	}

	return maze, nil
}

func randomCellType() (domain.CellType, error) {
	numb, err := rand.Int(rand.Reader, big.NewInt(15))
	if err != nil {
		return domain.Wall, fmt.Errorf("generate random number: %w", err)
	}

	switch numb.Int64() {
	case 0:
		return domain.Money, nil

	case 1:
		return domain.River, nil

	case 2:
		return domain.Sand, nil

	default:
		return domain.Passage, nil
	}
}
