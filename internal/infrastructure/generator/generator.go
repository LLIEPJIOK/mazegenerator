package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
	"golang.org/x/sync/errgroup"
)

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

func (g *Generator) createMazeFromCoord(
	height, width int,
	start domain.Coord,
	drawingChan chan<- domain.CellRenderData,
	processID int,
) ([][]domain.CellType, error) {
	maze := make([][]domain.CellType, height)

	for i := range height {
		maze[i] = make([]domain.CellType, width)
	}

	waitList := make([]domain.Coord, 0)

	for i := range len(g.dirRow) {
		newRowID, newColID := start.RowID+g.dirRow[i], start.ColID+g.dirCol[i]
		if min(newRowID, newColID) < 0 || newRowID >= height || newColID >= width {
			continue
		}

		waitList = append(waitList, domain.NewCoord(newRowID, newColID))
	}

	maze[start.RowID][start.ColID] = domain.Passage
	drawingChan <- domain.NewCellRenderData(start.RowID, start.ColID, domain.Passage, processID, 3000)

	for len(waitList) != 0 {
		randID, err := rand.Int(rand.Reader, big.NewInt(int64(len(waitList))))
		if err != nil {
			return nil, fmt.Errorf("generate random processID: %w", err)
		}

		randCoord := waitList[randID.Int64()]
		waitList[randID.Int64()], waitList[len(waitList)-1] = waitList[len(waitList)-1], waitList[randID.Int64()]
		waitList = waitList[:len(waitList)-1]

		cntWalls, cntBorders := 0, 0

		for i := range g.dirRow {
			newRowID, newColID := randCoord.RowID+g.dirRow[i], randCoord.ColID+g.dirCol[i]
			if min(newRowID, newColID) < 0 || newRowID >= height || newColID >= width {
				cntBorders++
				continue
			}

			if maze[newRowID][newColID] == domain.Wall {
				waitList = append(waitList, domain.NewCoord(newRowID, newColID))
				cntWalls++
			}
		}

		if cntWalls+cntBorders < 3 {
			waitList = waitList[:len(waitList)-cntWalls]
		} else {
			maze[randCoord.RowID][randCoord.ColID] = domain.Passage
			drawingChan <- domain.NewCellRenderData(randCoord.RowID, randCoord.ColID, domain.Passage, processID, 3000)
		}
	}

	return maze, nil
}

func (g *Generator) clearDeadEnd(
	maze [][]domain.CellType,
	start domain.Coord,
	drawingChan chan<- domain.CellRenderData,
	processID int,
) [][]domain.CellType {
	height := len(maze)
	newMaze := make([][]domain.CellType, height)

	for i := range height {
		width := len(maze[i])
		newMaze[i] = make([]domain.CellType, width)

		for j := range width {
			cntPassages := 0

			for k := range g.dirRow {
				newRowID, newColID := i+g.dirRow[k], j+g.dirCol[k]
				if min(newRowID, newColID) < 0 || newRowID >= height || newColID >= width {
					continue
				}

				if maze[newRowID][newColID] == domain.Passage {
					cntPassages++
				}
			}

			if cntPassages == 1 && (i != start.RowID || j != start.ColID) {
				newMaze[i][j] = domain.Wall
				drawingChan <- domain.NewCellRenderData(i, j, domain.Wall, processID, 100)
			} else {
				newMaze[i][j] = maze[i][j]
			}
		}
	}

	return newMaze
}

const clearDeadEndNumber = 10

func (g *Generator) generateMazeFromPoint(
	height, width int,
	start domain.Coord,
	drawingChan chan<- domain.CellRenderData,
	processID int,
) ([][]domain.CellType, error) {
	maze, err := g.createMazeFromCoord(height, width, start, drawingChan, processID)
	if err != nil {
		return nil, fmt.Errorf(
			"g.createMazeFromCoord(%d, %d, %#v): %w",
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

func mergeMazes(first, second [][]domain.CellType, drawingChan chan<- domain.CellRenderData,
	processID int,
) [][]domain.CellType {
	height := len(first)
	mergedMaze := make([][]domain.CellType, height)

	for i := range height {
		width := len(first[i])
		mergedMaze[i] = make([]domain.CellType, width)

		for j := range width {
			if first[i][j] == domain.Wall {
				mergedMaze[i][j] = second[i][j]
			} else {
				mergedMaze[i][j] = first[i][j]
			}

			drawingChan <- domain.NewCellRenderData(i, j, mergedMaze[i][j], processID, 50)
		}
	}

	return mergedMaze
}

func (g *Generator) GenerateMaze(
	height, width int,
	start, end domain.Coord,
	drawingChan chan<- domain.CellRenderData,
) (domain.Maze, error) {
	eg := &errgroup.Group{}

	var startMaze, endMaze [][]domain.CellType

	eg.Go(func() error {
		var err error

		startMaze, err = g.generateMazeFromPoint(height, width, start, drawingChan, 1)
		if err != nil {
			return fmt.Errorf("g.generateMazeFromPoint(%d, %d, %#v): %w", height, width, start, err)
		}

		return nil
	})

	eg.Go(func() error {
		var err error

		endMaze, err = g.generateMazeFromPoint(height, width, end, drawingChan, 2)
		if err != nil {
			return fmt.Errorf("g.generateMazeFromPoint(%d, %d, %#v): %w", height, width, end, err)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return domain.Maze{}, fmt.Errorf("errgroup: %w", err)
	}

	maze := mergeMazes(startMaze, endMaze, drawingChan, 0)

	return domain.NewMaze(height, width, maze), nil
}
