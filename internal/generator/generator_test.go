package generator_test

import (
	"fmt"
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
	"github.com/es-debug/backend-academy-2024-go-template/internal/generator"
	"github.com/stretchr/testify/require"
)

func compareMazes(t *testing.T, maze domain.Maze, drawMaze [][]domain.CellType) {
	for i, row := range maze.Cells {
		for j, v := range row {
			if drawMaze[i][j] != v {
				t.Fatalf("maze and drawing maze cells differs in (%d, %d)", i, j)
			}
		}
	}
}

func TestGenerateMaze(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		data      domain.MazeData
		algorithm generator.Algorithm
		ch        chan domain.CellPaintingData
	}{
		{
			data:      domain.NewMazeData(5, 5, domain.NewCoord(0, 0), domain.NewCoord(0, 4)),
			algorithm: generator.NewBacktrack(),
			ch:        make(chan domain.CellPaintingData),
		},
		{
			data:      domain.NewMazeData(10, 10, domain.NewCoord(0, 0), domain.NewCoord(9, 9)),
			algorithm: generator.NewPrim(),
			ch:        make(chan domain.CellPaintingData),
		},
		{
			data:      domain.NewMazeData(7, 7, domain.NewCoord(0, 3), domain.NewCoord(6, 3)),
			algorithm: generator.NewBacktrack(),
			ch:        make(chan domain.CellPaintingData),
		},
		{
			data:      domain.NewMazeData(6, 8, domain.NewCoord(0, 7), domain.NewCoord(5, 0)),
			algorithm: generator.NewPrim(),
			ch:        make(chan domain.CellPaintingData),
		},
		{
			data:      domain.NewMazeData(8, 8, domain.NewCoord(7, 0), domain.NewCoord(7, 7)),
			algorithm: generator.NewBacktrack(),
			ch:        make(chan domain.CellPaintingData),
		},
		{
			data:      domain.NewMazeData(12, 12, domain.NewCoord(0, 6), domain.NewCoord(11, 11)),
			algorithm: generator.NewPrim(),
			ch:        make(chan domain.CellPaintingData),
		},
		{
			data:      domain.NewMazeData(15, 10, domain.NewCoord(14, 0), domain.NewCoord(14, 9)),
			algorithm: generator.NewBacktrack(),
			ch:        make(chan domain.CellPaintingData),
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("#%d", i+1), func(t *testing.T) {
			t.Parallel()

			gen := generator.New(testCase.algorithm)

			var maze domain.Maze

			var err error

			go func() {
				defer close(testCase.ch)
				maze, err = gen.GenerateMaze(
					testCase.data,
					testCase.ch,
				)
			}()

			drawMaze := make([][]domain.CellType, testCase.data.Height)
			for i := range testCase.data.Height {
				drawMaze[i] = make([]domain.CellType, testCase.data.Width)
			}

			for cellData := range testCase.ch {
				drawMaze[cellData.Row][cellData.Col] = cellData.Tpe
			}

			require.NoError(t, err, "generate maze should return nil error")
			compareMazes(t, maze, drawMaze)
		})
	}
}
