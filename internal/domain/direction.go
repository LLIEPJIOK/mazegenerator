package domain

type Direction struct {
	Rows []int
	Cols []int
}

func DefaultDirection() Direction {
	return Direction{
		Rows: []int{-1, 1, 0, 0},
		Cols: []int{0, 0, -1, 1},
	}
}
