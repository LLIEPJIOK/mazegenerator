package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
	"golang.org/x/sync/errgroup"
)

type Generator struct {
	dirX []int
	dirY []int
}

func NewGenerator() *Generator {
	return &Generator{
		dirX: []int{-1, 1, 0, 0},
		dirY: []int{0, 0, -1, 1},
	}
}

func (g *Generator) createMazeFromCoord(
	height, width int,
	startCoord domain.Coord,
	drawingChan chan<- domain.CellRenderData,
	processID int,
) ([][]domain.CellType, error) {
	maze := make([][]domain.CellType, height)

	for i := range height {
		maze[i] = make([]domain.CellType, width)
	}

	waitList := make([]domain.Coord, 0)

	for i := range len(g.dirX) {
		nx, ny := startCoord.X+g.dirX[i], startCoord.Y+g.dirY[i]
		if min(nx, ny) < 0 || nx >= height || ny >= width {
			continue
		}

		waitList = append(waitList, domain.NewCoord(nx, ny))
	}

	maze[startCoord.X][startCoord.Y] = domain.Passage
	drawingChan <- domain.NewCellRenderData(startCoord.X, startCoord.Y, domain.Passage, processID, 3000)

	for len(waitList) != 0 {
		randID, err := rand.Int(rand.Reader, big.NewInt(int64(len(waitList))))
		if err != nil {
			return nil, fmt.Errorf("generate random processID: %w", err)
		}

		randCoord := waitList[randID.Int64()]
		waitList[randID.Int64()], waitList[len(waitList)-1] = waitList[len(waitList)-1], waitList[randID.Int64()]
		waitList = waitList[:len(waitList)-1]

		cntWalls, cntBorders := 0, 0

		for i := range g.dirX {
			nx, ny := randCoord.X+g.dirX[i], randCoord.Y+g.dirY[i]
			if min(nx, ny) < 0 || nx >= height || ny >= width {
				cntBorders++
				continue
			}

			if maze[nx][ny] == domain.Wall {
				waitList = append(waitList, domain.NewCoord(nx, ny))
				cntWalls++
			}
		}

		if cntWalls+cntBorders < 3 {
			waitList = waitList[:len(waitList)-cntWalls]
		} else {
			maze[randCoord.X][randCoord.Y] = domain.Passage
			drawingChan <- domain.NewCellRenderData(randCoord.X, randCoord.Y, domain.Passage, processID, 3000)
		}
	}

	return maze, nil
}

func (g *Generator) clearDeadEnd(
	maze [][]domain.CellType,
	startCoord domain.Coord,
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

			for k := range g.dirX {
				nx, ny := i+g.dirX[k], j+g.dirY[k]
				if min(nx, ny) < 0 || nx >= height || ny >= width {
					continue
				}

				if maze[nx][ny] == domain.Passage {
					cntPassages++
				}
			}

			if cntPassages == 1 && (i != startCoord.X || j != startCoord.Y) {
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
	startCoord domain.Coord,
	drawingChan chan<- domain.CellRenderData,
	processID int,
) ([][]domain.CellType, error) {
	maze, err := g.createMazeFromCoord(height, width, startCoord, drawingChan, processID)
	if err != nil {
		return nil, fmt.Errorf(
			"g.createMazeFromCoord(%d, %d, %#v): %w",
			height,
			width,
			startCoord,
			err,
		)
	}

	for range min(clearDeadEndNumber, min(height, width)-1) {
		maze = g.clearDeadEnd(maze, startCoord, drawingChan, processID)
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
) (*domain.Maze, error) {
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
		return nil, fmt.Errorf("errgroup: %w", err)
	}

	maze := mergeMazes(startMaze, endMaze, drawingChan, 0)

	return domain.NewMaze(height, width, maze), nil
}
