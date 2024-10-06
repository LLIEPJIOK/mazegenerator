package domain

type Maze struct {
	Height int
	Width  int
	Cells  [][]CellType
}

func NewMaze(height, width int, cells [][]CellType) Maze {
	return Maze{
		Height: height,
		Width:  width,
		Cells:  cells,
	}
}

type DrawingMaze struct {
	cells [][]map[int]CellType
}

func NewDrawingMaze(height, width int) DrawingMaze {
	cells := make([][]map[int]CellType, height)

	for i := range height {
		cells[i] = make([]map[int]CellType, width)
		for j := range width {
			cells[i][j] = make(map[int]CellType)
		}
	}

	return DrawingMaze{
		cells: cells,
	}
}

func (dm *DrawingMaze) GetCellType(x, y int) CellType {
	switch len(dm.cells[x][y]) {
	case 0:
		return Wall
	case 1:
		for _, v := range dm.cells[x][y] {
			return v
		}
	}

	return Guessing
}

func (dm *DrawingMaze) AddCellType(cellData CellRenderData) {
	if cellData.Tpe == Wall {
		delete(dm.cells[cellData.Row][cellData.Col], cellData.SenderID)
	} else {
		if cellData.SenderID == 0 {
			clear(dm.cells[cellData.Row][cellData.Col])
		}

		dm.cells[cellData.Row][cellData.Col][cellData.SenderID] = cellData.Tpe
	}
}
