package maze

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
	"github.com/es-debug/backend-academy-2024-go-template/internal/generator"
	"github.com/es-debug/backend-academy-2024-go-template/internal/painter"
	"github.com/es-debug/backend-academy-2024-go-template/internal/pathfinder"
	"github.com/es-debug/backend-academy-2024-go-template/internal/presentation"
)

const pathDrawingDelay = 50 * time.Millisecond

type pathFinder interface {
	ShortestPath(maze domain.Maze) ([]domain.Coord, bool)
}

func inputToMazeData(in *presentation.Input) domain.MazeData {
	return domain.NewMazeData(in.Height, in.Width, in.Start, in.End)
}

func generationAlgorithm(algo string) generator.Algorithm {
	switch algo {
	case "prim":
		return generator.NewPrim()

	case "backtrack":
		return generator.NewBacktrack()

	default:
		slog.Error("unknown generation algorithm")
	}

	return nil
}

func pathFinderAlgorithm(algo string) pathFinder {
	switch algo {
	case "dijkstra":
		return pathfinder.NewDijkstra()

	case "a-star":
		return pathfinder.NewAStar()

	default:
		slog.Error("unknown path finder algorithm")
	}

	return nil
}

func Start() error {
	input, output := os.Stdin, os.Stdout
	pres := presentation.New(input, output)

	inputData, err := pres.ProcessInput()
	if err != nil {
		return fmt.Errorf("process input: %w", err)
	}

	mazeData := inputToMazeData(inputData)

	gen := generator.New(generationAlgorithm(inputData.GenAlgo))
	paint := painter.New(output, mazeData.Height, mazeData.Width)
	pathFinder := pathFinderAlgorithm(inputData.PathFindAlgo)
	paintingChan := make(chan domain.PaintingData)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		paint.PaintGeneration(ctx, paintingChan)
	}()

	maze, err := gen.GenerateMaze(mazeData, paintingChan)
	if err != nil {
		cancel()
		return fmt.Errorf("generate maze: %w", err)
	}

	close(paintingChan)
	wg.Wait()

	if path, ok := pathFinder.ShortestPath(maze); ok {
		paint.PaintPath(path, pathDrawingDelay)
	} else {
		// ANSI code for red letters
		fmt.Println("\033[31mThere is no way between start and end points\033[0m")
	}

	return nil
}
