## About Maze Generator

It's a program for creating mazes and finding the shortest path between two points.

## Getting Started

1. Open a terminal that supports ANSI and run the following command:

   ```shell
   git clone git@github.com:central-university-dev/backend-academy-2024-fall.git
   ```

2. Navigate to the project folder:

   ```shell
   cd backend_academy_2024_project_2-go-LLIEPJIOK
   ```

3. Run the program:

   ```bash
   go run cmd/maze/main.go
   ```

4. Follow the instructions provided by the program. It's quite simple.

## Algorithms

### Maze Generation

Two algorithms are implemented for maze generation:

- Prim's Algorithm
- Backtracking Algorithm

### Pathfinding

Two algorithms are implemented for pathfinding:

- Dijkstra's Algorithm
- A\* (A-star) with Manhattan distance as the heuristic

## Cell Types

The following cell types are used in the maze:

- Wall
- Passage (cost = 3)
- Money (cost = 1)
- Sand (cost = 5)
- River (cost = 10)
- Ambiguous
- Path
