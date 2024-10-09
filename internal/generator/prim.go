package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
)

type Prim struct {
	direction
}

func NewPrim() *Prim {
	return &Prim{
		direction: defaultDirection(),
	}
}

func (p *Prim) createMazeCellsFromCoord(
	height, width int,
	start domain.Coord,
	drawingChan chan<- cell,
) ([][]domain.CellType, error) {
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
	drawingChan <- newCell(start.Row, start.Col, domain.Passage, drawingDelay)

	for len(waitList) != 0 {
		randID, err := rand.Int(rand.Reader, big.NewInt(int64(len(waitList))))
		if err != nil {
			return nil, fmt.Errorf("generate random processID: %w", err)
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
			drawingChan <- newCell(randCoord.Row, randCoord.Col, domain.Passage, drawingDelay)
		}
	}

	return cells, nil
}
