package main

import "fmt"

func main() {
	arr := []int32{1, 2, 3}
	EnumerateProduct(arr, func(view []int32) {
		fmt.Println(view)
	})

	EnumerateLabeledTree(3, func(tree [][]int32) {
		fmt.Println(tree)
	})
}

// n阶`有标号无根树`枚举，个数为n^{n-2}
// 凯莱定理：n个点的完全图的生成树有n^(n-2)个。
func EnumerateLabeledTree(n int32, f func(tree [][]int32)) {
	if n == 1 {
		tree := make([][]int32, 1)
		f(tree)
		return
	}

	buffer := make([]int32, n-1)
	buffer[n-2] = n - 1
	ends := make([]int32, n-2)
	for i := int32(0); i < n-2; i++ {
		ends[i] = n
	}
	EnumerateProduct(ends, func(code []int32) {
		copy(buffer, code)
		tree := FromPruferCode(buffer)
		f(tree)
	})
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

// [0,arr[0]) x [0,arr[1]) x ...
func EnumerateProduct(arr []int32, f func(view []int32)) {
	n := len(arr)
	var dfs func([]int32)
	dfs = func(view []int32) {
		if len(view) == n {
			f(view)
			return
		}
		for x := int32(0); x < arr[len(view)]; x++ {
			view = append(view, x)
			dfs(view)
			view = view[:len(view)-1]
		}
	}
	p := make([]int32, 0, len(arr))
	dfs(p)
}
