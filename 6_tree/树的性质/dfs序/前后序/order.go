//       0
//     / | \
//   1   2  3
//  / \     |
// 4   5    6
//
// !preOrder:
//  0 => [0, 7)
//  1 => [1, 4)
//  4 => [2, 3)
//  5 => [3, 4)
//  2 => [4, 5)
//  3 => [5, 7)
//  6 => [6, 7)
//
// !postOrder:
//  4 => [0, 1)
//  5 => [1, 2)
//  1 => [0, 3)
//  2 => [3, 4)
//  6 => [4, 5)
//  3 => [4, 6)
//  0 => [0, 7)

package main

import "fmt"

func main() {
	fmt.Println("Order")
	n := int32(7)
	tree := make([][]int32, n)
	edegs := [][]int32{{0, 1}, {0, 2}, {0, 3}, {1, 4}, {1, 5}, {2, 6}}
	for _, e := range edegs {
		tree[e[0]] = append(tree[e[0]], e[1])
		tree[e[1]] = append(tree[e[1]], e[0])
	}

	lid, rid := DfsPreOrder(tree, 0)
	fmt.Println(lid, rid)

	lid, rid = DfsPostOrder(tree, 0)
	fmt.Println(lid, rid)
}

// 前序遍历dfs序.
// !data[lid[i]] = values[i]
func DfsPreOrder(tree [][]int32, root int32) (lid, rid []int32) {
	n := int32(len(tree))
	lid, rid = make([]int32, n), make([]int32, n)
	dfn := int32(0)

	var dfs func(cur, pre int32)
	dfs = func(cur, pre int32) {
		lid[cur] = dfn
		dfn++
		for _, next := range tree[cur] {
			if next != pre {
				dfs(next, cur)
			}
		}
		rid[cur] = dfn
	}
	dfs(root, -1)
	return
}

// 后序遍历dfs序.
// !data[rid[i]-1] = values[i]
func DfsPostOrder(tree [][]int32, root int32) (lid, rid []int32) {
	n := int32(len(tree))
	lid, rid = make([]int32, n), make([]int32, n)
	dfn := int32(0)

	var dfs func(cur, pre int32)
	dfs = func(cur, pre int32) {
		lid[cur] = dfn
		for _, next := range tree[cur] {
			if next != pre {
				dfs(next, cur)
			}
		}
		dfn++
		rid[cur] = dfn
	}
	dfs(root, -1)
	return
}
