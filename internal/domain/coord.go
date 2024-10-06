package domain

type Coord struct {
	Row int
	Col int
}

func NewCoord(rowID, colID int) Coord {
	return Coord{
		Row: rowID,
		Col: colID,
	}
}
