// 四元环:C4
// https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
// https://blog.csdn.net/weixin_43466755/article/details/112985722

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	abc260_f()
}

func abc260_f() {
	// https://atcoder.jp/contests/abc260/tasks/abc260_f
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var S, T, M int
	fmt.Fscan(in, &S, &T, &M)

	edges := make([][2]int, M)
	for i := 0; i < M; i++ {
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
		edges[i][0]--
		edges[i][1]--
	}

	res, ok := FindC4(S+T, edges)
	if !ok {
		fmt.Fprintln(out, -1)
		return
	}

	fmt.Fprintln(out, res[0]+1, res[1]+1, res[2]+1, res[3]+1)

}

// !时间复杂度O(E*3/2)
// !四元环唯一表示两条无向边+两条有向边 u-v->w u-x->w 且 rank(u)最小rank(w)最大
// !这样才能保证一个四元环不被反复计算。
func CountC4(n int, edges [][2]int) (res int) {
	adjList1 := make([][]int, n)
	deg := make([]int, n)
	less := func(u, v int) bool { return deg[u] < deg[v] || (deg[u] == deg[v] && u < v) }

	for _, e := range edges {
		u, v := e[0], e[1]
		adjList1[u] = append(adjList1[u], v)
		adjList1[v] = append(adjList1[v], u)
		deg[u]++
		deg[v]++
	}

	adjList2 := make([][]int, n)
	for cur, nexts := range adjList1 {
		for _, next := range nexts {
			if less(cur, next) {
				adjList2[cur] = append(adjList2[cur], next)
			}
		}
	}

	count := make([]int, n)
	for cur, nexts := range adjList1 {
		for _, next1 := range nexts {
			for _, next2 := range adjList2[next1] {
				if less(cur, next2) {
					count[cur]++
					res += int(count[cur])
				}

			}

			for _, next1 := range nexts {
				for _, next2 := range adjList2[next1] {
					if less(cur, next2) {
						count[cur] = 0
					}
				}
			}
		}
	}

	return
}

// !无向图是否存在四元环 O(n^2)
func HasC4(n int, edges [][2]int) bool {
	adjList := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adjList[u] = append(adjList[u], v)
		adjList[v] = append(adjList[v], u)
	}

	check := make([][]bool, n)
	for i := range check {
		check[i] = make([]bool, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < len(adjList[i])-1; j++ {
			for k := j + 1; k < len(adjList[i]); k++ {
				if check[adjList[i][j]][adjList[i][k]] {
					return true
				}
				check[adjList[i][j]][adjList[i][k]] = true
				check[adjList[i][k]][adjList[i][j]] = true
			}
		}
	}

	return false
}

// !无向图寻找一个四元环.时间复杂度O(E*3/2).
// https://codeforces.com/problemset/problem/1468/M
func FindC4(n int, edges [][2]int) (c4 [4]int, ok bool) {
	deg := make([]int, n)
	for _, e := range edges {
		deg[e[0]]++
		deg[e[1]]++
	}
	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return deg[order[i]] > deg[order[j]] })
	rank := make([]int, n)
	for i, v := range order {
		rank[v] = i
	}

	adjList := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		rank1, rank2 := rank[u], rank[v]
		adjList[rank1] = append(adjList[rank1], rank2)
		adjList[rank2] = append(adjList[rank2], rank1)
	}
	for i := range adjList {
		sort.Sort(sort.Reverse(sort.IntSlice(adjList[i])))
	}

	pre := make([]int, n)
	for i := range pre {
		pre[i] = -1
	}

	for a := 0; a < n; a++ {
		for _, b := range adjList[a] {
			adjList[b] = adjList[b][:len(adjList[b])-1]
		}
		for _, b := range adjList[a] {
			for _, c := range adjList[b] {
				if pre[c] != -1 {
					return [4]int{order[a], order[b], order[c], order[pre[c]]}, true
				}
				pre[c] = b
			}
		}
		for _, b := range adjList[a] {
			for _, c := range adjList[b] {
				pre[c] = -1
			}
		}
	}

	return [4]int{-1, -1, -1, -1}, false
}
