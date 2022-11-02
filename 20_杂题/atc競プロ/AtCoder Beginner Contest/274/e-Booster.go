package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	street := make([]struct{ x, y int }, n)
	for i := range street {
		fmt.Fscan(in, &street[i].x, &street[i].y)
	}

	box := make([]struct{ x, y int }, m)
	for i := range box {
		fmt.Fscan(in, &box[i].x, &box[i].y)
	}

	// points = street + box + [(0, 0)]
	points := append([]struct{ x, y int }{}, street...)
	points = append(points, box...)
	points = append(points, struct{ x, y int }{0, 0})

	dist := make([][]float64, n+m+1)
	for i := range dist {
		x1, y1 := points[i].x, points[i].y
		dist[i] = make([]float64, n+m+1)
		for j := range dist[i] {
			x2, y2 := points[j].x, points[j].y
			dist[i][j] = math.Sqrt(float64((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)))
		}
	}

	memo := make([][]float64, n+m+1) // (index,visited)
	for i := range memo {
		memo[i] = make([]float64, 1<<(n+m+1))
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}

	target := (1 << n) - 1

	var dfs func(int, int) float64
	dfs = func(cur, visited int) float64 {
		if memo[cur][visited] != -1 {
			return memo[cur][visited]
		}

		if visited&target == target {
			speed := float64(int(1) << bits.OnesCount(uint(visited>>n)))
			memo[cur][visited] = dist[cur][n+m] / speed
			return memo[cur][visited]
		}

		speed := float64(int(1) << bits.OnesCount(uint(visited>>n)))
		res := math.MaxFloat64
		for next := 0; next < n+m; next++ {
			if visited&(1<<next) != 0 {
				continue
			}
			cost := dist[cur][next] / speed
			cand := cost + dfs(next, visited|(1<<next))
			if cand < res {
				res = cand
			}
		}

		memo[cur][visited] = res
		return res
	}

	res := dfs(n+m, 0)
	fmt.Fprintln(out, res)
}
