package generator

import (
	"fmt"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
	"golang.org/x/sync/errgroup"
)

type Algorithm interface {
	createMazeFromCoord(
		height, width int,
		start domain.Coord,
		drawingChan chan<- domain.CellRenderData,
		processID int,
	) (domain.Maze, error)
}

type Generator struct {
	dirRow []int
	dirCol []int
}

func New() *Generator {
	return &Generator{
		dirRow: []int{-1, 1, 0, 0},
		dirCol: []int{0, 0, -1, 1},
	}
}

func (g *Generator) clearDeadEnd(
	maze domain.Maze,
	start domain.Coord,
	drawingChan chan<- domain.CellRenderData,
	processID int,
) domain.Maze {
	newCells := make([][]domain.CellType, maze.Height)

	for i := range maze.Height {
		newCells[i] = make([]domain.CellType, maze.Width)

		for j := range maze.Width {
			cntPassages := 0

			for k := range g.dirRow {
				newRowID, newColID := i+g.dirRow[k], j+g.dirCol[k]
				if min(newRowID, newColID) < 0 || newRowID >= maze.Height || newColID >= maze.Width {
					continue
				}

				if maze.Cells[newRowID][newColID] == domain.Passage {
					cntPassages++
				}
			}

			if cntPassages == 1 && (i != start.Row || j != start.Col) {
				newCells[i][j] = domain.Wall
				drawingChan <- domain.NewCellRenderData(i, j, domain.Wall, processID, 100*time.Microsecond)
			} else {
				newCells[i][j] = maze.Cells[i][j]
			}
		}
	}

	return domain.NewMaze(maze.Height, maze.Width, newCells)
}

const clearDeadEndNumber = 10

func (g *Generator) generateMazeFromCoord(
	height, width int,
	start domain.Coord,
	algorithm Algorithm,
	drawingChan chan<- domain.CellRenderData,
	processID int,
) (domain.Maze, error) {
	maze, err := algorithm.createMazeFromCoord(height, width, start, drawingChan, processID)
	if err != nil {
		return domain.Maze{}, fmt.Errorf(
			"algorithm.createMazeFromCoord(%d, %d, %#v): %w",
			height,
			width,
			start,
			err,
		)
	}

	// minimum is needed so as not to remove all the cells if the sizes are small
	for range min(clearDeadEndNumber, min(height, width)-1) {
		maze = g.clearDeadEnd(maze, start, drawingChan, processID)
	}

	return maze, nil
}

func mergeMazes(first, second domain.Maze, drawingChan chan<- domain.CellRenderData,
	processID int,
) domain.Maze {
	mergedCells := make([][]domain.CellType, first.Height)

	for i := range first.Height {
		mergedCells[i] = make([]domain.CellType, first.Width)

		for j := range first.Width {
			if first.Cells[i][j] == domain.Wall {
				mergedCells[i][j] = second.Cells[i][j]
			} else {
				mergedCells[i][j] = first.Cells[i][j]
			}

			drawingChan <- domain.NewCellRenderData(i, j, mergedCells[i][j], processID, 50*time.Microsecond)
		}
	}

	return domain.NewMaze(first.Height, first.Width, mergedCells)
}

func (g *Generator) GenerateMaze(
	height, width int,
	start, end domain.Coord,
	algorithm Algorithm,
	drawingChan chan<- domain.CellRenderData,
) (domain.Maze, error) {
	eg := &errgroup.Group{}

	var startMaze, endMaze domain.Maze

	eg.Go(func() error {
		var err error

		startMaze, err = g.generateMazeFromCoord(height, width, start, algorithm, drawingChan, 1)
		if err != nil {
			return fmt.Errorf("g.generateMazeFromPoint(%d, %d, %#v): %w", height, width, start, err)
		}

		return nil
	})

	eg.Go(func() error {
		var err error

		endMaze, err = g.generateMazeFromCoord(height, width, end, algorithm, drawingChan, 2)
		if err != nil {
			return fmt.Errorf("g.generateMazeFromPoint(%d, %d, %#v): %w", height, width, end, err)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return domain.Maze{}, fmt.Errorf("errgroup: %w", err)
	}

	maze := mergeMazes(startMaze, endMaze, drawingChan, 0)

	return maze, nil
}
