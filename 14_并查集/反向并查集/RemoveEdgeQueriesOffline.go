// 反向并查集，删边变成加边，然后按照加边的顺序进行操作。
// RemoveEdgeQueriesOffline

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	abc120d()
	// yuki416()
}

// D - Decayed Bridges
// https://atcoder.jp/contests/abc120/tasks/abc120_d
// 桥按顺序断掉.求每座桥断掉之后,不能互相到达的岛屿对数(a<b).
func abc120d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int32
	fmt.Fscan(in, &n, &m)
	edges := make([][2]int32, 0, m)
	for i := int32(0); i < m; i++ {
		var a, b int32
		fmt.Fscan(in, &a, &b)
		edges = append(edges, [2]int32{a - 1, b - 1})
	}
	S := NewRemoveEdgeQueriesOffline(n)
	for i := int32(0); i < m; i++ {
		S.InitEdge(edges[i][0], edges[i][1])
		S.RemoveEdge(edges[i][0], edges[i][1])
	}

	diff := make([]int, m)
	S.Run(
		func(qi int32, uf *UnionFindArraySimple32, a, b int32) {
			ra, rb := uf.Find(a), uf.Find(b)
			if ra == rb {
				diff[qi] = 0
			} else {
				diff[qi] = int(uf.Size(ra)) * int(uf.Size(rb))
			}
			uf.Union(ra, rb, nil)
		},
	)

	res := make([]int, m)
	res[m-1] = int(n) * (int(n) - 1) / 2
	for i := int32(m - 2); i >= 0; i-- {
		res[i] = res[i+1] - diff[i+1]
	}
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// No.416 旅行会社 (启发式合并+反向并查集删边)
// https://yukicoder.me/problems/no/416
// 对每个点，求第几次删边操作后，无法从该点到达0号点.
// 如果一直可以到达0号点, 输出-1.
// 如果一开始就无法到达0号点, 输出0.
func yuki416() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int32
	fmt.Fscan(in, &n, &m, &q)
	S := NewRemoveEdgeQueriesOffline(n)
	for i := int32(0); i < m; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		S.InitEdge(u-1, v-1)
	}
	for i := int32(0); i < q; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		S.RemoveEdge(u-1, v-1)
	}

	res := make([]int32, n)
	for i := int32(0); i < n; i++ {
		res[i] = -2
	}
	groups := make([][]int32, n)
	for i := int32(0); i < n; i++ {
		groups[i] = []int32{i}
	}
	S.Run(
		func(qi int32, uf *UnionFindArraySimple32, a, b int32) {
			ra, rb := uf.Find(a), uf.Find(b)
			if ra == rb {
				return
			}
			r0 := uf.Find(0)

			// 连接后，原来与0不连通的点，现在连通了，说明删除的这条边是关键的.
			if r0 == ra || r0 == rb {
				other := r0 ^ ra ^ rb
				for _, v := range groups[other] {
					res[v] = qi // !qi: -1(始终保留的边) or 0~q-1(第0~q-1次删边操作)
				}
			}

			uf.Union(ra, rb, func(big, small int32) {
				for _, v := range groups[small] {
					groups[big] = append(groups[big], v)
				}
				groups[small] = nil
			})
		},
	)

	for _, v := range res[1:] {
		if v == -2 { // 未连接到0号点
			fmt.Fprintln(out, 0)
		} else if v == -1 { // 一直可以到达0号点
			fmt.Fprintln(out, -1)
		} else { // 第几次删边操作后，无法从该点到达0号点
			fmt.Fprintln(out, v+1)
		}
	}
}

type RemoveEdgeQueriesOffline struct {
	n            int32
	initialEdges [][2]int32
	removedEdges [][2]int32
	uf           *UnionFindArraySimple32
}

func NewRemoveEdgeQueriesOffline(n int32) *RemoveEdgeQueriesOffline {
	return &RemoveEdgeQueriesOffline{
		n:  n,
		uf: NewUnionFindArraySimple32(n),
	}
}

func (r *RemoveEdgeQueriesOffline) InitEdge(a, b int32) {
	if a > b {
		a, b = b, a
	}
	r.initialEdges = append(r.initialEdges, [2]int32{a, b})
}

func (r *RemoveEdgeQueriesOffline) RemoveEdge(a, b int32) {
	if a > b {
		a, b = b, a
	}
	r.removedEdges = append(r.removedEdges, [2]int32{a, b})
}

// 第qi次操作连接a和b.
//
// !qi: -1(始终保留的边) or 0~q-1(第0~q-1次删边操作).
func (r *RemoveEdgeQueriesOffline) Run(f func(qi int32, uf *UnionFindArraySimple32, a, b int32)) {
	removedSet := make(map[[2]int32]struct{}, len(r.removedEdges))
	for _, e := range r.removedEdges {
		removedSet[e] = struct{}{}
	}
	for _, e := range r.initialEdges {
		if _, has := removedSet[e]; !has {
			// keptEdge
			f(-1, r.uf, e[0], e[1])
		}
	}
	for qi := int32(len(r.removedEdges)) - 1; qi >= 0; qi-- {
		f(qi, r.uf, r.removedEdges[qi][0], r.removedEdges[qi][1])
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
