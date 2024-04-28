// api：
//
//	单点修改边权/点权
//	查询路径点权和边权和
//	查询子树点权和边权和

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	demo()
}

func demo() {
	//   0
	//  / \
	// 1   2
	//    / \
	//   3   4
	//      /
	//     5

	edges := [][3]int32{{0, 1, 1}, {0, 2, 2}, {2, 3, 3}, {2, 4, 4}, {4, 5, 5}}
	adjList := make([][]neighbor, 6)
	for _, e := range edges {
		u, v, w := e[0], e[1], e[2]
		adjList[u] = append(adjList[u], neighbor{v, int(w)})
		adjList[v] = append(adjList[v], neighbor{u, int(w)})
	}
	vertexCost := []int{0, 1, 2, 3, 4, 5}
	et := NewEulerTour(adjList, 0, vertexCost)
	_ = et
	fmt.Println(et.st.Query(0, 6))
	fmt.Println(et.Lca(4, 5))
	fmt.Println(et.LcaMulti([]int32{4, 5, 3}))
	fmt.Println(et.SubtreeVcost(0), et.SubtreeEcost(0))
	fmt.Println(et.SubtreeVcost(1), et.SubtreeEcost(1))
	fmt.Println(et.SubtreeVcost(2), et.SubtreeEcost(2))
	fmt.Println(et.SubtreeVcost(3), et.SubtreeEcost(3))
	fmt.Println(et.SubtreeVcost(4), et.SubtreeEcost(4))
	fmt.Println(et.SubtreeVcost(5), et.SubtreeEcost(5))

	et.AddVertex(0, 1)
	fmt.Println(et.SubtreeVcost(0), et.SubtreeEcost(0))
	et.AddVertex(1, 1)
	fmt.Println(et.SubtreeVcost(1), et.SubtreeEcost(1))
	et.SetVertex(1, 2)
	fmt.Println(et.SubtreeVcost(1), et.SubtreeEcost(1))

	et.AddEdge(0, 1, 1)
	fmt.Println(et.SubtreeVcost(0), et.SubtreeEcost(0))
}

type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }
func inv(a E) E   { return -a }

type EulerTour struct {
	n                                                int32
	depth, nodeIn, nodeOut                           []int32
	vertexCost                                       []E
	path                                             []int32
	vCostSubtree, vCostPath, eCostSubtree, eCostPath *bitGroup32
	mask                                             int32
	st                                               *SegmentTree32
}

type neighbor struct {
	next int32
	cost E
}

// root为-1表示无根.
func NewEulerTour(adjList [][]neighbor, root int32, vertexCost []E) *EulerTour {
	n := int32(len(adjList))
	if vertexCost == nil {
		vertexCost = make([]E, n)
		for i := range vertexCost {
			vertexCost[i] = e()
		}
	}

	path := make([]int32, 2*n)
	vCost1 := make([]E, 2*n)
	vCost2 := make([]E, 2*n)
	eCost1 := make([]E, 2*n)
	eCost2 := make([]E, 2*n)
	nodeIn := make([]int32, n)
	nodeOut := make([]int32, n)
	depth := make([]int32, n)
	for i := range depth {
		depth[i] = -1
	}
	curTime := int32(-1)

	dfs := func(curRoot int32) {
		depth[curRoot] = 0
		stack := []neighbor{{next: ^curRoot, cost: 0}, {next: curRoot, cost: 0}}
		for len(stack) > 0 {
			curTime++
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			v, ec := top.next, top.cost
			if v >= 0 {
				nodeIn[v] = curTime
				path[curTime] = v
				eCost1[curTime] = ec
				eCost2[curTime] = ec
				vCost1[curTime] = vertexCost[v]
				vCost2[curTime] = vertexCost[v]
				if len(adjList[v]) == 1 {
					nodeOut[v] = curTime + 1
				}
				for _, e := range adjList[v] {
					x := e.next
					c := e.cost
					if depth[x] != -1 {
						continue
					}
					depth[x] = depth[v] + 1
					stack = append(stack, neighbor{next: ^v, cost: c})
					stack = append(stack, neighbor{next: x, cost: c})
				}
			} else {
				v = ^v
				path[curTime] = v
				eCost1[curTime] = e()
				eCost2[curTime] = inv(ec)
				vCost1[curTime] = e()
				vCost2[curTime] = inv(vertexCost[v])
				nodeOut[v] = curTime
			}
		}
	}

	if root >= 0 {
		dfs(root)
	} else {
		for i := int32(0); i < n; i++ {
			if depth[i] == -1 {
				dfs(i)
			}
		}
	}

	res := &EulerTour{}
	res.n = n
	res.depth = depth
	res.nodeIn = nodeIn
	res.nodeOut = nodeOut
	res.vertexCost = vertexCost
	res.path = path
	res.vCostSubtree = newBITGroup32From(2*n, func(i int32) E { return vCost1[i] })
	res.vCostPath = newBITGroup32From(2*n, func(i int32) E { return vCost2[i] })
	res.eCostSubtree = newBITGroup32From(2*n, func(i int32) E { return eCost1[i] })
	res.eCostPath = newBITGroup32From(2*n, func(i int32) E { return eCost2[i] })
	bit := bits.Len32(uint32(len(path)))
	res.mask = (1 << bit) - 1
	a := make([]int32, len(path))
	for i, v := range path {
		a[i] = (depth[v] << bit) + int32(i) // (深度，编号)
	}
	res.st = NewSegmentTree32From(a)
	return res
}

func (et *EulerTour) Lca(u, v int32) int32 {
	if u == v {
		return u
	}
	l := min32(et.nodeIn[u], et.nodeIn[v])
	r := max32(et.nodeOut[u], et.nodeOut[v])
	ind := et.st.Query(l, r) & et.mask
	return et.path[ind]
}

