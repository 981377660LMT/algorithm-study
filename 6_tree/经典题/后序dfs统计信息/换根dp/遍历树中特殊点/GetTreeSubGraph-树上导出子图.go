package main

import (
	"bufio"
	"fmt"
	"os"
)

// D - Minimum Steiner Tree
// https://atcoder.jp/contests/abc368/tasks/abc368_d
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int32
	fmt.Fscan(in, &n, &k)
	edges := make([][2]int32, n-1)
	tree := make([][]int32, n)
	for i := int32(0); i < n-1; i++ {
		var a, b int32
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		edges[i] = [2]int32{a, b}
		tree[a] = append(tree[a], b)
		tree[b] = append(tree[b], a)
	}
	specials := make([]int32, k)
	for i := int32(0); i < k; i++ {
		fmt.Fscan(in, &specials[i])
		specials[i]--
	}

	_, _, ok := GetTreeSubGraph(n, tree, specials)
	res := 0
	for _, v := range ok {
		if v {
			res++
		}
	}
	fmt.Fprintln(out, res)
}

// 树上导出子图/树上诱导子图.
func GetTreeSubGraph(n int32, rawTree [][]int32, specials []int32) (newTree [][]int32, newEdges [][2]int32, inNewTree []bool) {
	if len(specials) == 0 {
		return make([][]int32, n), [][2]int32{}, make([]bool, n)
	}
	inNewTree = make([]bool, n)
	for _, v := range specials {
		inNewTree[v] = true
	}

	var dfs func(int32, int32) bool
	dfs = func(cur, pre int32) bool {
		for _, next := range rawTree[cur] {
			if next != pre && dfs(next, cur) {
				inNewTree[cur] = true
			}
		}
		return inNewTree[cur]
	}

	root := specials[0]
	dfs(root, -1)

	newEdges, newTree = [][2]int32{}, make([][]int32, n)
	for cur := int32(0); cur < n; cur++ {
		for _, next := range rawTree[cur] {
			if cur < next && inNewTree[cur] && inNewTree[next] {
				newEdges = append(newEdges, [2]int32{cur, next})
				newTree[cur] = append(newTree[cur], next)
				newTree[next] = append(newTree[next], cur)
			}
		}
	}

	return
}
