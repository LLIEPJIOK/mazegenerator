package domain

type Coord struct {
	X int
	Y int
}

func NewCoord(x, y int) Coord {
	return Coord{
		X: x,
		Y: y,
	}
}
