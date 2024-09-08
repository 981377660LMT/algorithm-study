// No.416 旅行会社 (启发式合并+反向并查集删边)
// https://yukicoder.me/problems/no/416
// 给定一个无向图.
// q次断开两条边.
// 对每个点，求第几次删边操作后，无法从该点到达0号点.
// 如果一直可以到达0号点, 输出-1.
// 如果一开始就无法到达0号点, 输出0.

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int32
	fmt.Fscan(in, &n, &m, &q)
	edges := make([][2]int32, 0, m)
	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		edges = append(edges, [2]int32{u - 1, v - 1})
	}
	queries := make([][2]int32, 0, q)
	for i := int32(0); i < q; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		queries = append(queries, [2]int32{u - 1, v - 1})
	}

	res := make([]int32, n)
	uf := NewUnionFindArraySimple32(n)
	groups := make([][]int32, n)
	for i := int32(0); i < n; i++ {
		groups[i] = []int32{i}
	}
	unionAt := func(time int32, a, b int32) {
		ra, rb := uf.Find(a), uf.Find(b)
		if ra == rb {
			return
		}
		r0 := uf.Find(0)

		// 连接后，原来与0不连通的点，现在连通了，说明删除的这条边是关键的.
		if r0 == ra || r0 == rb {
			other := r0 ^ ra ^ rb
			for _, v := range groups[other] {
				res[v] = time
			}
		}

		uf.Union(ra, rb, func(big, small int32) {
			for _, v := range groups[small] {
				groups[big] = append(groups[big], v)
			}
			groups[small] = nil
		})
	}

	// set(edges) - set(queries)
	keptEdges := func() map[[2]int32]struct{} {
		res, removed := make(map[[2]int32]struct{}), make(map[[2]int32]struct{}, q)
		for _, e := range queries {
			a, b := e[0], e[1]
			if a > b {
				a, b = b, a
			}
			removed[[2]int32{a, b}] = struct{}{}
		}
		for _, e := range edges {
			a, b := e[0], e[1]
			if a > b {
				a, b = b, a
			}
			if _, has := removed[[2]int32{a, b}]; !has {
				res[[2]int32{a, b}] = struct{}{}
			}
		}
		return res
	}()
	for e := range keptEdges {
		unionAt(-1, e[0], e[1])
	}
	for qi := q - 1; qi >= 0; qi-- {
		unionAt(qi+1, queries[qi][0], queries[qi][1])
	}

	for _, v := range res[1:] {
		fmt.Fprintln(out, v)
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

func (u *UnionFindArraySimple32) GetSize(key int32) int32 {
	return -u.data[u.Find(key)]
}
