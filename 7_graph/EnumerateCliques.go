// 枚举无向图的团(クリーク) & 最大团
// 在一个无向图中找出一个点数最多的完全图(每对顶点之间都恰连有一条边的图)。
// 时间复杂度O(n*2^sqrt(2*m))
// n,m<=200
// https://ei1333.github.io/library/graph/others/enumerate-cliques.hpp

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func maxClique(n int, edges [][]int) int {
	res := 0
	enumerateClique(n, edges, func(clique []int) {
		if len(clique) > res {
			res = len(clique)
		}
	})
	return res
}

func enumerateClique(n int, edges [][]int, cb func(clique []int)) {
	deg := make([]int, n)
	adjMatrix := make([][]bool, n)
	for i := range adjMatrix {
		adjMatrix[i] = make([]bool, n)
	}

	for _, e := range edges {
		u, v := e[0], e[1]
		deg[u]++
		deg[v]++
		adjMatrix[u][v] = true
		adjMatrix[v][u] = true
	}

	m := len(edges)
	lim := int(math.Sqrt(float64(m * 2)))

	// dont mutate rem
	findClique := func(rem []int, last bool) {
		tmp := 0
		if last {
			tmp = 1
		}

		neighbor := make([]int, len(rem)-tmp)
		for i := 0; i < len(neighbor); i++ {
			for j := 0; j < len(neighbor); j++ {
				if i != j && !adjMatrix[rem[i]][rem[j]] {
					neighbor[i] |= 1 << uint(j)
				}
			}
		}

		for i := 1 - tmp; i < 1<<uint(len(neighbor)); i++ {
			ok := true
			for j := 0; j < len(neighbor); j++ {
				if ((i>>uint(j))&1 == 1) && ((i & neighbor[j]) != 0) {
					ok = false
					break
				}
			}
			if ok {
				clique := []int{}
				if last {
					clique = append(clique, rem[len(rem)-1])
				}
				for j := 0; j < len(neighbor); j++ {
					if (i>>uint(j))&1 == 1 {
						clique = append(clique, rem[j])
					}
				}
				cb(clique)
			}
		}
	}

	used := make([]bool, n)
	queue := []int{}
	for i := 0; i < n; i++ {
		if deg[i] < lim {
			used[i] = true
			queue = append(queue, i)
		}
	}

	for len(queue) > 0 {
		index := queue[0]
		queue = queue[1:]
		rem := []int{}
		for k := 0; k < n; k++ {
			if adjMatrix[index][k] {
				rem = append(rem, k)
			}
		}
		rem = append(rem, index)
		findClique(rem, true)
		used[index] = true
		for k := 0; k < n; k++ {
			if adjMatrix[index][k] {
				adjMatrix[index][k] = false
				adjMatrix[k][index] = false
				deg[k]--
				if !used[k] && deg[k] < lim {
					used[k] = true
					queue = append(queue, k)
				}
			}
		}
	}

	rem := []int{}
	for i := 0; i < n; i++ {
		if !used[i] {
			rem = append(rem, i)
		}
	}
	findClique(rem, false)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][]int, 0, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		edges = append(edges, []int{u, v})
	}

	fmt.Fprintln(out, maxClique(n, edges))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
