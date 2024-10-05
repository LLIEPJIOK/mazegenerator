package domain

type Coord struct {
	RowID int
	ColID int
}

func NewCoord(rowID, colID int) Coord {
	return Coord{
		RowID: rowID,
		ColID: colID,
	}
}
