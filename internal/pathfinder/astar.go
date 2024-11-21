package pathfinder

import (
	"container/heap"
	"math"

	"github.com/LLIEPJIOK/mazegenerator/internal/domain"
)

type aStarItem struct {
	curCoord  domain.Coord
	prevCoord domain.Coord
	G         int
	H         int
	f         int
}

func newAStarItem(curCoord, prevCoord domain.Coord, g, h, f int) aStarItem {
	return aStarItem{
		curCoord:  curCoord,
		prevCoord: prevCoord,
		G:         g,
		H:         h,
		f:         f,
	}
}

func (a aStarItem) dist() int {
	return a.f
}

type AStar struct {
	dir domain.Direction
}

func NewAStar() *AStar {
	return &AStar{
		dir: domain.DefaultDirection(),
	}
}

// Function for calculating Manhattan distance.
func heuristic(a, b domain.Coord) int {
	return int(math.Abs(float64(a.Row-b.Row)) + math.Abs(float64(a.Col-b.Col)))
}

func (a *AStar) ShortestPath(maze domain.Maze, pathChan chan<- []domain.Coord) ([]domain.Coord, bool) {
	pq := make(priorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, newAStarItem(maze.Data.Start, domain.Coord{}, 0, heuristic(maze.Data.Start, maze.Data.End), 0))

	prevCoords := make(map[domain.Coord]domain.Coord)

	for len(pq) != 0 {
		curItem := heap.Pop(&pq).(aStarItem)

		if _, ok := prevCoords[curItem.curCoord]; ok {
			continue
		}

		prevCoords[curItem.curCoord] = curItem.prevCoord

		curPath := getPath(prevCoords, curItem.curCoord, maze.Data.Start)
		pathChan <- curPath

		if curItem.curCoord == maze.Data.End {
			path := getPath(prevCoords, maze.Data.End, maze.Data.Start)

			return path, true
		}

		for i := range a.dir.Cols {
			dRow, dCol := a.dir.Rows[i], a.dir.Cols[i]
			newRow, newCol := curItem.curCoord.Row+dRow, curItem.curCoord.Col+dCol

			if min(newCol, newRow) >= 0 && newRow < maze.Data.Height && newCol < maze.Data.Width &&
				maze.Cells[newRow][newCol].IsTraversable() {
				newG := curItem.G + maze.Cells[newRow][newCol].Cost()
				newH := heuristic(curItem.curCoord, domain.NewCoord(newRow, newCol))
				newDist := newG + newH
				heap.Push(
					&pq,
					newAStarItem(
						domain.NewCoord(newRow, newCol),
						curItem.curCoord,
						newG,
						newH,
						newDist,
					),
				)
			}
		}
	}

	pathChan <- nil

	return nil, false
}
