// Floyd-Warshall 算法求有向图的闭包传递问题。通俗地讲，就是可达性问题。
// O(n^3/w)

// NewTransitiveClosure(n) 构造一个n*n的传递闭包.
// AddDirectedEdge(from, to) 添加一条有向边.
// Build() 构造传递闭包.
// CanReach(from, to) 判断是否可达.

// 3000*3000 => 250ms
// 4000*4000 => 620ms
// 5000*5000 => 1.3s.

package main

import (
	"fmt"
	"time"
)

func main() {
	n := 5000
	time1 := time.Now()
	T := NewTransitiveClosure(n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			T.AddDirectedEdge(i, j)
		}
	}
	T.Build()
	time2 := time.Now()
	fmt.Println(fmt.Sprintf("%v*%v的传递闭包:%v", n, n, time2.Sub(time1))) // 5000*5000的传递闭包:1.3331916s
}

// https://leetcode.cn/problems/course-schedule-iv/
func checkIfPrerequisite(numCourses int, prerequisites [][]int, queries [][]int) []bool {
	trans := NewTransitiveClosure(numCourses)
	for _, p := range prerequisites {
		trans.AddDirectedEdge(p[0], p[1])
	}
	res := make([]bool, len(queries))
	for i, q := range queries {
		res[i] = trans.CanReach(q[0], q[1])
	}
	return res
}

// 有向图的传递闭包.
type TransitiveClosure struct {
	n        int
	canReach []_BitSet64
	hasBuilt bool
}

func NewTransitiveClosure(n int) *TransitiveClosure {
	canReach := make([]_BitSet64, n)
	for i := range canReach {
		canReach[i] = NewBitset(n)
	}
	return &TransitiveClosure{n: n, canReach: canReach}
}

func (tc *TransitiveClosure) AddDirectedEdge(from, to int) {
	if tc.hasBuilt {
		panic("can't add edge after build")
	}
	tc.canReach[from].Set(to)
}

func (tc *TransitiveClosure) Build() {
	if tc.hasBuilt {
		panic("can't build twice")
	}
	n, canReach := tc.n, tc.canReach
	for k := 0; k < n; k++ {
		cacheK := canReach[k]
		for i := 0; i < n; i++ {
			cacheI := canReach[i]
			if cacheI.Has(k) {
				cacheI.IOr(cacheK)
			}
		}
	}
	tc.hasBuilt = true
}

func (tc *TransitiveClosure) CanReach(from, to int) bool {
	if !tc.hasBuilt {
		tc.Build()
	}
	return tc.canReach[from].Has(to)
}

type _BitSet64 []uint64

func NewBitset(n int) _BitSet64 { return make(_BitSet64, n>>6+1) }

func (b _BitSet64) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 }
func (b _BitSet64) Set(p int)      { b[p>>6] |= 1 << (p & 63) }

func (b _BitSet64) IOr(c _BitSet64) _BitSet64 {
	for i, v := range c {
		b[i] |= v
	}
	return b
}
