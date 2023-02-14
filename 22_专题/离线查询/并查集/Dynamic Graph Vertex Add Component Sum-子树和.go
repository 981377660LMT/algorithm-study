// Dynamic Graph Vertex Add Component Sum
// 连接断开边/单点修改权值/查询子树和

// 0 root1 root2 root3 root4 断开(root1-root2) 连接(root3-root4)
// 1 root x 将root的值加上x
// 2 root1 root2 输出root1所在子树的和,其中root2是root1的父亲节点

// 离线查询
// n<=2e5
// !技巧:查询子树和=断开父亲边+查询连通块和+连回父亲边

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	sums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &sums[i])
	}
	edges := make([][]int, 0, n-1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		edges = append(edges, []int{u, v})
	}
	queries := make([][]int, 0, q)
	for i := 0; i < q; i++ {
		var op, root1, root2, root3, root4, add, parent int
		fmt.Fscan(in, &op)
		if op == 0 {
			fmt.Fscan(in, &root1, &root2, &root3, &root4)
			queries = append(queries, []int{op, root1, root2, root3, root4})
		} else if op == 1 {
			fmt.Fscan(in, &root1, &add)
			queries = append(queries, []int{op, root1, add})
		} else {
			fmt.Fscan(in, &root1, &parent)
			queries = append(queries, []int{op, root1, parent})
		}
	}

	odc := NewOfflineDynamicConnectivity(n, q)
	for i, v := range sums {
		odc.Uf.Add(i, v) // 开始时每个点的权值
	}
	for _, edge := range edges {
		u, v := edge[0], edge[1]
		odc.AddEdge(u, v, 0) // 最开始时的边
	}

	for i, query := range queries { // 处理边的添加和删除
		op := query[0]
		if op == 0 {
			root1, root2, root3, root4 := query[1], query[2], query[3], query[4]
			odc.RemoveEdge(root1, root2, i)
			odc.AddEdge(root3, root4, i)
		} else if op == 2 {
			root, parent := query[1], query[2]
			odc.RemoveEdge(parent, root, i)
			odc.AddEdge(parent, root, i+1) // 下个时刻加回来
		}
	}
	odc.Build()

	// 处理权值查询
	res := []int{}
	f := func(k int) {
		op := queries[k][0]
		if op == 1 {
			root, add := queries[k][1], queries[k][2]
			odc.Uf.Add(root, add)
		} else if op == 2 {
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
	Uf         *UndoDSU // todo:interface
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
	stack := []int{1}
	for len(stack) > 0 {
		k := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if k >= 0 {
			if o.size+o.q <= k {
				continue
			}
			stack = append(stack, ^k)
			for _, e := range o.seg[k] {
				o.Uf.Union(e/o.n, e%o.n)
			}
			if o.size <= k {
				cb(k - o.size)
			} else {
				stack = append(stack, k<<1|1, k<<1)
			}
		} else {
			for i := 0; i < len(o.seg[^k]); i++ {
				o.Uf.Undo()
			}
		}
	}
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
		ws[i] = weights[i]
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
		uf.weights[x] += delta
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
		uf.weights[x] -= uf.weights[y]
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
