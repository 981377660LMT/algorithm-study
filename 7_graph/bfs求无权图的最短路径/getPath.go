package main

// 求from到to的一条路径.
func GetPath(graph [][]int32, from, to int32) []int32 {
	n := int32(len(graph))
	visited := make([]bool, n)
	visited[from] = true
	stack := []int32{from}
	pre := make([]int32, n)
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		if cur == to {
			break
		}
		stack = stack[:len(stack)-1]
		for _, next := range graph[cur] {
			if !visited[next] {
				visited[next] = true
				pre[next] = cur
				stack = append(stack, next)
			}
		}
	}
	if !visited[to] {
		return nil
	}
	path := []int32{to}
	for v := to; v != from; v = pre[v] {
		path = append(path, pre[v])
	}
	for i := 0; i < len(path)/2; i++ {
		path[i], path[len(path)-1-i] = path[len(path)-1-i], path[i]
	}
	return path
}
