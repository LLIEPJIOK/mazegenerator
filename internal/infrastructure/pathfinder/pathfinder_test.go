package pathfinder_test

import (
	"fmt"
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/pathfinder"
	"github.com/stretchr/testify/require"
)

func TestFindPathIfPathExists(t *testing.T) {
	testCases := []struct {
		height       int
		width        int
		start        domain.Coord
		end          domain.Coord
		cells        [][]domain.CellType
		shortestDist int
	}{
		{
			height: 5,
			width:  5,
			start:  domain.NewCoord(0, 0),
			end:    domain.NewCoord(4, 4),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
			},
			shortestDist: 8,
		},
		{
			height: 5,
			width:  5,
			start:  domain.NewCoord(0, 0),
			end:    domain.NewCoord(0, 4),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
			},
			shortestDist: 4,
		},
		{
			height: 6,
			width:  6,
			start:  domain.NewCoord(0, 0),
			end:    domain.NewCoord(5, 0),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1, 1},
				{0, 0, 0, 1, 0, 1},
				{0, 1, 0, 1, 0, 1},
				{0, 1, 0, 0, 0, 1},
				{0, 1, 1, 1, 0, 1},
				{1, 1, 1, 1, 1, 1},
			},
			shortestDist: 15,
		},
		{
			height: 7,
			width:  7,
			start:  domain.NewCoord(0, 3),
			end:    domain.NewCoord(6, 3),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1, 1, 1},
				{0, 0, 0, 1, 0, 0, 1},
				{1, 1, 0, 1, 0, 1, 1},
				{1, 1, 0, 0, 0, 1, 0},
				{1, 1, 1, 1, 0, 1, 1},
				{1, 0, 0, 0, 0, 0, 1},
				{1, 1, 1, 1, 1, 1, 1},
			},
			shortestDist: 14,
		},
		{
			height: 4,
			width:  4,
			start:  domain.NewCoord(0, 0),
			end:    domain.NewCoord(3, 0),
			cells: [][]domain.CellType{
				{1, 1, 0, 1},
				{0, 1, 0, 1},
				{0, 1, 0, 0},
				{1, 1, 1, 1},
			},
			shortestDist: 5,
		},
		{
			height: 5,
			width:  5,
			start:  domain.NewCoord(0, 0),
			end:    domain.NewCoord(4, 0),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1},
				{1, 0, 0, 0, 1},
				{1, 0, 1, 0, 1},
				{1, 0, 0, 0, 1},
				{1, 1, 1, 1, 1},
			},
			shortestDist: 4,
		},
		{
			height: 6,
			width:  6,
			start:  domain.NewCoord(0, 0),
			end:    domain.NewCoord(5, 0),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1, 1},
				{0, 0, 0, 1, 0, 1},
				{0, 1, 0, 1, 0, 1},
				{0, 1, 1, 1, 0, 1},
				{0, 1, 0, 0, 0, 1},
				{1, 1, 1, 1, 1, 1},
			},
			shortestDist: 11,
		},
	}

	pathFinder := pathfinder.New()

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("#%d", i+1), func(t *testing.T) {
			maze := domain.NewMaze(testCase.height, testCase.width, testCase.cells)
			path, ok := pathFinder.FindPath(maze, testCase.start, testCase.end)

			require.True(t, ok, "path must exist")
			require.NotEqual(t, 0, len(path), "path mustn't be empty")
			require.Equal(t, testCase.start, path[0], "path must start from start point")
			require.Equal(t, testCase.end, path[len(path)-1], "path must end in end point")

			dist := 0

			for i := 1; i < len(path); i++ {
				coord := path[i]
				require.NotEqual(
					t,
					domain.Wall,
					maze.Cells[coord.RowID][coord.ColID],
					"path mustn't go throw",
				)

				dist += int(maze.Cells[coord.RowID][coord.ColID])
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
		height int
		width  int
		start  domain.Coord
		end    domain.Coord
		cells  [][]domain.CellType
	}{
		{
			height: 5,
			width:  5,
			start:  domain.NewCoord(0, 0),
			end:    domain.NewCoord(4, 4),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 1},
				{0, 0, 0, 0, 1},
			},
		},
		{
			height: 7,
			width:  7,
			start:  domain.NewCoord(0, 3),
			end:    domain.NewCoord(6, 3),
			cells: [][]domain.CellType{
				{1, 1, 1, 1, 1, 1, 1},
				{0, 1, 0, 1, 0, 0, 1},
				{1, 1, 0, 1, 0, 1, 1},
				{1, 1, 0, 0, 0, 0, 0},
				{1, 1, 1, 1, 0, 1, 1},
				{0, 0, 0, 0, 0, 0, 1},
				{1, 1, 1, 1, 1, 1, 1},
			},
		},
		{
			height: 4,
			width:  4,
			start:  domain.NewCoord(0, 0),
			end:    domain.NewCoord(3, 0),
			cells: [][]domain.CellType{
				{1, 1, 0, 1},
				{0, 1, 1, 1},
				{0, 1, 0, 0},
				{1, 0, 1, 1},
			},
		},
	}

	pathFinder := pathfinder.New()

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("#%d", i+1), func(t *testing.T) {
			maze := domain.NewMaze(testCase.height, testCase.width, testCase.cells)
			_, ok := pathFinder.FindPath(maze, testCase.start, testCase.end)

			require.False(t, ok, "path mustn't exist")
		})
	}
}
