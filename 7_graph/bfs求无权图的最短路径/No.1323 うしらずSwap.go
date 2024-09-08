package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const INF int = 1e18
const INF32 int32 = 1e9 + 10

func main() {
	yuki1323()
	// abc301_e()
}

// No.1323 うしらずSwap (交换位置，移形换位，分类讨论)
// https://yukicoder.me/problems/no/1323
// 给定一个带障碍的网格图和两个人的位置.
// 求两个人交换位置的最短路径总和.注意两个人不能同时移动到同一个位置.
// 如果不能交换位置，输出-1.
//
// 预处理出两个人到所有点的最短路.
// 1.最短路径不唯一 -> 2*D
// 2.最短路径唯一 -> 2*D+2
// 3.经过度数 >= 3 的点 -> 2*(dist1[x][y]+dist2[x][y])+4
// 4.两条边不相交的路径 -> D+dist[gx][gy]
func yuki1323() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var H, W, sx, sy, gx, gy int32
	fmt.Fscan(in, &H, &W, &sx, &sy, &gx, &gy)
	sx, sy, gx, gy = sx-1, sy-1, gx-1, gy-1
	G := make([]string, H)
	for i := int32(0); i < H; i++ {
		fmt.Fscan(in, &G[i])
	}

	dist1 := GridBfs(H, W, sx, sy, 4, func(x, y int32) bool { return G[x][y] == '#' })
	dist2 := GridBfs(H, W, gx, gy, 4, func(x, y int32) bool { return G[x][y] == '#' })
	D := dist1[gx][gy]
	if D == INF32 {
		fmt.Fprintln(out, -1)
		return
	}

	// 点是否在最短路径上.
	onPath := func(x, y int32) bool {
		d1, d2 := dist1[x][y], dist2[x][y]
		return d1 < INF32 && d2 < INF32 && d1+d2 == D
	}

	{
		// !最短路径不唯一 -> 2*D
		count := int32(0)
		for x := int32(0); x < H; x++ {
			for y := int32(0); y < W; y++ {
				if onPath(x, y) {
					count++
				}
			}
		}
		if count >= D+2 {
			fmt.Fprintln(out, 2*D)
			return
		}
	}

	res := INF32
	// 经过度数 >= 3 的点.
	{
		for x := int32(0); x < H; x++ {
			for y := int32(0); y < W; y++ {
				if G[x][y] == '#' {
					continue
				}
				deg := int32(0)
				for i := int32(0); i < 4; i++ {
					if G[x+dir8[i][0]][y+dir8[i][1]] != '#' {
						deg++
					}
				}
				if deg >= 3 {
					mid := onPath(x, y)
					if x == sx && y == sy {
						mid = false
					}
					if x == gx && y == gy {
						mid = false
					}
					if mid {
						fmt.Fprintln(out, 2*D+2)
						return
					} else if dist1[x][y] < INF32 {
						res = min32(res, 2*(dist1[x][y]+dist2[x][y])+4)
					}
				}
			}
		}
	}

	// 两条边不相交的路径.
	newG := make([][]byte, H)
	for i := range newG {
		newG[i] = []byte(G[i])
	}
	for x := int32(0); x < H; x++ {
		for y := int32(0); y < W; y++ {
			if onPath(x, y) {
				newG[x][y] = '#'
			}
		}
	}
	newG[sx][sy], newG[gx][gy] = '.', '.'
	dist := make([][]int32, H)
	for i := range dist {
		dist[i] = make([]int32, W)
		for j := range dist[i] {
			dist[i][j] = INF32
		}
	}
	dist[sx][sy] = 0
	queue := [][2]int32{{sx, sy}}
	for len(queue) > 0 {
		x, y := queue[0][0], queue[0][1]
		queue = queue[1:]
		for i := int32(0); i < 4; i++ {
			nx, ny := x+dir8[i][0], y+dir8[i][1]
			if newG[nx][ny] == '#' {
				continue
			}
			if x == sx && y == sy && nx == gx && ny == gy {
				continue
			}
			if dist[nx][ny] > dist[x][y]+1 {
				dist[nx][ny] = dist[x][y] + 1
				queue = append(queue, [2]int32{nx, ny})
			}
		}
	}
	res = min32(res, D+dist[gx][gy])

	if res == INF32 {
		fmt.Fprintln(out, -1)
		return
	}
	fmt.Fprintln(out, res)
}

