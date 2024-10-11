package painter

import "github.com/es-debug/backend-academy-2024-go-template/internal/domain"

type PaintingMaze struct {
	cells [][]map[int]domain.CellType
}

func newPaintingMaze(height, width int) PaintingMaze {
	cells := make([][]map[int]domain.CellType, height)

	for i := range height {
		cells[i] = make([]map[int]domain.CellType, width)
		for j := range width {
			cells[i][j] = make(map[int]domain.CellType)
		}
	}

	return PaintingMaze{
		cells: cells,
	}
}

func (dm *PaintingMaze) GetCellType(x, y int) domain.CellType {
	switch len(dm.cells[x][y]) {
	case 0:
		return domain.Wall
	case 1:
		for _, v := range dm.cells[x][y] {
			return v
		}
	}

	return domain.Guessing
}

func (dm *PaintingMaze) AddCellType(cellData domain.CellPaintingData) {
	if cellData.Tpe == domain.Wall {
		delete(dm.cells[cellData.Row][cellData.Col], cellData.SenderID)
	} else {
		if cellData.SenderID == 0 {
			clear(dm.cells[cellData.Row][cellData.Col])
		}

		dm.cells[cellData.Row][cellData.Col][cellData.SenderID] = cellData.Tpe
	}
}
