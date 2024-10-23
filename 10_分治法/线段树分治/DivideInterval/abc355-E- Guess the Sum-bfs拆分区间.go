// 交互题
// abc355-E- Guess the Sum-bfs拆分区间
// 给定任意区间[l,r)，求将区间拆分成若干个形如[2^i*j,2^i*(j+1))的区间.
// !操作满足阿贝尔群.
//
// !注意：树状数组的拆分不是最少.
// !O(logn)拆分方法

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	abc355_e()
}

// 交互题
// abc355-E- Guess the Sum-bfs拆分区间
// 给定任意区间[l,r)，求将区间拆分成若干个形如[2^i*j,2^i*(j+1))的区间.
// !操作满足阿贝尔群.
//
// !注意每次输出后调用out.Flush()刷新输出缓冲区.
func abc355_e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	query := func(start, end int32) int {
		i, j := Format(start, end)
		fmt.Println("?", i, j)
		out.Flush()
		var res int
		fmt.Fscan(in, &res)
		return res
	}

	output := func(res int) {
		fmt.Println("!", res)
		out.Flush()
	}

	var LOG, L, R int32
	fmt.Fscan(in, &LOG, &L, &R)

	n := int32(1 << LOG)
	adjList := make([][]int32, n+1)
	for k := int32(0); k < LOG+1; k++ {
		for l := int32(0); l < n; l += 1 << k {
			r := l + (1 << k)
			adjList[l] = append(adjList[l], r)
			adjList[r] = append(adjList[r], l)
		}
	}

	_, pre := bfs(adjList, L)
	path := restorePath(pre, R+1)

	res := 0
	for i := 0; i < len(path)-1; i++ {
		start, end := path[i], path[i+1]
		sign := 1
		if start > end {
			start, end = end, start
			sign = -1
		}
		res += query(start, end) * sign
	}

	res = (res%100 + 100) % 100
	output(res)
}

func Format(start, end int32) (i, j int32) {
	n := end - start
	i = int32(bits.Len(uint(n - 1)))
	j = start >> i
	return
}

const INF32 int32 = 1e9 + 10

func bfs(adjList [][]int32, start int32) (dist, pre []int32) {
	n := int32(len(adjList))
	dist = make([]int32, n)
	pre = make([]int32, n)
	for i := int32(0); i < n; i++ {
		dist[i] = INF32
		pre[i] = -1
	}
	queue := []int32{start}
	dist[start] = 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next := range adjList[cur] {
			if dist[next] == INF32 {
				dist[next] = dist[cur] + 1
				pre[next] = cur
				queue = append(queue, next)
			}
		}
	}
	return
}

// 还原路径/dp复原.
func restorePath(pre []int32, target int32) []int32 {
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
