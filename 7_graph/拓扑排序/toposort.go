package main

// https://leetcode.cn/problems/course-schedule/
func canFinish(numCourses int, prerequisites [][]int) bool {
	adjList := make([][]int, numCourses)
	deg := make([]int, numCourses)
	for _, e := range prerequisites {
		u, v := e[0], e[1]
		adjList[u] = append(adjList[u], v)
		deg[v]++
	}
	_, ok := TopoSort(numCourses, adjList, deg, true)
	return ok
}

func TopoSort(n int, adjList [][]int, deg []int, directed bool) (order []int, ok bool) {
	startDeg := 0
	if !directed {
		startDeg = 1
	}
	queue := make([]int, 0, n)
	for v := 0; v < n; v++ {
		if deg[v] == startDeg {
			queue = append(queue, v)
		}
	}
	visited := make([]bool, n)
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		order = append(order, cur)
		visited[cur] = true
		for _, next := range adjList[cur] {
			if visited[next] {
				continue
			}
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

// func TopoSortDfs(dag [][]int) (order []int, ok bool) {
// 	n := len(dag)
// 	visited, temp := make([]bool, n), make([]bool, n)
// 	var dfs func(int) bool
// 	dfs = func(i int) bool {
// 		if temp[i] {
// 			return false
// 		}
// 		if !visited[i] {
// 			temp[i] = true
// 			for _, v := range dag[i] {
// 				if !dfs(v) {
// 					return false
// 				}
// 			}
// 			visited[i] = true
// 			order = append(order, i)
// 			temp[i] = false
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
