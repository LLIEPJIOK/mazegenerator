package maze

import (
	"context"
	"fmt"
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

func Start() error {
	input, output := os.Stdin, os.Stdout
	pres := presentation.New(input, output)

	data, err := pres.ProcessInput()
	if err != nil {
		return fmt.Errorf("process input: %w", err)
	}

	gen := generator.New(generator.NewPrim())
	paint := painter.New(output, data.Height, data.Width)
	pathFinder := pathfinder.NewDijkstra()
	paintingChan := make(chan domain.PaintingData)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		paint.PaintGeneration(ctx, paintingChan)
	}()

	maze, err := gen.GenerateMaze(data, paintingChan)
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
