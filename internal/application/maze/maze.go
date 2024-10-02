package maze

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/generator"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/painter"
)

func Start() error {
	ch := make(chan domain.CellRenderData)

	paint := painter.New(35, 70, os.Stdout)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		paint.Paint(context.Background(), ch)
	}()

	gen := generator.New()

	_, err := gen.GenerateMaze(35, 70, domain.NewCoord(0, 0), domain.NewCoord(0, 69), ch)
	if err != nil {
		return fmt.Errorf("generate maze: %w", err)
	}

	close(ch)

	wg.Wait()

	return nil
}
