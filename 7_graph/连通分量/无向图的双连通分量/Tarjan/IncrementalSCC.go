// IncrementalSCC
// 应用场景:
// 1. 维护各个时间点的强连通分量
// 2. 以时刻为边权建立最小生成树，求两点首次在同一个 scc 中的时刻

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	CF1989F()
	// yosupo()
}

// Simultaneous Coloring
// https://www.luogu.com.cn/problem/CF1989F
//
// 给定一个ROW * COL的棋盘，你可以把一行染成R或者一列染成B，
// 你可以同时染色多行或者多列，并如果有交叉的染色你可以控制它为 R 或 B。
// 定义一次染色的代价为选择的行和列的数量的平方(如果正好选择1个，代价为0);
// 现在有 q 个要求，每个要求为方格 (x,y) 为 R 或 B，
// 问你将达成前 i(i≤q) 个要求的代价分别为多少。
//
// 棋盘模型转二分图，如果要求染红则从左部（行）连向右部（列），否则从右部连向左部
// !答案即为大小 >1 的 SCC 大小平方和
func CF1989F() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var ROW, COL, Q int32
	fmt.Fscan(in, &ROW, &COL, &Q)

	edges := make([][2]int32, 0, Q)
	for i := int32(0); i < Q; i++ {
		var x, y int32
		var c string
		fmt.Fscan(in, &x, &y, &c)
		x--
		y--
		if c == "R" {
			edges = append(edges, [2]int32{ROW + y, x})
		} else {
			edges = append(edges, [2]int32{x, ROW + y})
		}
	}

	time := IncrementalScc(ROW+COL, edges)
	events := make([][]int32, Q+1)
	for i, t := range time {
		if t != INF32 {
			events[t] = append(events[t], int32(i))
		}
	}

	cost := func(v int) int {
		if v <= 1 {
			return 0
		}
		return v * v
	}

	res := 0
	uf := NewUnionFindArraySimple32(ROW + COL)
	for t := int32(1); t <= Q; t++ {
		eids := events[t]
		for _, eid := range eids {
			u, v := edges[eid][0], edges[eid][1]
			uf.Union(u, v, func(big, small int32) {
				size1, size2 := int(uf.GetSize(big)), int(uf.GetSize(small))
				res -= cost(size1)
				res -= cost(size2)
				res += cost(size1 + size2)
			})
		}
		fmt.Fprintln(out, res)
	}
}

// https://judge.yosupo.jp/problem/incremental_scc
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var n, m int32
	fmt.Fscan(in, &n, &m)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	edges := make([][2]int32, 0, m)
	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		edges = append(edges, [2]int32{u, v})

	}

	time := IncrementalScc(n, edges)
	groupByTime := make([][]int32, m+1)
	for i, t := range time {
		if t != INF32 {
			groupByTime[t] = append(groupByTime[t], int32(i))
		}
	}

	uf := NewUnionFindArraySimple32(n)
	res := 0
	// 每个时间点的强连通分量
	for t := int32(1); t <= m; t++ {
		eids := groupByTime[t]
		for _, eid := range eids {
			u, v := edges[eid][0], edges[eid][1]
			uf.Union(u, v, func(big, small int32) {
				res += nums[big] * nums[small]
				res %= MOD
				nums[big] += nums[small]
				nums[big] %= MOD
			})
		}
		fmt.Fprintln(out, res)
	}
}

const INF32 int32 = 1e9 + 10

// 按照顺序连接编号为 0, 1, 2, ... 的边。
// 对于每条边 i，返回它在包含的环中的最小时间.
// 如果它不在环中，返回 INF32.
// !用途：mst + path max query 可以求出两点首次在同一个 scc 中的时刻.
func IncrementalScc(n int32, directedEdges [][2]int32) (mergeTime []int32) {
	m := int32(len(directedEdges))
	mergeTime = make([]int32, m)
	for i := range mergeTime {
		mergeTime[i] = INF32
	}
	data := make([][3]int32, 0, m) // (i,from,to)
	for i, e := range directedEdges {
		data = append(data, [3]int32{int32(i), e[0], e[1]})
	}

	newId := make([]int32, n)
	for i := range newId {
		newId[i] = -1
	}

	// L 时刻不在环中，R 时刻在环中
	var dfs func([][3]int32, int32, int32)
	dfs = func(data [][3]int32, L, R int32) {
		if len(data) == 0 || R == L+1 {
			return
		}
		mid := (L + R) >> 1
		n := int32(0)
		for j := range data {
			a, b := data[j][1], data[j][2]
			if newId[a] == -1 {
				newId[a] = n
				n++
			}
			if newId[b] == -1 {
				newId[b] = n
				n++
			}
		}
		newGraph := make([][]int32, n)
		for j := range data {
			i, a, b := data[j][0], data[j][1], data[j][2]
			a, b = newId[a], newId[b]
			if i < mid {
				newGraph[a] = append(newGraph[a], b)
			}
		}

		belong := StronglyConnectedComponent(n, newGraph)

		var data1, data2 [][3]int32
		for j := range data {
			i, a, b := data[j][0], data[j][1], data[j][2]
			a, b = newId[a], newId[b]
			if i < mid {
				if belong[a] == belong[b] {
					if mid < mergeTime[i] {
						mergeTime[i] = mid
					}
					data1 = append(data1, [3]int32{i, a, b})
				} else {
					data2 = append(data2, [3]int32{i, belong[a], belong[b]})
				}
			} else {
				data2 = append(data2, [3]int32{i, belong[a], belong[b]})
			}
		}
		for j := range data {
			a, b := data[j][1], data[j][2]
			newId[a], newId[b] = -1, -1
		}
		dfs(data1, L, mid)
		dfs(data2, mid, R)
	}

	dfs(data, 0, m+1)
	return
}

func StronglyConnectedComponent(n int32, graph [][]int32) (belong []int32) {
	dfsOrder := make([]int32, n)
	dfsId := int32(0)
	stack := make([]int32, n)
	stackPtr := int32(0)
	inStack := make([]bool, n)
	belong = make([]int32, n)
	compId := int32(0)

	var dfs func(int32) int32
	dfs = func(cur int32) int32 {
		dfsId++
		dfsOrder[cur] = dfsId
		curLow := dfsId
		stack[stackPtr] = cur
		stackPtr++
		inStack[cur] = true
		for _, next := range graph[cur] {
			if dfsOrder[next] == 0 {
				nextLow := dfs(next)
				if nextLow < curLow {
					curLow = nextLow
				}
			} else if inStack[next] && dfsOrder[next] < curLow {
				curLow = dfsOrder[next]
			}
		}
		if dfsOrder[cur] == curLow {
			for {
				top := stack[stackPtr-1]
				stackPtr--
				inStack[top] = false
				belong[top] = compId
				if top == cur {
					break
				}
			}
			compId++
		}
		return curLow
	}

	for i, order := range dfsOrder {
		if order == 0 {
			dfs(int32(i))
		}
	}
	return
}

type UnionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple32(n int32) *UnionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple32) Union(key1, key2 int32, beforeMerge func(big, small int32)) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	if beforeMerge != nil {
		beforeMerge(root1, root2)
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = root1
	u.Part--
	return true
}

func (u *UnionFindArraySimple32) Find(key int32) int32 {
	if u.data[key] < 0 {
		return key
	}
	u.data[key] = u.Find(u.data[key])
	return u.data[key]
}

func (u *UnionFindArraySimple32) GetSize(key int32) int32 {
	return -u.data[u.Find(key)]
}