func (et *EulerTour) LcaMulti(points []int32) int32 {
	l, r := et.n+1, -et.n-1
	for _, e := range points {
		l = min32(l, et.nodeIn[e])
		r = max32(r, et.nodeOut[e])
	}
	ind := et.st.Query(l, r) & et.mask
	return et.path[ind]
}

func (et *EulerTour) SubtreeVcost(v int32) E {
	l := et.nodeIn[v]
	r := et.nodeOut[v]
	return et.vCostSubtree.QueryRange(l, r)
}

func (et *EulerTour) SubtreeEcost(v int32) E {
	l := et.nodeIn[v]
	r := et.nodeOut[v]
	return et.eCostSubtree.QueryRange(l+1, r)
}

func (et *EulerTour) PathVcost(u, v int32) E {
	lca := et.Lca(u, v)
	res := et._pathVcost(u)
	res = op(res, et._pathVcost(v))
	inv := inv(et._pathVcost(lca))
	res = op(res, inv)
	res = op(res, inv)
	res = op(res, et.vertexCost[lca])
	return res
}

func (et *EulerTour) PathEcost(u, v int32) E {
	lca := et.Lca(u, v)
	res := et._pathEcost(u)
	res = op(res, et._pathEcost(v))
	inv := inv(et._pathEcost(lca))
	res = op(res, inv)
	res = op(res, inv)
	return res
}

func (et *EulerTour) AddVertex(root int32, e E) {
	l := et.nodeIn[root]
	r := et.nodeOut[root]
	et.vCostSubtree.Update(l, e)
	et.vCostPath.Update(l, e)
	et.vCostPath.Update(r, inv(e))
	et.vertexCost[root] = op(et.vertexCost[root], e)
}

func (et *EulerTour) SetVertex(root int32, e E) {
	et.AddVertex(root, op(e, inv(et.vertexCost[root])))
}

func (et *EulerTour) AddEdge(u, v int32, e E) {
	if et.depth[u] < et.depth[v] {
		u, v = v, u
	}
	l := et.nodeIn[u]
	r := et.nodeOut[u]
	invE := inv(e)
	et.eCostSubtree.Update(l, e)
	et.eCostSubtree.Update(r+1, invE)
	et.eCostPath.Update(l, e)
	et.eCostPath.Update(r+1, invE)
}

func (et *EulerTour) SetEdge(u, v int32, e E) {
	et.AddEdge(u, v, op(e, inv(et.PathEcost(u, v))))
}

func (et *EulerTour) _pathVcost(v int32) E {
	return et.vCostPath.QueryPrefix(et.nodeIn[v] + 1)
}

func (et *EulerTour) _pathEcost(v int32) E {
	return et.eCostPath.QueryPrefix(et.nodeIn[v] + 1)
}

type bitGroup32 struct {
	n     int32
	data  []E
	total E
}

func newBITGroup32(n int32) *bitGroup32 {
	data := make([]E, n)
	for i := range data {
		data[i] = e()
	}
	return &bitGroup32{n: n, data: data, total: e()}
}

func newBITGroup32From(n int32, f func(index int32) E) *bitGroup32 {
	total := e()
	data := make([]E, n)
	for i := range data {
		data[i] = f(int32(i))
		total = op(total, data[i])
	}
	for i := int32(1); i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] = op(data[j-1], data[i-1])
		}
	}
	return &bitGroup32{n: n, data: data, total: total}
}

func (fw *bitGroup32) Update(i int32, x E) {
	fw.total = op(fw.total, x)
	for i++; i <= fw.n; i += i & -i {
		fw.data[i-1] = op(fw.data[i-1], x)
	}
}

func (fw *bitGroup32) QueryAll() E { return fw.total }

// [0, end)
func (fw *bitGroup32) QueryPrefix(end int32) E {
	if end > fw.n {
		end = fw.n
	}
	res := e()
	for end > 0 {
		res = op(res, fw.data[end-1])
		end &= end - 1
	}
	return res
}

// [start, end)
func (fw *bitGroup32) QueryRange(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > fw.n {
		end = fw.n
	}
	if start == 0 {
		return fw.QueryPrefix(end)
	}
	if start > end {
		return e()
	}
	pos, neg := e(), e()
	for end > start {
		pos = op(pos, fw.data[end-1])
		end &= end - 1
	}
	for start > end {
		neg = op(neg, fw.data[start-1])
		start &= start - 1
	}
	return op(pos, inv(neg))
}

// PointSetRangeMin

const INF32 int32 = 1 << 30

type SegData = int32

func (*SegmentTree32) e() SegData              { return INF32 }
func (*SegmentTree32) op(a, b SegData) SegData { return min32(a, b) }
func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

type SegmentTree32 struct {
	n, size int32
	seg     []SegData
}

func NewSegmentTree32(n int32, f func(int32) SegData) *SegmentTree32 {
	res := &SegmentTree32{}
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]SegData, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := int32(0); i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func NewSegmentTree32From(leaves []SegData) *SegmentTree32 {
	res := &SegmentTree32{}
	n := int32(len(leaves))
	size := int32(1)
	for size < n {
		size <<= 1
	}
	seg := make([]SegData, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := int32(0); i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *SegmentTree32) Get(index int32) SegData {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTree32) Set(index int32, value SegData) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree32) Update(index int32, value SegData) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTree32) Query(start, end int32) SegData {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}
func (st *SegmentTree32) QueryAll() SegData { return st.seg[1] }
func (st *SegmentTree32) GetAll() []SegData {
	res := make([]SegData, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}
