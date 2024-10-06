package maze

import (
	"fmt"
	"os"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/ui"
)

func Start() error {
	userInterface := ui.New(os.Stdin, os.Stdout)
	if err := userInterface.Run(); err != nil {
		return fmt.Errorf("userInterface.Run(): %w", err)
	}

	return nil
}
