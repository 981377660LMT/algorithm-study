// 反向并查集，删边变成加边，然后按照加边的顺序进行操作。
// RemoveEdgeQueriesOffline

package main

func main() {

}

type RemoveEdgeQueriesOffline struct {
	n            int32
	initialEdges [][2]int32
	removedEdges [][2]int32
	uf           *UnionFindArraySimple32
}

func NewRemoveEdgeQueriesOffline(n, m, q int32) *RemoveEdgeQueriesOffline {
	return &RemoveEdgeQueriesOffline{
		n:            n,
		initialEdges: make([][2]int32, 0, m),
		removedEdges: make([][2]int32, 0, m),
		uf:           NewUnionFindArraySimple32(n),
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

// 第kth次操作连接a和b.
// k=-1表示始终保留这条边.
// 第一次删边操作的k从0开始.
func (r *RemoveEdgeQueriesOffline) Run(f func(kth int32, uf *UnionFindArraySimple32, a, b int32)) {
	removed := make(map[[2]int32]struct{})
	for _, e := range r.removedEdges {
		removed[e] = struct{}{}
	}
	keptEdges := make(map[[2]int32]struct{})
	for _, e := range r.initialEdges {
		if _, has := removed[e]; !has {
			keptEdges[e] = struct{}{}
		}
	}
	for e := range keptEdges {
		f(-1, r.uf, e[0], e[1])
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

func (u *UnionFindArraySimple32) GetSize(key int32) int32 {
	return -u.data[u.Find(key)]
}
