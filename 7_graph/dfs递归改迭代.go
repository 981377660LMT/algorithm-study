// dfs递归转迭代(dfs树搜索)
// https://nachiavivias.github.io/cp-library/column/2022/01.html
// 1. 不太适合的一种方法是自己实现递归函数生成的栈
// 2. 合适的方法是：
//   - 定义状态
//   - 建图(也可以隐式建图，邻接表用迭代器表示)
//   - 记录dfs树中每个点的父节点、邻接表的迭代器索引

package main

import "fmt"

func main() {
	demo()
}

func demo() {
	// dfs求子树大小
	//   0
	//  /  \
	// 1    3
	// |   / \
	// 2  4   5
	//    |   |
	//    6   7
	graph := [][]int32{
		{1, 3},
		{0, 2},
		{1},
		{0, 4, 5},
		{3, 6},
		{3, 7},
		{4},
		{5},
	}

	subSize := make([]int32, len(graph))
	Dfs(graph, 0,
		func(v int32) {
			fmt.Println("down", v)
			subSize[v] = 1
		},
		func(child, parent int32) {
			fmt.Println("up", child, parent)
			subSize[parent] += subSize[child]
		},
	)
	fmt.Println(subSize)
}

func Dfs(graph [][]int32, start int32, down func(v int32), up func(child, parent int32)) {
	n := int32(len(graph))
	parent := make([]int32, n) // DFS树每个点的父节点
	for i := range parent {
		parent[i] = -1
	}
	iter := make([]int32, n) // 每个点的邻接表的迭代器
	cur := start             // 探索中的点
	parent[cur] = -2
	for cur >= 0 {
		if iter[cur] == 0 {
			// !首次访问
			down(cur)
		}

		if iter[cur] == int32(len(graph[cur])) {
			// !回溯
			p := parent[cur]
			if p >= 0 {
				up(cur, p)
			}
			cur = p
			continue
		}

		next := graph[cur][iter[cur]]
		iter[cur]++
		if parent[next] != -1 {
			// 返回边 cur->next
			continue
		}

		// !DFS树边 cur->next
		parent[next] = cur
		cur = next
	}
}
