package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/LLIEPJIOK/mazegenerator/internal/application/maze"
)

func main() {
	if err := maze.Start(); err != nil {
		slog.Error(fmt.Sprintf("maze.Start(): %s", err))
		os.Exit(1)
	}
}
