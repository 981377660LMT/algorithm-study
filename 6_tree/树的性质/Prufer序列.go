package main

import "fmt"

func main() {
	//    0
	//   / \
	//  1   2
	//     / \
	//    3   4

	edges := [][]int32{{0, 1}, {0, 2}, {2, 3}, {2, 4}}
	tree := make([][]int32, 5)
	for _, e := range edges {
		tree[e[0]] = append(tree[e[0]], e[1])
		tree[e[1]] = append(tree[e[1]], e[0])
	}
	fmt.Println(ToPruferCode(tree))                  // [0 2 2 4]
	fmt.Println(FromPruferCode([]int32{0, 2, 2, 4})) // [[1 2] [0] [0 3 4] [2] [2]]
}

func ToPruferCode(tree [][]int32) []int32 {
	n := int32(len(tree))
	parent := make([]int32, n)
	{
		var dfs func(v, p int32)
		dfs = func(v, p int32) {
			parent[v] = p
			for _, to := range tree[v] {
				if to != p {
					dfs(to, v)
				}
			}
		}
		dfs(n-1, -1)
	}

	degree := make([]int32, n)
	for cur := int32(0); cur < n; cur++ {
		degree[cur] = int32(len(tree[cur]))
	}
	res := make([]int32, n-1)
	p := int32(0)
	leaf := int32(-1)
	for i := int32(0); i < n-1; i++ {
		if leaf == -1 {
			for degree[p] != 1 {
				p++
			}
			leaf = p
		}
		res[i] = parent[leaf]
		degree[leaf]--
		degree[parent[leaf]]--
		if degree[parent[leaf]] == 1 && parent[leaf] < p {
			leaf = parent[leaf]
		} else {
			leaf = -1
		}
	}
	return res
}

func FromPruferCode(code []int32) [][]int32 {
	n := int32(len(code)) + 1
	tree := make([][]int32, n)
	if n == 1 {
		return tree
	}
	degree := make([]int32, n)
	for i := int32(0); i < n; i++ {
		degree[i] = 1
	}
	for _, v := range code {
		degree[v]++
	}

	p := int32(0)
	leaf := int32(-1)
	for i := int32(0); i < n-1; i++ {
		if leaf == -1 {
			for degree[p] != 1 {
				p++
			}
			leaf = p
		}
		tree[code[i]] = append(tree[code[i]], leaf)
		tree[leaf] = append(tree[leaf], code[i])
		degree[leaf]--
		degree[code[i]]--
		if code[i] < p && degree[code[i]] == 1 {
			leaf = code[i]
		} else {
			leaf = -1
		}
	}
	return tree
}
