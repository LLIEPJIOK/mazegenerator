package domain

type MazeData struct {
	Height int
	Width  int
	Start  Coord
	End    Coord
}

func NewMazeData(height, width int, start, end Coord) MazeData {
	return MazeData{
		Height: height,
		Width:  width,
		Start:  start,
		End:    end,
	}
}

type Maze struct {
	Data  MazeData
	Cells [][]CellType
}

func NewMaze(data MazeData, cells [][]CellType) Maze {
	return Maze{
		Data:  data,
		Cells: cells,
	}
}
