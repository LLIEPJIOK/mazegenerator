package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
)

type Prim struct {
	dirRow []int
	dirCol []int
}

func NewPrim() *Prim {
	return &Prim{
		dirRow: []int{-1, 1, 0, 0},
		dirCol: []int{0, 0, -1, 1},
	}
}

func (p *Prim) createMazeFromCoord(
	height, width int,
	start domain.Coord,
	drawingChan chan<- domain.CellRenderData,
	processID int,
) (domain.Maze, error) {
	cells := make([][]domain.CellType, height)

	for i := range height {
		cells[i] = make([]domain.CellType, width)
	}

	waitList := make([]domain.Coord, 0)

	for i := range len(p.dirRow) {
		newRowID, newColID := start.Row+p.dirRow[i], start.Col+p.dirCol[i]
		if min(newRowID, newColID) < 0 || newRowID >= height || newColID >= width {
			continue
		}

		waitList = append(waitList, domain.NewCoord(newRowID, newColID))
	}

	cells[start.Row][start.Col] = domain.Passage
	drawingChan <- domain.NewCellRenderData(start.Row, start.Col, domain.Passage, processID, 3*time.Millisecond)

	for len(waitList) != 0 {
		randID, err := rand.Int(rand.Reader, big.NewInt(int64(len(waitList))))
		if err != nil {
			return domain.Maze{}, fmt.Errorf("generate random processID: %w", err)
		}

		randCoord := waitList[randID.Int64()]
		waitList[randID.Int64()], waitList[len(waitList)-1] = waitList[len(waitList)-1], waitList[randID.Int64()]
		waitList = waitList[:len(waitList)-1]

		cntWalls, cntBorders := 0, 0

		for i := range p.dirRow {
			newRowID, newColID := randCoord.Row+p.dirRow[i], randCoord.Col+p.dirCol[i]
			if min(newRowID, newColID) < 0 || newRowID >= height || newColID >= width {
				cntBorders++
				continue
			}

			if cells[newRowID][newColID] == domain.Wall {
				waitList = append(waitList, domain.NewCoord(newRowID, newColID))
				cntWalls++
			}
		}

		if cntWalls+cntBorders < 3 {
			waitList = waitList[:len(waitList)-cntWalls]
		} else {
			cells[randCoord.Row][randCoord.Col] = domain.Passage
			drawingChan <- domain.NewCellRenderData(randCoord.Row, randCoord.Col, domain.Passage, processID, 3*time.Millisecond)
		}
	}

	return domain.NewMaze(height, width, cells), nil
}
