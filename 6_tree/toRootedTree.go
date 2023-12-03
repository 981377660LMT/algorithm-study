// https://nyaannyaan.github.io/library/tree/convert-tree.hpp

package main

import "fmt"

func main() {
	tree := [][]int{{1, 2}, {0, 3, 4}, {0, 5, 6}, {1}, {1}, {2}, {2}}
	rootedTree := ToRootedTree(tree, 0)
	fmt.Println(rootedTree)
}

// 无根树转有根树.
func ToRootedTree(tree [][]int, root int) [][]int {
	n := len(tree)
	rootedTree := make([][]int, n)
	visited := make([]bool, n)
	visited[root] = true
	queue := []int{root}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, next_ := range tree[cur] {
			if !visited[next_] {
				visited[next_] = true
				queue = append(queue, next_)
				rootedTree[cur] = append(rootedTree[cur], next_)
			}
		}
	}
	return rootedTree
}
