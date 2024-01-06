// 无向图删除一条边后的联通性(是否连接、联通分量数、联通分量大小)

package main

import "fmt"

func main() {
	es := [][]int{{0, 1}, {1, 2}, {2, 0}, {1, 3}, {3, 4}, {4, 5}, {5, 3}}

	adjList := make([][]Neighbor, 6)
	for i, e := range es {
		u, v := e[0], e[1]
		adjList[u] = append(adjList[u], Neighbor{v, i})
		adjList[v] = append(adjList[v], Neighbor{u, i})
	}
	edgeWithId := make([]Edge, len(es))
	for i, e := range es {
		edgeWithId[i] = Edge{e[0], e[1], i}
	}
	C := NewRemoveOneEdgeConnectivity(adjList, edgeWithId)
	fmt.Println(C.GetCompSize(5, 2))
	fmt.Println(C.GetCompSize(3, 2))
	fmt.Println(C.IsConnected(3, 1, 2))
}

type Neighbor struct{ next, eid int }
type Edge struct{ u, v, eid int }

// 无向图删除一条边后的联通性(是否连接、联通分量数、联通分量大小)
type RemoveOneEdgeConnectivity struct {
	compBase int32
	idToEdge [][2]int32
	lid, rid []int32
	root     []int32
	isBridge []bool
}

func NewRemoveOneEdgeConnectivity(graph [][]Neighbor, edges []Edge) *RemoveOneEdgeConnectivity {
	res := &RemoveOneEdgeConnectivity{}
	res._build(graph, edges)
	return res
}

// 删除 removedEdgeId 后的联通分量数.
func (r *RemoveOneEdgeConnectivity) GetCompCount(removedEdgeId int) int {
	if r.isBridge[removedEdgeId] {
		return int(r.compBase) + 1
	}
	return int(r.compBase)
}

// 删除 removedEdgeId 后, v 所在的联通分量大小.
func (r *RemoveOneEdgeConnectivity) GetCompSize(removedEdgeId, v int) int {
	rt := r.root[v]
	if !r.isBridge[removedEdgeId] {
		return int(r._subtreeSize(rt))
	}
	a, b := r.idToEdge[removedEdgeId][0], r.idToEdge[removedEdgeId][1]
	if rt != r.root[a] {
		return int(r._subtreeSize(rt))
	}
	if r._inSubtree(int32(v), b) {
		return int(r._subtreeSize(b))
	}
	return int(r._subtreeSize(rt) - r._subtreeSize(b))
}

// 删除 removedEdgeId 后, u 和 v 是否连通.
func (r *RemoveOneEdgeConnectivity) IsConnected(removedEdgeId, u, v int) bool {
	if r.root[u] != r.root[v] {
		return false
	}
	if !r.isBridge[removedEdgeId] {
		return true
	}
	b := r.idToEdge[removedEdgeId][1]
	return r._inSubtree(int32(u), b) == r._inSubtree(int32(v), b)
}

func (r *RemoveOneEdgeConnectivity) _build(graph [][]Neighbor, edges []Edge) {
	n, m := int32(len(graph)), int32(len(edges))
	r.idToEdge = make([][2]int32, m)
	for i := range edges {
		e := &edges[i]
		r.idToEdge[e.eid] = [2]int32{int32(e.u), int32(e.v)}
	}
	r.root = make([]int32, n)
	r.lid = make([]int32, n)
	r.rid = make([]int32, n)
	r.isBridge = make([]bool, m)
	low := make([]int32, n)
	for i := int32(0); i < n; i++ {
		r.root[i] = -1
		r.lid[i] = -1
		r.rid[i] = -1
		low[i] = -1
	}

	dfn := int32(0)
	var dfs func(int32, int)
	dfs = func(v int32, lastEid int) {
		low[v] = dfn
		r.lid[v] = dfn
		dfn++
		nexts := graph[v]
		for i := range nexts {
			e := &nexts[i]
			if e.eid == lastEid {
				continue
			}
			next32 := int32(e.next)
			if r.root[next32] == -1 {
				r.root[next32] = r.root[v]
				dfs(next32, e.eid)
				low[v] = min32(low[v], low[next32])
				r.isBridge[e.eid] = low[next32] == r.lid[next32]
			} else {
				low[v] = min32(low[v], r.lid[next32])
			}
		}
		r.rid[v] = dfn
	}

	for i := int32(0); i < n; i++ {
		if r.root[i] == -1 {
			r.compBase++
			r.root[i] = i
			dfs(i, -1)
		}
	}
	for i := range r.idToEdge {
		a, b := r.idToEdge[i][0], r.idToEdge[i][1]
		if r._inSubtree(a, b) {
			r.idToEdge[i] = [2]int32{b, a}
		}
	}
}

func (r *RemoveOneEdgeConnectivity) _inSubtree(a, b int32) bool {
	return r.lid[b] <= r.lid[a] && r.lid[a] < r.rid[b]
}

func (r *RemoveOneEdgeConnectivity) _subtreeSize(v int32) int32 {
	return r.rid[v] - r.lid[v]
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