// E - Pac-Takahashi
// https://atcoder.jp/contests/abc301/tasks/abc301_e
// 最短路+状压dp
//
// bfs 预处理出所有点之间(起点、关键点、终点)的最短路
// 状态压缩dp，dp[s][i] 表示经过 s 中的点，最后到达 i 的最短路
func abc301_e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var H, W, T int32
	fmt.Fscan(in, &H, &W, &T)
	A := make([]string, H)
	for i := int32(0); i < H; i++ {
		fmt.Fscan(in, &A[i])
	}

	pos := [][2]int32{}
	for x := int32(0); x < H; x++ {
		for y := int32(0); y < W; y++ {
			if A[x][y] == 'S' {
				pos = append(pos, [2]int32{x, y})
			}
		}
	}
	for x := int32(0); x < H; x++ {
		for y := int32(0); y < W; y++ {
			if A[x][y] == 'o' {
				pos = append(pos, [2]int32{x, y})
			}
		}
	}
	for x := int32(0); x < H; x++ {
		for y := int32(0); y < W; y++ {
			if A[x][y] == 'G' {
				pos = append(pos, [2]int32{x, y})
			}
		}
	}

	N := int32(len(pos))
	dist := make([][]int32, N)
	for i := range dist {
		dist[i] = make([]int32, N)
	}
	for i := int32(0); i < N; i++ {
		sx, sy := pos[i][0], pos[i][1]
		mat := GridBfs(H, W, sx, sy, 4, func(x, y int32) bool {
			return A[x][y] == '#'
		})
		for j := int32(0); j < N; j++ {
			a, b := pos[j][0], pos[j][1]
			dist[i][j] = mat[a][b]
		}
	}

	dp := make([][]int, 1<<N)
	for i := range dp {
		dp[i] = make([]int, N)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	dp[1][0] = 0

	for s := int32(1); s < 1<<N; s++ {
		for a := int32(0); a < N; a++ {
			if dp[s][a] == INF {
				continue
			}
			for b := int32(0); b < N; b++ {
				t := s | 1<<b
				if s == t {
					continue
				}
				if tmp := dp[s][a] + int(dist[a][b]); tmp < dp[t][b] {
					dp[t][b] = tmp
				}
			}
		}
	}

	if dist[0][N-1] > T {
		fmt.Fprintln(out, -1)
		return
	}

	res := int32(0)
	for s := int32(1); s < 1<<N; s++ {
		if dp[s][N-1] > int(T) {
			continue
		}
		n := int32(bits.OnesCount32(uint32(s)) - 2)
		if n > res {
			res = n
		}
	}

	fmt.Fprintln(out, res)
}

var dir8 = [][]int32{{1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}

// 网格最短路.
func GridBfs(
	row, col int32, sx, sy int32, dmax int32,
	isWall func(int32, int32) bool,
) [][]int32 {
	isIn := func(x, y int32) bool {
		if x < 0 || row <= x {
			return false
		}
		if y < 0 || col <= y {
			return false
		}
		return !isWall(x, y)
	}

	dist := make([][]int32, row)
	for i := range dist {
		dist[i] = make([]int32, col)
		for j := range dist[i] {
			dist[i][j] = INF32
		}
	}
	queue := [][2]int32{{sx, sy}}
	add := func(x, y, d int32) {
		if !isIn(x, y) {
			return
		}
		if dist[x][y] > d {
			dist[x][y] = d
			queue = append(queue, [2]int32{x, y})
		}
	}
	add(sx, sy, 0)

	for len(queue) > 0 {
		x, y := queue[0][0], queue[0][1]
		queue = queue[1:]
		for i := int32(0); i < dmax; i++ {
			nx, ny := x+dir8[i][0], y+dir8[i][1]
			add(nx, ny, dist[x][y]+1)
		}
	}
	return dist
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
