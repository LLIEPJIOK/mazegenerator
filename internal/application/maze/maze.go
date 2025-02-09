package maze

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/LLIEPJIOK/mazegenerator/internal/domain"
	"github.com/LLIEPJIOK/mazegenerator/internal/generator"
	"github.com/LLIEPJIOK/mazegenerator/internal/painter"
	"github.com/LLIEPJIOK/mazegenerator/internal/pathfinder"
	"github.com/LLIEPJIOK/mazegenerator/internal/presentation"
)

const pathDrawingDelay = 50 * time.Millisecond

type pathFinder interface {
	ShortestPath(maze domain.Maze, pathChan chan<- []domain.Coord) ([]domain.Coord, bool)
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
	paint := painter.New(output, mazeData)
	pathFinder := pathFinderAlgorithm(inputData.PathFindAlgo)
	paintingChan := make(chan domain.CellPaintingData)

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

	pathChan := make(chan []domain.Coord)

	wg.Add(1)

	go func() {
		defer wg.Done()
		paint.PaintPath(pathChan, pathDrawingDelay)
	}()

	_, ok := pathFinder.ShortestPath(maze, pathChan)

	close(pathChan)
	wg.Wait()

	if !ok {
		fmt.Println("There is no way between start and end points")
	}

	return nil
}
