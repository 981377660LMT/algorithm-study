package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var SEED uint = uint(time.Now().UnixNano()/2 + 1)

// https://www.luogu.com.cn/problem/P5043
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var q int
	fmt.Fscan(in, &q)

	// 接下来 q 行，每行包含若干个整数，表示一个树。第一个整数 n 表示点数。
	// 接下来 n  个整数，依次表示编号为 1 到 n  的每个点的父亲结点的编号。根节点父亲结点编号为 0。
	// 输出 q 行，每行一个整数，表示与每个树同构的树的最小编号。
	res := make([]int, q)
	visited := make(map[uint]int)

	for i := range res {
		var n int
		fmt.Fscan(in, &n)
		parents := make([]int, n)
		for j := range parents {
			fmt.Fscan(in, &parents[j])
		}

		tree := make([][]int, n)
		root := -1
		for cur, pre := range parents {
			pre--
			if pre == -1 {
				root = cur
				continue
			}
			tree[pre] = append(tree[pre], cur)
			tree[cur] = append(tree[cur], pre)
		}

		centroids := findCentroids(n, tree, root)
		hashes := make([]uint, len(centroids))
		for i, centroid := range centroids {
			hashes[i] = treeHash(n, tree, centroid, SEED)
		}

		hash := min(hashes...)
		if id, ok := visited[hash]; ok {
			res[i] = id
		} else {
			res[i] = i + 1
			visited[hash] = i + 1
		}
	}

	for _, v := range res {
		fmt.Fprintln(out, v)
	}

}

type Edge struct{ u, v int }

func treeHash(n int, tree [][]int, root int, seed uint) uint {
	bases := make([]uint, n)
	for i := range bases {
		bases[i] = fastRand(&seed)
	}
	depths := make([]int, n)
	hashes := make([]uint, n)
	for i := range hashes {
		hashes[i] = 1
	}

	var dfs func(cur, pre int) int
	dfs = func(cur, pre int) int {
		dep := 0
		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			dep = max(dep, dfs(next, cur)+1)
		}
		base := bases[dep]
		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			hashes[cur] *= base + hashes[next]
		}
		depths[cur] = dep
		return dep
	}

	dfs(root, -1)
	return hashes[root]
}

func fastRand(seed *uint) uint {
	*seed ^= *seed << 13
	*seed ^= *seed >> 17
	*seed ^= *seed << 5
	return *seed
}

func findCentroids(n int, tree [][]int, root int) (centroids []int) {
	weight := make([]int, n)
	subSize := make([]int, n)
	var dfs func(cur, pre int)
	dfs = func(cur, pre int) {
		subSize[cur] = 1
		for _, next := range tree[cur] {
			if next == pre {
				continue
			}
			dfs(next, cur)
			subSize[cur] += subSize[next]
			weight[cur] = max(weight[cur], subSize[next])
		}
		weight[cur] = max(weight[cur], n-subSize[cur])
		if weight[cur] <= n/2 {
			centroids = append(centroids, cur)
		}
	}

	dfs(root, -1)
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(nums ...uint) uint {
	res := nums[0]
	for _, v := range nums[1:] {
		if v < res {
			res = v
		}
	}
	return res
}
