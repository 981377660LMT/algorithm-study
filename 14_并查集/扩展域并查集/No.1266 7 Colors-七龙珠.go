package main

import (
	"bufio"
	"fmt"
	"os"
)

// No.1266 7 Colors (七龙珠)
// https://yukicoder.me/problems/no/1266
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const K int32 = 7

	var n, m, q int32
	fmt.Fscan(in, &n, &m, &q)
	states := make([][]bool, n)
	for i := range states {
		var s string
		fmt.Fscan(in, &s)
		states[i] = make([]bool, K)
		for j, c := range s {
			if c == '1' {
				states[i][j] = true
			}
		}
	}
	graph := make([][]int32, n)
	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}

	idx := func(v, c int32) int32 { return K*v + c }

	uf := NewUnionFindArraySimple32(K * n)

	add := func(v, k int32) {
		states[v][k] = true
		for _, to := range graph[v] {
			if states[to][k] {
				uf.Union(idx(v, k), idx(to, k), nil)
			}
		}
		a := k - 1
		b := k + 1
		if a < 0 {
			a = K - 1
		}
		if b == K {
			b = 0
		}
		if states[v][a] {
			uf.Union(idx(v, a), idx(v, k), nil)
		}
		if states[v][b] {
			uf.Union(idx(v, b), idx(v, k), nil)
		}
	}

	for v := int32(0); v < n; v++ {
		for k := int32(0); k < K; k++ {
			if states[v][k] {
				add(v, k)
			}
		}
	}

	for i := int32(0); i < q; i++ {
		var t, x, y int32
		fmt.Fscan(in, &t, &x, &y)
		if t == 1 {
			x, y = x-1, y-1
			add(x, y)
		} else if t == 2 {
			x--
			res := uf.Size(idx(x, 0))
			fmt.Fprintln(out, res)
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

func (u *UnionFindArraySimple32) Union(key1, key2 int32, beforeMerge func(big, small int32)) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	if beforeMerge != nil {
		beforeMerge(root1, root2)
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = root1
	u.Part--
	return true
}

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
