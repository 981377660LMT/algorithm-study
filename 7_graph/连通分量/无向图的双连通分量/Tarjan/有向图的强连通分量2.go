package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://atcoder.jp/contests/abc357/tasks/abc357_e
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	nexts := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nexts[i])
		nexts[i]--
	}

	adjList := make([][]int32, n)
	for cur, next := range nexts {
		adjList[cur] = append(adjList[cur], next)
	}

	counter := CountVisitedNodes(n, adjList)
	res := 0
	for _, v := range counter {
		res += int(v)
	}
	fmt.Fprintln(out, res)
}

// 2876. 有向图访问计数
// https://leetcode.cn/problems/count-visited-nodes-in-a-directed-graph/description/
// 给定一个有向图，对每个结点 0 <= i < n，统计从 i 出发可以访问到的结点数量。
func CountVisitedNodes(n int32, adjList [][]int32) []int32 {
	groups, belong := FindScc(n, adjList)
	dag, _ := ToDag(adjList, groups, belong, nil)

	memo := make([]int32, len(groups))
	for i := range memo {
		memo[i] = -1
	}
	var dfs func(int32) int32
	dfs = func(cur int32) int32 {
		if memo[cur] != -1 {
			return memo[cur]
		}
		res := int32(len(groups[cur]))
		for _, next := range dag[cur] {
			res += dfs(next)
		}
		memo[cur] = res
		return res
	}

	res := make([]int32, n)
	for i := int32(0); i < n; i++ {
		res[i] = dfs(belong[i])
	}
	return res
}

// tarjan 算法求有向图的 scc.
// 返回值为每个 scc 组里包含的点，每个点所在 scc 的编号(0 ~ len(groups)-1).
// !groups 按照拓扑序输出.
func FindScc(n int32, graph [][]int32) (groups [][]int32, belong []int32) {
	dfsOrder := make([]int32, n)
	dfsId := int32(0)
	stack := []int32{}
	inStack := make([]bool, n)

	var dfs func(int32) int32
	dfs = func(cur int32) int32 {
		dfsId++
		dfsOrder[cur] = dfsId
		curLow := dfsId
		stack = append(stack, cur)
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
			group := []int32{}
			for {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				inStack[top] = false
				group = append(group, top)
				if top == cur {
					break
				}
			}
			groups = append(groups, group)
		}
		return curLow
	}

	for i, order := range dfsOrder {
		if order == 0 {
			dfs(int32(i))
		}
	}

	for i, j := 0, len(groups)-1; i < j; i, j = i+1, j-1 {
		groups[i], groups[j] = groups[j], groups[i]
	}
	belong = make([]int32, n)
	for i := int32(0); i < int32(len(groups)); i++ {
		for _, v := range groups[i] {
			belong[v] = i
		}
	}
	return
}

// scc缩点成dag.
func ToDag(
	graph [][]int32, groups [][]int32, belong []int32,
	forEachEdge func(from, fromId, to, toId int32),
) (dag [][]int32, indeg []int32) {
	m := int32(len(groups))
	dag = make([][]int32, m)
	visitedEdge := map[int]struct{}{}
	indeg = make([]int32, m)
	for cur, nexts := range graph {
		curId := belong[cur]
		for _, next := range nexts {
			nextId := belong[next]
			if curId != nextId {
				hash := int(curId)*int(m) + int(nextId)
				if _, ok := visitedEdge[hash]; ok {
					continue
				}
				visitedEdge[hash] = struct{}{}
				dag[curId] = append(dag[curId], nextId)
				indeg[nextId]++
			}
			if forEachEdge != nil {
				forEachEdge(int32(cur), curId, next, nextId)
			}
		}
	}
	return
}
