package main

// https://leetcode.cn/problems/course-schedule/
func canFinish(numCourses int, prerequisites [][]int) bool {
	adjList := make([][]int, numCourses)
	for _, e := range prerequisites {
		u, v := e[0], e[1]
		adjList[u] = append(adjList[u], v)
	}
	_, ok := TopoSort(numCourses, adjList, true)
	return ok
}

func canFinish2(numCourses int, prerequisites [][]int) bool {
	vertices := make([]int, numCourses)
	for i := range vertices {
		vertices[i] = i
	}
	edges := make([][2]int, len(prerequisites))
	for i, e := range prerequisites {
		edges[i] = [2]int{e[0], e[1]}
	}
	_, ok := TopoSortMap(vertices, edges, true)
	return ok
}

// 拓扑排序环检测.
func HasCycle(n int, adjList [][]int, directed bool) bool {
	deg := make([]int, n)
	startDeg := 0
	if directed {
		for _, adj := range adjList {
			for _, j := range adj {
				deg[j]++
			}
		}
	} else {
		for i, adj := range adjList {
			deg[i] = len(adj)
		}
		startDeg = 1
	}

	queue := []int{}
	for v := 0; v < n; v++ {
		if deg[v] == startDeg {
			queue = append(queue, v)
		}
	}
	count := 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		count++
		for _, next := range adjList[cur] {
			deg[next]--
			if deg[next] == startDeg {
				queue = append(queue, next)
			}
		}
	}

	return count < n
}

// 拓扑排序结果.
func TopoSort(n int, adjList [][]int, directed bool) (order []int, ok bool) {
	deg := make([]int, n)
	startDeg := 0
	if directed {
		for _, adj := range adjList {
			for _, j := range adj {
				deg[j]++
			}
		}
	} else {
		for i, adj := range adjList {
			deg[i] = len(adj)
		}
		startDeg = 1
	}

	queue := []int{}
	for v := 0; v < n; v++ {
		if deg[v] == startDeg {
			queue = append(queue, v)
		}
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		order = append(order, cur)
		for _, next := range adjList[cur] {
			deg[next]--
			if deg[next] == startDeg {
				queue = append(queue, next)
			}
		}
	}

	if len(order) < n {
		return nil, false
	}
	return order, true
}

func TopoSortMap[T comparable](vertices []T, edges [][2]T, directed bool) (order []T, ok bool) {
	verticesSet := make(map[T]struct{}, len(vertices))
	for _, v := range vertices {
		verticesSet[v] = struct{}{}
	}
	deg, adjMap := make(map[T]int32, len(verticesSet)), make(map[T][]T, len(verticesSet))
	for v := range verticesSet {
		deg[v] = 0
		adjMap[v] = []T{}
	}

	add := func(from, to T) {
		adjMap[from] = append(adjMap[from], to)
		deg[to]++
	}
	if directed {
		for _, e := range edges {
			add(e[0], e[1])
		}
	} else {
		for _, e := range edges {
			add(e[0], e[1])
			add(e[1], e[0])
		}
	}

	startDeg := int32(0)
	if !directed {
		startDeg = 1
	}
	queue := []T{}
	for v, d := range deg {
		if d == startDeg {
			queue = append(queue, v)
		}
	}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		order = append(order, cur)
		for _, next := range adjMap[cur] {
			deg[next]--
			if deg[next] == startDeg {
				queue = append(queue, next)
			}
		}
	}

	if len(order) < len(deg) {
		return nil, false
	}
	return order, true
}

// func TopoSortDfs(dag [][]int) (order []int, ok bool) {
// 	n := len(dag)
// 	visited, onPath := make([]bool, n), make([]bool, n)
// 	var dfs func(int) bool
// 	dfs = func(i int) bool {
// 		if onPath[i] {
// 			return false
// 		}
// 		if !visited[i] {
// 			onPath[i] = true
// 			for _, v := range dag[i] {
// 				if !dfs(v) {
// 					return false
// 				}
// 			}
// 			visited[i] = true
// 			order = append(order, i)
// 			onPath[i] = false
// 		}
// 		return true
// 	}

// 	for i := 0; i < n; i++ {
// 		if !visited[i] {
// 			if !dfs(i) {
// 				return nil, false
// 			}
// 		}
// 	}

// 	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
// 		order[i], order[j] = order[j], order[i]
// 	}
// 	return order, true
// }
