package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
)

type Backtrack struct {
	dirRow []int
	dirCol []int
}

func NewBacktrack() *Backtrack {
	return &Backtrack{
		dirRow: []int{-1, 1, 0, 0},
		dirCol: []int{0, 0, -1, 1},
	}
}

const (
	forkCoeff = 3
)

func (b *Backtrack) createMazeFromCoord(
	height, width int,
	start domain.Coord,
	drawingChan chan<- domain.CellRenderData,
	processID int,
) (domain.Maze, error) {
	cells := make([][]domain.CellType, height)

	for i := range height {
		cells[i] = make([]domain.CellType, width)
	}

	stack := []domain.Coord{start}

	for len(stack) > 0 {
		curCoord := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		cntWalls := 0

		for i := range b.dirRow {
			newRow, newCol := curCoord.Row+b.dirRow[i], curCoord.Col+b.dirCol[i]
			if min(newRow, newCol) < 0 || newRow >= height || newCol >= width ||
				cells[newRow][newCol] == domain.Wall {
				cntWalls++
				continue
			}
		}

		if cntWalls < 3 {
			continue
		}

		cells[curCoord.Row][curCoord.Col] = domain.Passage
		drawingChan <- domain.NewCellRenderData(curCoord.Row, curCoord.Col, domain.Passage, processID, 3*time.Millisecond)

		prevRands := make(map[int64]struct{})

		for len(prevRands) != forkCoeff {
			randID, err := rand.Int(rand.Reader, big.NewInt(int64(len(b.dirRow))))
			if err != nil {
				return domain.Maze{}, fmt.Errorf("generate random processID: %w", err)
			}

			if _, ok := prevRands[randID.Int64()]; ok {
				continue
			}

			newRow, newCol := curCoord.Row+b.dirRow[randID.Int64()], curCoord.Col+b.dirCol[randID.Int64()]

			if newRow >= 0 && newRow < height && newCol >= 0 && newCol < width && cells[newRow][newCol] == domain.Wall {
				stack = append(stack, domain.NewCoord(newRow, newCol))
			}

			prevRands[randID.Int64()] = struct{}{}
		}
	}

	return domain.NewMaze(height, width, cells), nil
}
