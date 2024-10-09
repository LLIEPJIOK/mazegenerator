package pathfinder

import (
	"container/heap"
	"slices"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
)

type queueItem struct {
	curCoord  domain.Coord
	prevCoord domain.Coord
	dist      int
}

func newQueueItem(curCoord, prevCoord domain.Coord, dist int) queueItem {
	return queueItem{
		curCoord:  curCoord,
		prevCoord: prevCoord,
		dist:      dist,
	}
}

type priorityQueue []queueItem

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(queueItem))
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]

	return item
}

type PathFinder struct {
	dirRow []int
	dirCol []int
}

func New() *PathFinder {
	return &PathFinder{
		dirRow: []int{-1, 1, 0, 0},
		dirCol: []int{0, 0, -1, 1},
	}
}

func getPath(
	prevCoords map[domain.Coord]domain.Coord,
	curPoint, startPoint domain.Coord,
) []domain.Coord {
	path := make([]domain.Coord, 0)

	for curPoint != startPoint {
		path = append(path, curPoint)
		curPoint = prevCoords[curPoint]
	}

	path = append(path, curPoint)

	slices.Reverse(path)

	return path
}

func (p *PathFinder) FindPath(maze domain.Maze) ([]domain.Coord, bool) {
	pq := make(priorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, newQueueItem(maze.Data.Start, domain.Coord{}, 0))

	prevCoords := make(map[domain.Coord]domain.Coord)

	for len(pq) != 0 {
		curItem := heap.Pop(&pq).(queueItem)

		if _, ok := prevCoords[curItem.curCoord]; ok {
			continue
		}

		prevCoords[curItem.curCoord] = curItem.prevCoord

		if curItem.curCoord == maze.Data.End {
			path := getPath(prevCoords, maze.Data.End, maze.Data.Start)

			return path, true
		}

		for i := range p.dirCol {
			dRow, dCol := p.dirRow[i], p.dirCol[i]
			newRow, newCol := curItem.curCoord.Row+dRow, curItem.curCoord.Col+dCol

			if min(newCol, newRow) >= 0 && newRow < maze.Data.Height && newCol < maze.Data.Width &&
				maze.Cells[newRow][newCol] != domain.Wall {
				heap.Push(
					&pq,
					newQueueItem(
						domain.NewCoord(newRow, newCol),
						curItem.curCoord,
						curItem.dist+int(maze.Cells[newRow][newCol]),
					),
				)
			}
		}
	}

	return nil, false
}
