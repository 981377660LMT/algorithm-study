package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(RandomTree(6))
	fmt.Println(RandomTree32(6))
}

func RandomTree(n int) [][]int {
	if n <= 2 {
		g := make([][]int, n)
		if n == 2 {
			g[0] = append(g[0], 1)
			g[1] = append(g[1], 0)
		}
		return g
	}
	prufer := make([]int, n-2)
	for i := 0; i < n-2; i++ {
		prufer[i] = 1 + rand.Intn(n)
	}
	deg := make([]int, n+1)
	for _, p := range prufer {
		deg[p]++
	}
	prufer = append(prufer, n)
	parents := make([]int, n+1)
	i, j := 0, 1
	for i < n-1 {
		for deg[j] > 0 {
			j++
		}
		parents[j] = prufer[i]
		for i < n-2 {
			p := prufer[i]
			deg[p]--
			if p > j || deg[p] > 0 {
				break
			}
			parents[p] = prufer[i+1]
			i++
		}
		i++
		j++
	}
	parents = parents[1:]
	tree := make([][]int, n)
	for i := 1; i < n; i++ {
		p := parents[i-1]
		tree[i-1] = append(tree[i-1], p-1)
		tree[p-1] = append(tree[p-1], i-1)
	}
	return tree
}

func RandomTree32(n int32) [][]int32 {
	if n <= 2 {
		g := make([][]int32, n)
		if n == 2 {
			g[0] = append(g[0], 1)
			g[1] = append(g[1], 0)
		}
		return g
	}
	prufer := make([]int32, n-2)
	for i := int32(0); i < n-2; i++ {
		prufer[i] = 1 + int32(rand.Intn(int(n)))
	}
	deg := make([]int32, n+1)
	for _, p := range prufer {
		deg[p]++
	}
	prufer = append(prufer, n)
	parents := make([]int32, n+1)
	i, j := int32(0), int32(1)
	for i < n-1 {
		for deg[j] > 0 {
			j++
		}
		parents[j] = prufer[i]
		for i < n-2 {
			p := prufer[i]
			deg[p]--
			if p > j || deg[p] > 0 {
				break
			}
			parents[p] = prufer[i+1]
			i++
		}
		i++
		j++
	}
	parents = parents[1:]
	tree := make([][]int32, n)
	for i := int32(1); i < n; i++ {
		p := parents[i-1]
		tree[i-1] = append(tree[i-1], p-1)
		tree[p-1] = append(tree[p-1], i-1)
	}
	return tree
}
