// https://www.luogu.com.cn/problem/CF1416D
// https://www.luogu.com.cn/blog/wlx050125/solution-cf1416d
// TODO
package main

import (
	"bufio"
	"fmt"
	"os"
)

// 给定一个无向图，每个节点有一个权值，且每个节点的权值都不相同.
// 进行q次操作.
// 0 v : 查询与v相连的所有节点中权值最大的节点u的权值，然后将u这个点的权值变为0.
// 1 ei : 将第ei条边删除.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)

	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}

	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
		edges[i][0]--
		edges[i][1]--
	}

	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1])
		queries[i][0]--
		queries[i][1]--
	}
}
