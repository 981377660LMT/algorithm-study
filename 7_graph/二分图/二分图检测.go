// 二分图检测

package main

func IsBipartite(n int, graph [][]int) (colors []int, ok bool) {
	colors = make([]int, n)
	for i := range colors {
		colors[i] = -1
	}
	bfs := func(start int) bool {
		colors[start] = 0
		queue := []int{start}
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			for _, next := range graph[cur] {
				if colors[next] == -1 {
					colors[next] = colors[cur] ^ 1
					queue = append(queue, next)
				} else if colors[next] == colors[cur] {
					return false
				}
			}
		}
		return true
	}

	for i := range colors {
		if colors[i] == -1 && !bfs(i) {
			return nil, false
		}
	}
	return colors, true
}

// 扩展域并查集判断二分图.
//  配合 `OfflineDynamicConnectivity` 可以支持动态修改边的二分图判断.
//  https://www.luogu.com.cn/problem/solution/P5787
func IsBipartite2(n int, adjList [][]int) bool {
	uf := NewUnionFindArray(n * 2)
	for cur := 0; cur < n; cur++ {
		for _, next := range adjList[cur] {
			uf.Union(cur, next+n)
			uf.Union(cur+n, next)
			if uf.Find(cur) == uf.Find(next) {
				return false
			}
		}
	}
	return true
}

type UnionFindArray struct {
	data []int
}

func NewUnionFindArray(n int) *UnionFindArray {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArray{data: data}
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1, root2 = root2, root1
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}
