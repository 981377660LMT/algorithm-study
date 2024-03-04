// 多组查询visited技巧
// 一个技巧是，用visitedTime数组代替布尔类型的visited数组
// 带时间戳的visited数组，适用与多组查询的场景

package main

func good(n int32, operations [][2]int32) {
	visitedTime := make([]int32, n)
	for i := range visitedTime {
		visitedTime[i] = -1
	}
	for i, op := range operations {
		index := op[0]
		// mutate ...
		if visitedTime[index] < int32(i) {
			visitedTime[index] = int32(i)
			// query ...
		}
	}
}

func bad(n int32, operations [][2]int32) {
	visited := make([]bool, n)
	for _, op := range operations {
		mutate := []int32{}
		index := op[0]
		// mutate ...
		if !visited[index] {
			visited[index] = true
			mutate = append(mutate, index)
			// query ...
		}
		for _, v := range mutate {
			visited[v] = false
		}
	}
}
