// 给定一张 n 个点、m 条边的有向图，每个点都有一个权值 f[i]，每条边都有一个权值 t[i]。
// !求图中的一个环，使“环上各点的权值之和”除以“环上各边的权值之和”最大。
// 输出这个最大值。
// 注意：数据保证至少存在一个环。
// n<=1000 m<=5000

// 01分数规划 二分法
// https://www.acwing.com/solution/content/40640/
// !求 ∑f/∑t 的最大值
// ! ∑f/∑t > mid 即 ∑(mid*ti - fi) < 0  即存在负环
// ! mid∗ti−fi 为边权

package main

import (
	"bufio"
	"fmt"
	"os"
)

const EPS float64 = 1e-8

func 分数规划01(n int, edges [][3]int, values []int) float64 {
	graph := make([][][2]int, n)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		graph[u] = append(graph[u], [2]int{v, w})
	}

	check := func(mid float64) bool {
		dist := make([]float64, n)
		queue := make([]int, 0, n)
		inQueue := make([]bool, n)
		relaxedCount := make([]int, n)
		for i := 0; i < n; i++ {
			queue = append(queue, i)
			inQueue[i] = true
			relaxedCount[i] = 1
		}

		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			inQueue[cur] = false
			for _, e := range graph[cur] {
				next, weight := e[0], e[1]
				cand := float64(dist[cur]) - float64(values[cur]) + mid*float64(weight) // 新的边权为  mid*weight - values[cur]
				if cand < dist[next] {                                                  // !如果要正环这里需要改成 >
					dist[next] = cand
					if !inQueue[next] {
						relaxedCount[next]++
						if relaxedCount[next] >= n+1 { // +1是虚拟源点
							return true
						}
						inQueue[next] = true
						queue = append(queue, next)
					}
				}
			}
		}

		return false
	}

	left, right := float64(0), 1e9
	for left <= right {
		mid := float64((left + right) / 2)
		if check(mid) {
			left = mid + EPS
		} else {
			right = mid - EPS
		}
	}
	return right
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	values := make([]int, n)
	for i := range values {
		fmt.Fscan(in, &values[i])
	}
	edges := make([][3]int, m)
	for i := range edges {
		fmt.Fscan(in, &edges[i][0], &edges[i][1], &edges[i][2])
		edges[i][0]--
		edges[i][1]--
	}

	res := 分数规划01(n, edges, values)

	// 保留两位小数
	fmt.Fprintf(out, "%.2f", res)
}
