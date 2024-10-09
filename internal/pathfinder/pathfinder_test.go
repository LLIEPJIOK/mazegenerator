package pathfinder_test

import (
	"fmt"
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
	"github.com/es-debug/backend-academy-2024-go-template/internal/pathfinder"
	"github.com/stretchr/testify/require"
)

type PathFinder interface {
	ShortestPath(maze domain.Maze) ([]domain.Coord, bool)
}

func squareDist(first, second domain.Coord) int {
	dRow, dCol := first.Row-second.Row, first.Col-second.Col

	return dRow*dRow + dCol*dCol
}

func TestFindPathIfPathExists(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		data         domain.MazeData
		cells        [][]domain.CellType
		pathFinder   PathFinder
		shortestDist int
	}{
		{
			data: domain.NewMazeData(5, 5, domain.NewCoord(0, 0), domain.NewCoord(4, 4)),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
			},
			pathFinder:   pathfinder.NewAStar(),
			shortestDist: 8,
		},
		{
			data: domain.NewMazeData(5, 5, domain.NewCoord(0, 0), domain.NewCoord(0, 4)),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
			},
			pathFinder:   pathfinder.NewDijkstra(),
			shortestDist: 4,
		},
		{
			data: domain.NewMazeData(6, 6, domain.NewCoord(0, 0), domain.NewCoord(5, 0)),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1, 1},
				{0, 0, 0, 1, 0, 1},
				{0, 1, 0, 1, 0, 1},
				{0, 1, 0, 0, 0, 1},
				{0, 1, 1, 1, 0, 1},
				{1, 1, 1, 1, 1, 1},
			},
			pathFinder:   pathfinder.NewAStar(),
			shortestDist: 15,
		},
		{
			data: domain.NewMazeData(7, 7, domain.NewCoord(0, 3), domain.NewCoord(6, 3)),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1, 1, 1},
				{0, 0, 0, 1, 0, 0, 1},
				{1, 1, 0, 1, 0, 1, 1},
				{1, 1, 0, 0, 0, 1, 0},
				{1, 1, 1, 1, 0, 1, 1},
				{1, 0, 0, 0, 0, 0, 1},
				{1, 1, 1, 1, 1, 1, 1},
			},
			pathFinder:   pathfinder.NewDijkstra(),
			shortestDist: 14,
		},
		{
			data: domain.NewMazeData(4, 4, domain.NewCoord(0, 0), domain.NewCoord(3, 0)),
			cells: [][]domain.CellType{
				{1, 1, 0, 1},
				{0, 1, 0, 1},
				{0, 1, 0, 0},
				{1, 1, 1, 1},
			},
			pathFinder:   pathfinder.NewAStar(),
			shortestDist: 5,
		},
		{
			data: domain.NewMazeData(5, 5, domain.NewCoord(0, 0), domain.NewCoord(4, 0)),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1},
				{1, 0, 0, 0, 1},
				{1, 0, 1, 0, 1},
				{1, 0, 0, 0, 1},
				{1, 1, 1, 1, 1},
			},
			pathFinder:   pathfinder.NewDijkstra(),
			shortestDist: 4,
		},
		{
			data: domain.NewMazeData(6, 6, domain.NewCoord(0, 0), domain.NewCoord(5, 0)),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1, 1},
				{0, 0, 0, 1, 0, 1},
				{0, 1, 0, 1, 0, 1},
				{0, 1, 1, 1, 0, 1},
				{0, 1, 0, 0, 0, 1},
				{1, 1, 1, 1, 1, 1},
			},
			pathFinder:   pathfinder.NewAStar(),
			shortestDist: 11,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("#%d", i+1), func(t *testing.T) {
			t.Parallel()

			maze := domain.NewMaze(testCase.data, testCase.cells)
			path, ok := testCase.pathFinder.ShortestPath(maze)

			require.True(t, ok, "path must exist")
			require.NotEqual(t, 0, len(path), "path mustn't be empty")
			require.Equal(t, testCase.data.Start, path[0], "path must start from start point")
			require.Equal(t, testCase.data.End, path[len(path)-1], "path must end in end point")

			dist := 0

			for i := 1; i < len(path); i++ {
				coord := path[i]
				require.NotEqual(
					t,
					domain.Wall,
					maze.Cells[coord.Row][coord.Col],
					"path mustn't go throw",
				)

				require.Equal(t, 1, squareDist(path[i-1], path[i]))

				dist += int(maze.Cells[coord.Row][coord.Col])
			}

			require.Equalf(
				t,
				testCase.shortestDist,
				dist,
				"invalid shortest path",
				testCase.shortestDist,
				dist,
			)
		})
	}
}

func TestFindPathIfPathNotExists(t *testing.T) {
	testCases := []struct {
		data       domain.MazeData
		cells      [][]domain.CellType
		pathFinder PathFinder
	}{
		{
			data: domain.NewMazeData(5, 5, domain.NewCoord(0, 0), domain.NewCoord(4, 4)),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
			},
			pathFinder: pathfinder.NewDijkstra(),
		},
		{
			data: domain.NewMazeData(7, 7, domain.NewCoord(0, 3), domain.NewCoord(6, 3)),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1, 1, 1},
				{0, 1, 0, 1, 0, 0, 1},
				{1, 1, 0, 1, 0, 1, 1},
				{1, 1, 0, 0, 0, 0, 0},
				{1, 1, 1, 1, 0, 1, 1},
				{0, 0, 0, 0, 0, 0, 1},
				{1, 1, 1, 1, 1, 1, 1},
			},
			pathFinder: pathfinder.NewAStar(),
		},
		{
			data: domain.NewMazeData(4, 4, domain.NewCoord(0, 0), domain.NewCoord(3, 0)),
			cells: [][]domain.CellType{
				{1, 1, 0, 1},
				{0, 1, 1, 1},
				{0, 1, 0, 0},
				{1, 0, 1, 1},
			},
			pathFinder: pathfinder.NewDijkstra(),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("#%d", i+1), func(t *testing.T) {
			maze := domain.NewMaze(testCase.data, testCase.cells)
			_, ok := testCase.pathFinder.ShortestPath(maze)

			require.False(t, ok, "path mustn't exist")
		})
	}
}
