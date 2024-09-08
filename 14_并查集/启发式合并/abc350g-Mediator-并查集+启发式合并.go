package main

import (
	"bufio"
	"fmt"
	"os"
)

// [ABC350G] Mediator (并查集+启发式合并)
// https://www.luogu.com.cn/problem/AT_abc350_g
// 初始时，有 N 个点，编号为 0 到 N-1,没有边存在.
// 有 Q 次操作，每次操作有三个整数 a, b, c.
// 1 u v: 在 u 和 v 之间添加一条边,保证 u 和 v 不在同一个连通块中.
// 2 u v: 询问是否存在和 u,v 都相邻的点，若存在输出编号，若不存在输出0.
//
// !每次并查集合并后，都可以对整个小的分组进行某种更新, 实时维护每个点的父节点.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	const MOD int = 998244353

	var n, q int
	fmt.Fscan(in, &n, &q)

	adjList := make([][]int32, n)
	parent := make([]int32, n)
	for i := range parent {
		parent[i] = -1
	}
	uf := NewUnionFindArraySimple32(int32(n))

	var updateParent func(cur, pre int32)
	updateParent = func(cur, pre int32) {
		parent[cur] = pre
		for _, next := range adjList[cur] {
			if next != pre {
				updateParent(next, cur)
			}
		}
	}

	// 添加边.
	addEdge := func(a, b int32) {
		adjList[a] = append(adjList[a], b)
		adjList[b] = append(adjList[b], a)
		if uf.Size(a) > uf.Size(b) {
			a, b = b, a
		}
		uf.Union(a, b)
		updateParent(a, b)
	}

	// 查询a和b是否有共同的邻居.
	query := func(a, b int32) (int32, bool) {
		if parent[a] == parent[b] && parent[a] != -1 {
			return parent[a], true
		}
		if parent[a] != -1 && parent[parent[a]] == b {
			return parent[a], true
		}
		if parent[b] != -1 && parent[parent[b]] == a {
			return parent[b], true
		}
		return 0, false
	}

	preRes := 0
	for i := 0; i < q; i++ {
		var A, B, C int
		fmt.Fscan(in, &A, &B, &C)
		A = 1 + (((A * (1 + preRes)) % MOD) % 2)
		B = 1 + (((B * (1 + preRes)) % MOD) % n)
		C = 1 + (((C * (1 + preRes)) % MOD) % n)
		B--
		C--
		if A == 1 {
			addEdge(int32(B), int32(C))
		} else {
			res, ok := query(int32(B), int32(C))
			if !ok {
				preRes = 0
				fmt.Fprintln(out, preRes)
			} else {
				preRes = int(res + 1)
				fmt.Fprintln(out, preRes)
			}
		}
	}
}

type UnionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple32(n int32) *UnionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple32) Union(key1, key2 int32) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = int32(root1)
	u.Part--
	return true
}

// 非递归版本.
func (u *UnionFindArraySimple32) Find(key int32) int32 {
	root := key
	for u.data[root] >= 0 {
		root = u.data[root]
	}
	for key != root {
		key, u.data[key] = u.data[key], root
	}
	return root
}

func (u *UnionFindArraySimple32) Size(key int32) int32 {
	return -u.data[u.Find(key)]
}
