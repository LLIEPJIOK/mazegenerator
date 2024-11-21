package pathfinder

import (
	"slices"

	"github.com/LLIEPJIOK/mazegenerator/internal/domain"
)

func getPath(
	prevCoords map[domain.Coord]domain.Coord,
	curPoint, startPoint domain.Coord,
) []domain.Coord {
	path := make([]domain.Coord, 0)

	for curPoint != startPoint {
		path = append(path, curPoint)
		curPoint = prevCoords[curPoint]
	}

	path = append(path, curPoint)

	slices.Reverse(path)

	return path
}
