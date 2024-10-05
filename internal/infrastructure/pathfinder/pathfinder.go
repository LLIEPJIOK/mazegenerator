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

type PriorityQueue []queueItem

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(queueItem))
}

func (pq *PriorityQueue) Pop() interface{} {
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

func (p *PathFinder) FindPath(maze domain.Maze, start, end domain.Coord) ([]domain.Coord, bool) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, newQueueItem(start, domain.Coord{}, 0))

	prevCoords := make(map[domain.Coord]domain.Coord)

	for len(pq) != 0 {
		curItem := heap.Pop(&pq).(queueItem)

		if _, ok := prevCoords[curItem.curCoord]; ok {
			continue
		}

		prevCoords[curItem.curCoord] = curItem.prevCoord

		if curItem.curCoord == end {
			path := getPath(prevCoords, end, start)

			return path, true
		}

		for i := range p.dirCol {
			dRow, dCol := p.dirRow[i], p.dirCol[i]
			newRowID, newColID := curItem.curCoord.RowID+dRow, curItem.curCoord.ColID+dCol

			if min(newColID, newRowID) >= 0 && newRowID < maze.Height && newColID < maze.Width &&
				maze.Cells[newRowID][newColID] != domain.Wall {
				heap.Push(
					&pq,
					newQueueItem(
						domain.NewCoord(newRowID, newColID),
						curItem.curCoord,
						curItem.dist+int(maze.Cells[newRowID][newColID]),
					),
				)
			}
		}
	}

	return nil, false
}
