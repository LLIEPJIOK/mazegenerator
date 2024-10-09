package generator

type direction struct {
	dirRow []int
	dirCol []int
}

func defaultDirection() direction {
	return direction{
		dirRow: []int{-1, 1, 0, 0},
		dirCol: []int{0, 0, -1, 1},
	}
}
