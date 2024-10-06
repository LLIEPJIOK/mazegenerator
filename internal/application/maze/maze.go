package maze

import (
	"fmt"
	"os"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/generator"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/painter"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/pathfinder"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/ui"
)

func Start() error {
	out := os.Stdout
	gen := generator.New()
	paint := painter.New(out)
	pathFinder := pathfinder.New()

	userInterface := ui.NewUI(os.Stdin, out, gen, paint, pathFinder)
	if err := userInterface.Run(); err != nil {
		return fmt.Errorf("userInterface.Run(): %w", err)
	}

	return nil
}
