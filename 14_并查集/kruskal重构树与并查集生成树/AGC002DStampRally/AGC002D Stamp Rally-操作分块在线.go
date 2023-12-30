// AGC002D Stamp Rally-操作分块在线
// https://hoikoro.hatenablog.com/entry/2017/12/14/040750

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	graph := make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([][3]int, q)
	for i := 0; i < q; i++ {
		var x, y, z int
		fmt.Fscan(in, &x, &y, &z)
		x--
		y--
		queries[i] = [3]int{x, y, z}
	}

	res := StampRally(graph, queries)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// https://atcoder.jp/contests/agc002/tasks/agc002_d
// 一张连通图，q 次询问从两个点 x 和 y 出发，
// 希望经过的点数量等于 z（每个点可以重复经过，但是重复经过只计算一次）
// 求经过的边最大编号最小是多少。
//
// 很容易想到离线加边，判断每一个询问是否满足条件，但必须每次检查每一个询问，时间复杂度O(mqlogn)
// 上述过程每一次检验非常浪费时间，考虑分块，每插入O(sqrt(m))次检验一次，，对于每一个询问倒着扫一遍用带撤销并查集断一下就行了。
// 时间复杂度O(qsqrt(m)logn).
func StampRally(graph [][]int, queries [][3]int) []int {}
