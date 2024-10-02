package generator_test

import (
	"fmt"
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/generator"
	"github.com/stretchr/testify/require"
)

func compareMazes(t *testing.T, maze domain.Maze, drawMaze domain.DrawingMaze) {
	for i, row := range maze.Cells {
		for j, v := range row {
			if drawMaze.GetCellType(i, j) != v {
				t.Fatalf("maze and drawing maze cells differs in (%d, %d)", i, j)
			}
		}
	}
}

func TestGenerateMaze(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		height int
		width  int
		start  domain.Coord
		end    domain.Coord
		ch     chan domain.CellRenderData
	}{
		{
			height: 5,
			width:  5,
			start:  domain.NewCoord(0, 0),
			end:    domain.NewCoord(0, 4),
			ch:     make(chan domain.CellRenderData),
		},
		{
			height: 10,
			width:  10,
			start:  domain.NewCoord(0, 0),
			end:    domain.NewCoord(9, 9),
			ch:     make(chan domain.CellRenderData),
		},
		{
			height: 7,
			width:  7,
			start:  domain.NewCoord(0, 3),
			end:    domain.NewCoord(6, 3),
			ch:     make(chan domain.CellRenderData),
		},
		{
			height: 6,
			width:  8,
			start:  domain.NewCoord(0, 7),
			end:    domain.NewCoord(5, 0),
			ch:     make(chan domain.CellRenderData),
		},
		{
			height: 8,
			width:  8,
			start:  domain.NewCoord(7, 0),
			end:    domain.NewCoord(7, 7),
			ch:     make(chan domain.CellRenderData),
		},
		{
			height: 12,
			width:  12,
			start:  domain.NewCoord(0, 6),
			end:    domain.NewCoord(11, 11),
			ch:     make(chan domain.CellRenderData),
		},
		{
			height: 15,
			width:  10,
			start:  domain.NewCoord(14, 0),
			end:    domain.NewCoord(14, 9),
			ch:     make(chan domain.CellRenderData),
		},
	}

	gen := generator.New()

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			t.Parallel()

			var maze domain.Maze

			var err error

			go func() {
				defer close(testCase.ch)
				maze, err = gen.GenerateMaze(testCase.height, testCase.width, testCase.start, testCase.end, testCase.ch)
			}()

			drawMaze := domain.NewDrawingMaze(testCase.height, testCase.width)

			for cellData := range testCase.ch {
				drawMaze.AddCellType(cellData)
			}

			require.NoError(t, err, "generate maze should return nil error")
			compareMazes(t, maze, drawMaze)
		})
	}
}
