package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
)

type Backtrack struct {
	dir domain.Direction
}

func NewBacktrack() *Backtrack {
	return &Backtrack{
		dir: domain.DefaultDirection(),
	}
}

const (
	forkCoeff = 3
)

func (b *Backtrack) createMazeCellsFromCoord(
	height, width int,
	start domain.Coord,
	drawingChan chan<- cell,
) ([][]domain.CellType, error) {
	cells := make([][]domain.CellType, height)

	for i := range height {
		cells[i] = make([]domain.CellType, width)
	}

	stack := []domain.Coord{start}

	for len(stack) > 0 {
		curCoord := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		cntWalls := 0

		for i := range b.dir.Rows {
			newRow, newCol := curCoord.Row+b.dir.Rows[i], curCoord.Col+b.dir.Cols[i]
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
		drawingChan <- newCell(curCoord.Row, curCoord.Col, domain.Passage, drawingDelay)

		prevRands := make(map[int64]struct{})

		for len(prevRands) != forkCoeff {
			randID, err := rand.Int(rand.Reader, big.NewInt(int64(len(b.dir.Rows))))
			if err != nil {
				return nil, fmt.Errorf("generate random processID: %w", err)
			}

			if _, ok := prevRands[randID.Int64()]; ok {
				continue
			}

			newRow, newCol := curCoord.Row+b.dir.Rows[randID.Int64()], curCoord.Col+b.dir.Cols[randID.Int64()]

			if newRow >= 0 && newRow < height && newCol >= 0 && newCol < width && cells[newRow][newCol] == domain.Wall {
				stack = append(stack, domain.NewCoord(newRow, newCol))
			}

			prevRands[randID.Int64()] = struct{}{}
		}
	}

	return cells, nil
}
