//
// Italiano's dynamic reachability data structure for DAG
//
// Description:
//   It is a data structure that admits the following operations:
//     add_edge(s, t):     insert edge (s,t) to the network if
//                         it does not make a cycle
//
//     is_reachable(s, t): return true iff there is a path s --> t
//
// Algorithm:
//   We maintain reachability trees T(u) for all u in V.
//   Then is_reachable(s, t) is solved by checking "t in T(u)".
//   For add_edge(s, t), if is_reachable(s, t) or is_reachable(t, s) then
//   no update is performed. Otherwise, we meld T(s) and T(t).
//
// Complexity:
//   !amortized O(n) per update
//
// Verified:
//   SPOJ 9458: Ghosts having fun
//
// References:
//   Giuseppe F. Italiano (1988):
//   Finding paths and deleting edges in directed acyclic graphs.
//   Information Processing Letters, vol. 28, no. 1, pp. 5--11.
//
// 维护DAG可到达性.

package main

import "fmt"

func main() {
	D := NewDAGReachability(5)
	D.AddEdge(0, 1)
	fmt.Println(D.IsReachable(0, 1))
	fmt.Println(D.IsReachable(1, 0))
	D.AddEdge(1, 2)
	fmt.Println(D.IsReachable(0, 2))
	fmt.Println(D.IsReachable(2, 0))
	D.AddEdge(2, 3)
	fmt.Println(D.IsReachable(0, 3))
	fmt.Println(D.IsReachable(3, 0))
	D.AddEdge(3, 4)
	fmt.Println(D.IsReachable(0, 4))
	fmt.Println(D.IsReachable(4, 0))
}

// 动态DAG可到达性.
type DAGReachability struct {
	n      int
	parent [][]int
	child  [][][]int
}

func NewDAGReachability(n int) *DAGReachability {
	res := &DAGReachability{}

	parent := make([][]int, n)
	for i := 0; i < n; i++ {
		parent[i] = make([]int, n)
		for j := 0; j < n; j++ {
			parent[i][j] = -1
		}
	}

	child := make([][][]int, n)
	for i := 0; i < n; i++ {
		child[i] = make([][]int, n)
	}

	res.n = n
	res.parent = parent
	res.child = child
	return res
}

// 判断是否可到达.
func (dr *DAGReachability) IsReachable(from, to int) bool {
	return from == to || dr.parent[from][to] >= 0
}

// 添加边.如果添加后形成环, 则返回false.
func (dr *DAGReachability) AddEdge(from, to int) bool {
	if dr.IsReachable(to, from) {
		return false
	}
	if dr.IsReachable(from, to) {
		return true
	}
	for p := 0; p < dr.n; p++ {
		if dr.IsReachable(p, from) && !dr.IsReachable(p, to) {
			dr.meld(p, to, from, to)
		}
	}
	return true
}

func (dr *DAGReachability) meld(root, sub, u, v int) {
	dr.parent[root][v] = u
	dr.child[root][u] = append(dr.child[root][u], v)
	for _, c := range dr.child[sub][v] {
		if !dr.IsReachable(root, c) {
			dr.meld(root, sub, v, c)
		}
	}
}
