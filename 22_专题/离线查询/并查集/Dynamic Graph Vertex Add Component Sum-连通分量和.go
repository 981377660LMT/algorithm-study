// 离线动态连通性查询
// https://suisen-cp.github.io/cp-library-cpp/library/algorithm/offline_dynamic_connectivity_component_sum.hpp

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// 0 u v 连接u v (保证u v不连接)
	// 1 u v 断开u v  (保证u v连接)
	// 2 u x 将u的值加上x
	// 3 u 输出u所在连通块的值

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	sums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &sums[i])
	}
	queries := make([][]int, 0, q)
	for i := 0; i < q; i++ {
		var op, u, v, add int
		fmt.Fscan(in, &op)
		if op == 0 || op == 1 {
			fmt.Fscan(in, &u, &v)
			queries = append(queries, []int{op, u, v})
		} else if op == 2 {
			fmt.Fscan(in, &u, &add)
			queries = append(queries, []int{op, u, add})
		} else {
			fmt.Fscan(in, &u)
			queries = append(queries, []int{op, u})
		}
	}

	odc := NewOfflineDynamicConnectivity(n, q)
	for i, v := range sums {
		odc.Uf.Add(i, v) // 开始时每个点的权值
	}
	for i, query := range queries { // 处理边的添加和删除
		op := query[0]
		if op == 0 {
			u, v := query[1], query[2]
			odc.AddEdge(u, v, i)
		} else if op == 1 {
			u, v := query[1], query[2]
			odc.RemoveEdge(u, v, i)
		}
	}
	odc.Build()

	// 处理权值查询
	res := []int{}
	f := func(k int) {
		op := queries[k][0]
		if op == 2 {
			index, add := queries[k][1], queries[k][2]
			odc.Uf.Add(index, add)
		} else if op == 3 {
			index := queries[k][1]
			res = append(res, odc.Uf.Get(index))
		}
	}
	odc.Run(f)

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type S = int

type OfflineDynamicConnectivity struct {
	Uf         *UndoDSU
	Part       int // 连通分量数
	n, q, size int
	seg        [][]int
	edges      [][]int // (edgeId,time,edge)
	edgeId     map[int]int
	remain     map[int]struct{}
}

// 离线动态连通性查询
//  n: 顶点 0~n-1
//  q: 查询 0~q-1
func NewOfflineDynamicConnectivity(n, q int) *OfflineDynamicConnectivity {
	log := bits.Len(uint(q - 1))
	size := 1 << log
	seg := make([][]int, size<<1)
	uf := NewUndoDSU(n)
	return &OfflineDynamicConnectivity{
		Uf:     uf,
		Part:   n,
		n:      n,
		q:      q,
		size:   size,
		seg:    seg,
		edges:  [][]int{},
		edgeId: map[int]int{},
		remain: map[int]struct{}{},
	}
}

// 时刻time添加一条无向边(u,v)
func (o *OfflineDynamicConnectivity) AddEdge(u, v, time int) {
	if u > v {
		u, v = v, u
	}
	tuple := u*o.n + v
	o.remain[tuple] = struct{}{}
	o.edgeId[tuple] = time
}

// 时刻time删除一条无向边(u,v)
func (o *OfflineDynamicConnectivity) RemoveEdge(u, v, time int) {
	if u > v {
		u, v = v, u
	}
	tuple := u*o.n + v
	delete(o.remain, tuple)
	o.edges = append(o.edges, []int{o.edgeId[tuple], time, tuple})
}

func (o *OfflineDynamicConnectivity) Build() {
	for e := range o.remain {
		o.edges = append(o.edges, []int{o.edgeId[e], o.q, e})
	}
	for i := range o.edges {
		l, r, e := o.edges[i][0], o.edges[i][1], o.edges[i][2]
		o.add(l, r, e)
	}
}

// 执行所有查询
//  cb: func(k int): 当前处于第k个查询(0-based)
func (o *OfflineDynamicConnectivity) Run(cb func(k int)) {
	var dfs func(cur int)
	dfs = func(cur int) {
		if o.size+o.q <= cur {
			return
		}

		add := 0
		for _, e := range o.seg[cur] {
			ok := o.Uf.Union(e/o.n, e%o.n) // AddEdge
			if ok {
				add++
			}
		}
		o.Part -= add

		if cur >= o.size {
			cb(cur - o.size)
		} else {
			dfs(cur << 1)
			dfs(cur<<1 | 1)
		}

		for i := 0; i < len(o.seg[cur]); i++ {
			o.Uf.Undo() // RemoveEdge
		}
		o.Part += add
	}
	dfs(1)
}

func (o *OfflineDynamicConnectivity) add(l, r, e int) {
	l += o.size
	r += o.size
	for l < r {
		if l&1 == 1 {
			o.seg[l] = append(o.seg[l], e)
			l++
		}
		if r&1 == 1 {
			r--
			o.seg[r] = append(o.seg[r], e)
		}
		l >>= 1
		r >>= 1
	}
}

func NewUndoDSU(n int) *UndoDSU { return NewUndoDSUWithWeights(make([]S, n)) }
func NewUndoDSUWithWeights(weights []S) *UndoDSU {
	n := len(weights)
	ps, ws := make([]int, n), make([]S, n)
	for i := 0; i < n; i++ {
		ps[i] = -1
		ws[i] = weights[i] // e()
	}
	history := [][]int{}
	return &UndoDSU{parentSize: ps, weights: ws, history: history}
}

type UndoDSU struct {
	parentSize []int
	weights    []S
	history    [][]int
}

func (uf *UndoDSU) Add(index int, delta S) {
	x := index
	for x >= 0 {
		uf.weights[x] += delta // op()
		x = uf.parentSize[x]
	}
}

func (uf *UndoDSU) Get(index int) S { return uf.weights[uf.Find(index)] }

func (uf *UndoDSU) Undo() bool {
	if len(uf.history) == 0 {
		return false
	}
	lastY := uf.history[len(uf.history)-1]
	uf.history = uf.history[:len(uf.history)-1]
	lastX := uf.history[len(uf.history)-1]
	uf.history = uf.history[:len(uf.history)-1]
	y, py := lastY[0], lastY[1]
	x, px := lastX[0], lastX[1]
	if uf.parentSize[x] != px {
		uf.weights[x] -= uf.weights[y] // inv()
	}
	uf.parentSize[x] = px
	uf.parentSize[y] = py
	return true
}

func (uf *UndoDSU) Reset() {
	for len(uf.history) > 0 {
		uf.Undo()
	}
}

func (uf *UndoDSU) Find(x int) int {
	cur := x
	for uf.parentSize[cur] >= 0 {
		cur = uf.parentSize[cur]
	}
	return cur
}

func (uf *UndoDSU) Union(x, y int) bool {
	x, y = uf.Find(x), uf.Find(y)
	if -uf.parentSize[x] < -uf.parentSize[y] {
		x, y = y, x
	}
	uf.history = append(uf.history, []int{x, uf.parentSize[x]})
	uf.history = append(uf.history, []int{y, uf.parentSize[y]})
	if x == y {
		return false
	}
	uf.parentSize[x] += uf.parentSize[y]
	uf.parentSize[y] = x
	uf.weights[x] += uf.weights[y]
	return true
}

func (uf *UndoDSU) IsConnected(x, y int) bool { return uf.Find(x) == uf.Find(y) }
