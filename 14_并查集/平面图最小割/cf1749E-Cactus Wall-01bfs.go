// 2 种遇到有后效性动态规划的处理方法：高斯消元和最短路
// https://www.luogu.com.cn/article/urd6r0r7

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF32 int32 = 1e9 + 10

// cf1749E-Cactus Wall
// https://www.luogu.com.cn/problem/CF1749E
// 为了抵御怪物的袭击，Monocarp 决定利用他的院子构建一座仙人掌屏障。
// 院子可以看作 n 行m 列的网格，每个格子内最多可以种一颗仙人掌。
// 怪物可以从第一行任意一个没有仙人掌的格子进入院子，它每次能走到一个四联通相邻（即两个格子有公共边）且没有仙人掌的格子。
// 如果怪物怎么走都无法到达最后一行的某个没有仙人掌的格子，那么院子里就成功构建了一座屏障。
// !任意两个四联通相邻的格子不能同时种上仙人掌。
// 现在可能有些格子已经种上了仙人掌。
// !Monocarp 想要知道在不铲除已有仙人掌的前提下，至少要再种多少颗仙人掌才能构成一座屏障。
// n*m<=2e5.
//
// 1.每个可以种植仙人掌的格子，都和其`斜对角`可以种植仙人掌的格子连边。
// 如果已经有仙人掌，则边权为0，否则边权为1。
// 2.每一行的第一个格子，都和源点S连边，边权为0/1。
// 3.每一行的最后一个格子，都和汇点T连边，边权为0/1。
// 4.求S到T的最短路，如果最短路长度为-1，则无解。否则为最小割。
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func(grid [][]bool) (res [][]bool, ok bool) {
		ROW, COL := int32(len(grid)), int32(len(grid[0]))
		isIn := func(x, y int32) bool {
			return 0 <= x && x < ROW && 0 <= y && y < COL
		}
		idx := func(x, y int32) int32 { return COL*x + y }
		canPlant := func(x, y int32) bool { // 当前格子是否可以种植仙人掌
			if !isIn(x, y) {
				return false
			}
			for dx := int32(-1); dx <= 1; dx++ {
				for dy := int32(-1); dy <= 1; dy++ {
					if abs32(dx)+abs32(dy) == 1 {
						nx, ny := x+dx, y+dy
						if isIn(nx, ny) && grid[nx][ny] {
							return false
						}
					}
				}
			}
			return true
		}

		S, T := ROW*COL, ROW*COL+1
		adjList := make([][][2]int32, ROW*COL+2)
		for x := int32(0); x < ROW; x++ {
			for y := int32(0); y < COL; y++ {
				if !canPlant(x, y) {
					continue
				}
				curPos := idx(x, y)
				for dx := int32(-1); dx <= 1; dx++ {
					for dy := int32(-1); dy <= 1; dy++ {
						if abs32(dx)+abs32(dy) != 2 {
							continue
						}
						// 斜对角连边
						nx, ny := x+dx, y+dy
						if !canPlant(nx, ny) {
							continue
						}
						cost := int32(1)
						if grid[nx][ny] {
							cost = 0
						}
						nextPos := idx(nx, ny)
						adjList[curPos] = append(adjList[curPos], [2]int32{nextPos, cost})
					}
				}
			}
		}

		for x := int32(0); x < ROW; x++ {
			if !canPlant(x, 0) {
				continue
			}
			cost := int32(1)
			if grid[x][0] {
				cost = 0
			}
			adjList[S] = append(adjList[S], [2]int32{idx(x, 0), cost})
		}
		for x := int32(0); x < ROW; x++ {
			if !canPlant(x, COL-1) {
				continue
			}
			cost := int32(1)
			if grid[x][COL-1] {
				cost = 0
			}
			adjList[idx(x, COL-1)] = append(adjList[idx(x, COL-1)], [2]int32{T, cost})
		}

		dist, pre := bfs0132(adjList, S)
		if dist[T] == INF32 {
			return nil, false
		}

		path := RestorePath32(T, pre)
		res = make([][]bool, ROW)
		for i := int32(0); i < ROW; i++ {
			res[i] = make([]bool, COL)
		}
		for _, pos := range path {
			if pos < ROW*COL {
				x, y := pos/COL, pos%COL
				res[x][y] = true
			}
		}

		// 保留原有的仙人掌
		for x := int32(0); x < ROW; x++ {
			for y := int32(0); y < COL; y++ {
				if grid[x][y] {
					res[x][y] = true
				}
			}
		}
		return res, true
	}

	var T int32
	fmt.Fscan(in, &T)
	for i := int32(0); i < T; i++ {
		var row, col int32
		fmt.Fscan(in, &row, &col)
		grid := make([][]bool, row)
		for i := int32(0); i < row; i++ {
			var s string
			fmt.Fscan(in, &s)
			grid[i] = make([]bool, col)
			for j, c := range s {
				grid[i][j] = c == '#'
			}
		}

		res, ok := solve(grid)
		if !ok {
			fmt.Fprintln(out, "NO")
		} else {
			fmt.Fprintln(out, "YES")
			for _, row := range res {
				for _, b := range row {
					if b {
						fmt.Fprint(out, "#")
					} else {
						fmt.Fprint(out, ".")
					}
				}
				fmt.Fprintln(out)
			}
		}
	}
}

// 还原路径/dp复原.
func RestorePath32(target int32, pre []int32) []int32 {
	path := []int32{target}
	cur := target
	for pre[cur] != -1 {
		cur = pre[cur]
		path = append(path, cur)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func bfs0132(adjList [][][2]int32, start int32) (dist, pre []int32) {
	n := int32(len(adjList))
	dist = make([]int32, n)
	pre = make([]int32, n)
	for i := int32(0); i < n; i++ {
		dist[i] = INF32
		pre[i] = -1
	}
	queue := NewDeque[int32](n)
	queue.Append(start)
	dist[start] = 0
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		for _, e := range adjList[cur] {
			next, weight := e[0], e[1]
			cand := dist[cur] + weight
			if cand < dist[next] {
				dist[next] = cand
				pre[next] = cur
				if weight == 0 {
					queue.AppendLeft(next)
				} else {
					queue.Append(next)
				}
			}
		}
	}
	return
}

type Deque[D any] struct{ l, r []D }

func NewDeque[D any](cap int32) *Deque[D] {
	return &Deque[D]{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)}
}

func (q Deque[D]) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q Deque[D]) Size() int {
	return len(q.l) + len(q.r)
}

func (q *Deque[D]) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque[D]) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque[D]) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque[D]) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque[D]) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque[D]) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque[D]) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
