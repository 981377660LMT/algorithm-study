// https://ei1333.github.io/library/graph/others/chromatic-number.hpp
// 一般图的彩色数 O(n*2^n)
// グラフの彩色数を求める. 彩色数とは,
// !隣接する頂点が異なる色となるように彩色するために必要な最小色数である.
// !(每个点的邻接点不能有相同的颜色)

// あるグラフが k彩色可能であることと, k個の独立集合で被覆できることは必要十分条件である.
// つまり独立集合から
// k個の頂点を選んで被覆する方法の総数が求まれば良い.
// これはbit DPと包除原理を用いて効率的に求められる.

// 使い方
// chromatic_number(n, edges) で, n頂点のグラフに対する彩色数を求める.

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// https://judge.yosupo.jp/problem/chromatic_number
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][]int, 0, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		edges = append(edges, []int{u, v})
	}

	fmt.Fprintln(out, ChromaticNnumber(n, edges))
}

func ChromaticNnumber(n int, edges [][]int) int {
	matrix := make([][]bool, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]bool, n)
	}
	for _, e := range edges {
		u, v := e[0], e[1]
		matrix[u][v] = true
		matrix[v][u] = true
	}

	es := make([]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if matrix[i][j] {
				es[i] |= 1 << j
			}
		}
	}

	ind := make([]int, 1<<n)
	ind[0] = 1
	for S := 1; S < 1<<n; S++ {
		u := bits.TrailingZeros(uint(S))
		ind[S] = ind[S^(1<<u)] + ind[(S^(1<<u))&^es[u]]
	}

	cnt := make([]int, (1<<n)+1)
	for i := 0; i < 1<<n; i++ {
		isOdd := bits.OnesCount(uint(i)) & 1
		cnt[ind[i]] += 1 - (isOdd << 1) // idOdd ? -1 : 1
	}

	hist := make([][2]int, 0)
	for i := 1; i <= 1<<n; i++ {
		if cnt[i] != 0 {
			hist = append(hist, [2]int{i, cnt[i]})
		}
	}

	mods := [3]int{1000000007, 1000000011, 1000000021}

	res := n
	for k := 0; k < 3; k++ {
		buf := make([][2]int, len(hist))
		copy(buf, hist)
		for c := 1; c < res; c++ {
			sum := 0
			for i := range buf {
				buf[i][1] = buf[i][1] * buf[i][0] % mods[k]
				sum += buf[i][1]
			}
			if sum%mods[k] != 0 {
				res = c
			}
		}
	}

	return res
}
