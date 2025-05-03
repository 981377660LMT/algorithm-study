package main

import "fmt"

func main() {
	// Example usage of Bfs function
	start := 1
	getNext := func(v int) []int {
		return []int{v + 1, v + 2}
	}
	f := func(v int, d int) (stop bool) {
		fmt.Printf("Visiting node %d at depth %d\n", v, d)
		if v > 5 {
			stop = true
		}
		return stop
	}

	Bfs(start, getNext, f)

	// Example usage of Dfs function
	startDfs := 1
	getNextDfs := func(v int) []int {
		return []int{v + 1, v + 2}
	}
	fDfs := func(v int, d int) (stop bool) {
		fmt.Printf("Visiting node %d at depth %d\n", v, d)
		if v > 5 {
			stop = true
		}
		return stop
	}
	Dfs(startDfs, getNextDfs, fDfs)
}

func Bfs[T comparable](
	start T,
	getNext func(T) []T,
	f func(v T, d int) (stop bool),
) {
	if f(start, 0) {
		return
	}

	visited := make(map[T]struct{})
	visited[start] = struct{}{}
	queue := []T{start}

	dist := 1
	for len(queue) > 0 {
		nextQueue := []T{}
		for _, cur := range queue {
			for _, next := range getNext(cur) {
				if _, has := visited[next]; !has {
					if f(next, dist) {
						return
					}
					visited[next] = struct{}{}
					nextQueue = append(nextQueue, next)
				}
			}
		}
		queue = nextQueue
		dist++
	}
}

func Dfs[T comparable](
	start T,
	getNext func(T) []T,
	f func(v T, d int) (stop bool),
) {
	visited := make(map[T]struct{})
	var dfs func(v T, d int) bool
	dfs = func(v T, d int) bool {
		if f(v, d) {
			return true
		}
		visited[v] = struct{}{}
		for _, next := range getNext(v) {
			if _, has := visited[next]; !has {
				if dfs(next, d+1) {
					return true
				}
			}
		}
		return false
	}
	dfs(start, 0)
}
