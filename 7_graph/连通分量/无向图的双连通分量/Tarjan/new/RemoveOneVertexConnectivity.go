package main

import (
	"fmt"
)

func main() {
	edges := [][]int{{0, 1}, {1, 2}, {2, 0}, {1, 3}, {3, 4}, {4, 5}, {5, 3}}
	adjList := make([][]Neighbor, 6)
	for i, e := range edges {
		u, v := e[0], e[1]
		adjList[u] = append(adjList[u], Neighbor{v, i})
		adjList[v] = append(adjList[v], Neighbor{u, i})
	}

	C := NewRemoveOneVertexConnectivity(adjList)
	fmt.Println(C.GetCompSize(5, 2))
	fmt.Println(C.IsConnected(3, 4, 2))
}

// 无向图删除一个顶点后的联通性(是否连接、联通分量数、联通分量大小)
type Neighbor struct {
	next, id int
}

type RemoveOneVertexConnectivity struct {
	root          []int32
	lid, rid, low []int32
	children      [][]int32
	removedSize   []int32
	removedComp   []int32
}

func NewRemoveOneVertexConnectivity(graph [][]Neighbor) *RemoveOneVertexConnectivity {
	res := &RemoveOneVertexConnectivity{}
	res._build(graph)
	return res
}

// 删除 removed 后的联通分量数.
func (r *RemoveOneVertexConnectivity) GetCompCount(removed int) int {
	return int(r.removedComp[removed])
}

// 删除 removed 后, v 所在的联通分量大小.
func (r *RemoveOneVertexConnectivity) GetCompSize(removed, v int) int {
	if v == removed {
		return 0
	}
	r32, v32 := int32(removed), int32(v)
	if r.root[v32] != r.root[r32] {
		return int(r._subtreeSize(r.root[v32]))
	}
	if r32 == r.root[v32] {
		return int(r._subtreeSize(r._jump(r32, v32)))
	}
	if !r._inSubtree(v32, r32) {
		return int(r.removedSize[r32])
	}
	v32 = r._jump(r32, v32)
	if r.lid[r32] <= r.low[v32] {
		return int(r._subtreeSize(v32))
	}
	return int(r.removedSize[r32])
}

// 删除 removed 后, u 和 v 是否连通.
func (r *RemoveOneVertexConnectivity) IsConnected(removed, u, v int) bool {
	if u == removed || v == removed {
		return false
	}
	r32 := int32(removed)
	u32, v32 := int32(u), int32(v)
	if r.root[u32] != r.root[v32] {
		return false
	}
	if r.root[u32] != r.root[r32] {
		return true
	}
	inU := r._inSubtree(u32, r32)
	inV := r._inSubtree(v32, r32)
	if inU {
		u32 = r._jump(r32, u32)
		inU = r.low[u32] >= r.lid[r32]
	}
	if inV {
		v32 = r._jump(r32, v32)
		inV = r.low[v32] >= r.lid[r32]
	}
	if inU != inV {
		return false
	}
	if inU {
		return u32 == v32
	}
	return true
}

func (r *RemoveOneVertexConnectivity) _build(graph [][]Neighbor) {
	n := int32(len(graph))
	r.root = make([]int32, n)
	r.lid = make([]int32, n)
	r.rid = make([]int32, n)
	r.low = make([]int32, n)
	r.children = make([][]int32, n)
	for i := int32(0); i < n; i++ {
		r.root[i] = -1
		r.lid[i] = -1
		r.rid[i] = -1
		r.low[i] = -1
	}

	dfsId := int32(0)
	compCount := int32(0)
	var dfs func(cur int32, preId int)
	dfs = func(cur int32, preId int) {
		r.low[cur] = dfsId
		r.lid[cur] = dfsId
		dfsId++
		nexts := graph[cur]
		for i := range nexts {
			e := &nexts[i]
			if e.id == preId {
				continue
			}
			next32 := int32(e.next)
			if r.root[next32] == -1 {
				r.root[next32] = r.root[cur]
				r.children[cur] = append(r.children[cur], next32)
				dfs(next32, e.id)
				r.low[cur] = min32(r.low[cur], r.low[next32])
			} else {
				r.low[cur] = min32(r.low[cur], r.lid[next32])
			}
		}
		r.rid[cur] = dfsId
	}

	for i := int32(0); i < n; i++ {
		if r.root[i] == -1 {
			compCount++
			r.root[i] = i
			dfs(i, -1)
		}
	}
	r.removedSize = make([]int32, n)
	r.removedComp = make([]int32, n)
	for i := range r.removedComp {
		r.removedComp[i] = compCount
	}
	for i := int32(0); i < n; i++ {
		if r.root[i] == i {
			r.removedComp[i] += int32(len(r.children[i])) - 1
		} else {
			r.removedSize[i] = r._subtreeSize(r.root[i]) - 1
			for _, c := range r.children[i] {
				if r.low[c] >= r.lid[i] {
					r.removedSize[i] -= r._subtreeSize(c)
					r.removedComp[i]++
				}
			}
		}
	}
}

func (r *RemoveOneVertexConnectivity) _inSubtree(a, b int32) bool {
	return r.lid[b] <= r.lid[a] && r.lid[a] < r.rid[b]
}

func (r *RemoveOneVertexConnectivity) _subtreeSize(v int32) int32 {
	return r.rid[v] - r.lid[v]
}

func (r *RemoveOneVertexConnectivity) _jump(root, v int32) int32 {
	n := len(r.children[root])
	k := binarySearch(func(x int) bool {
		return r.lid[r.children[root][x]] <= r.lid[v]
	}, 0, n)
	return r.children[root][k]
}

func binarySearch(check func(int) bool, ok, ng int) int {
	for abs(ok-ng) > 1 {
		x := (ng + ok) / 2
		if check(x) {
			ok = x
		} else {
			ng = x
		}
	}
	return ok
}

func min32(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
