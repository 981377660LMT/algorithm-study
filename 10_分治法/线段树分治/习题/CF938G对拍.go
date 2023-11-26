package main

import (
	"container/list"
	"fmt"
	"math/bits"
	"math/rand"
	"sort"
	"strings"
)

const INF int = 1e18

// 给出一个连通带权无向图,边有边权,要求支持 q 个操作:
// 1 u v w : 在原图中加入一条 u 到 v 的边, 边权为 w.
// 2 u v : 把 u 到 v 的边删除.
// 3 u v : 询问 u 到 v 的异或最短路.
// 保证任意时刻图是连通的.
// 1<=n,m,q<=2e5
func main() {

	n := 2
	q := 3

	checkIsConnected := func(edges [][3]int) bool {
		uf := NewUnionFindArray(n)
		for _, edge := range edges {
			u, v := edge[0], edge[1]
			uf.Union(u, v)
		}
		return uf.Part == 1
	}

	edgeMap := map[[2]int]struct{}{}
	edges := [][3]int{}
	addEdge := func(u, v, w int) bool {
		if u == v {
			return false
		}
		if u > v {
			u, v = v, u
		}
		if _, ok := edgeMap[[2]int{u, v}]; !ok {
			edgeMap[[2]int{u, v}] = struct{}{}
			return true
		} else {
			return false
		}
	}

	removeEdge := func(u, v int) bool {
		if u == v {
			return false
		}
		if u > v {
			u, v = v, u
		}
		if _, ok := edgeMap[[2]int{u, v}]; ok {
			delete(edgeMap, [2]int{u, v})
			newEdges := [][3]int{}
			for hash := range edgeMap {
				u, v := hash[0], hash[1]
				newEdges = append(newEdges, [3]int{u, v, 0})
			}
			if checkIsConnected(newEdges) {
				return true
			} else {
				edgeMap[[2]int{u, v}] = struct{}{}
				return false
			}
		} else {
			return false
		}
	}

	for i := 0; i < n-1; i++ {
		w := rand.Intn(100)
		addEdge(i, i+1, w)
		edges = append(edges, [3]int{i, i + 1, w})
	}
	for i := 0; i < n; i++ {
		u, v, w := rand.Intn(n), rand.Intn(n), rand.Intn(100)
		if addEdge(u, v, w) {
			edges = append(edges, [3]int{u, v, w})
		}
	}

	if !checkIsConnected(edges) {
		panic("not connected")
	}

	operations := make([][4]int, q)
	for i := range operations {
		op := rand.Intn(3) + 1
		if op == 1 {
			for {
				u, v, w := rand.Intn(n), rand.Intn(n), rand.Intn(100)
				if addEdge(u, v, w) {
					operations[i] = [4]int{op, u, v, w}
					break
				}
			}
		} else if op == 2 && len(edgeMap) > 0 {
			for {
				u, v := rand.Intn(n), rand.Intn(n)
				if removeEdge(u, v) {
					operations[i] = [4]int{op, u, v, 0}
					break
				}
			}
		} else {
			for {
				u, v := rand.Intn(n), rand.Intn(n)
				if u != v {
					operations[i] = [4]int{op, u, v, 0}
					break
				}
			}
		}
	}

	// 	5 5
	// 1 2 3
	// 2 3 4
	// 3 4 5
	// 4 5 6
	// 1 5 1
	// 5
	// 3 1 5
	// 1 1 3 1
	// 3 1 5
	// 2 1 5
	// 3 1 5
	reset := func() {
		n, q = 5, 5
		edges = [][3]int{{0, 1, 3}, {1, 2, 4}, {2, 3, 5}, {3, 4, 6}, {0, 4, 1}}
		operations = [][4]int{{3, 0, 4, 0}, {1, 0, 2, 1}, {3, 0, 4, 0}, {2, 0, 4, 0}, {3, 0, 4, 0}}
	}
	_ = reset
	reset()

	R := NewAddRemoveQuery(false)
	seg := NewSegmentTreeDivideAndConquerCopy()
	edgeWeight := make(map[int]int)
	for i := range edges {
		u, v, w := edges[i][0], edges[i][1], edges[i][2]
		if u > v {
			u, v = v, u
		}
		edgeHash := u*n + v
		edgeWeight[edgeHash] = w
		R.Add(-1, edgeHash) // 初始时刻的边
	}

	queryEdges := make([][2]int, 0, q)
	for i, query := range operations {
		op, u, v, w := query[0], query[1], query[2], query[3]
		if u > v {
			u, v = v, u
		}
		edgeHash := u*n + v
		if op == 1 {
			R.Add(len(queryEdges), edgeHash)
			edgeWeight[edgeHash] = w
		} else if op == 2 {
			R.Remove(len(queryEdges), edgeHash)
		} else {
			seg.AddQuery(len(queryEdges), i)
			queryEdges = append(queryEdges, [2]int{u, v})
		}
	}
	events := R.GetEvents(INF)
	for _, e := range events {
		seg.AddMutation(e.start, e.end, e.value) // value为edgeHash
	}

	res := make([]int, q)
	for i := range res {
		res[i] = INF
	}

	uf := NewUnionFindWithDistAndUndo(n)

	initVectorSpace := NewVectorSpace(nil)
	seg.Run(
		&initVectorSpace,
		func(vectorSpace *State, edgeHash int) {
			vs := *vectorSpace
			u, v := edgeHash/n, edgeHash%n
			w := edgeWeight[edgeHash]
			// fmt.Println("add", u, v, w, *vectorSpace)
			root1, x1 := uf.Find(u)
			root2, x2 := uf.Find(v)
			uf.Union(u, v, w)
			if root1 != root2 {
			} else {
				cycleXor := x1 ^ x2 ^ w
				vs.Add(cycleXor)
			}
		},
		func(vectorSpace *State) *State {
			tmp := *vectorSpace
			// fmt.Println("copy", tmp)
			res := tmp.Copy()
			return &res
		},
		func(vectorSpace *State, queryId int) {
			vs := *vectorSpace
			if queryId >= 0 && queryId < len(operations) {
				query := operations[queryId]
				u, v := query[1], query[2]
				// fmt.Println("query", u, v, vs)

				dist := uf.Dist(u, v)
				// fmt.Println(dist, uf.GetGroups())
				res[queryId] = vs.Min(dist)
			}
		},
		func() {
			uf.Undo()
			// fmt.Println("undo", uf.data.history)
		},
	)

	for _, v := range res {
		if v != INF {
			fmt.Println(v)
		}
	}
}

type State = *VectorSpace

type segMutation struct{ start, end, id int }
type segQuery struct{ time, id int }

// 线段树分治copy版.
// 如果修改操作难以撤销，可以在每个节点处保存一份副本.
// !调用O(n)次拷贝注意不要超出内存.
type SegmentTreeDivideAndConquerCopy struct {
	initState  *State
	mutate     func(state *State, mutationId int)
	copy       func(state *State) *State
	query      func(state *State, queryId int)
	undo       func() // 可选属性
	mutations  []segMutation
	operations []segQuery
	nodes      [][]int
}

func NewSegmentTreeDivideAndConquerCopy() *SegmentTreeDivideAndConquerCopy {
	return &SegmentTreeDivideAndConquerCopy{}
}

// 在时间范围`[startTime, endTime)`内添加一个编号为`id`的变更.
func (o *SegmentTreeDivideAndConquerCopy) AddMutation(startTime, endTime int, id int) {
	if startTime >= endTime {
		return
	}
	o.mutations = append(o.mutations, segMutation{startTime, endTime, id})
}

// 在时间`time`时添加一个编号为`id`的查询.
func (o *SegmentTreeDivideAndConquerCopy) AddQuery(time int, id int) {
	o.operations = append(o.operations, segQuery{time, id})
}

// dfs 遍历整棵线段树来得到每个时间点的答案.
//
//		initState: 数据结构的初始状态.
//		mutate: 添加编号为`mutationId`的变更后产生的副作用.
//		copy: 拷贝一份数据结构的副本.
//		query: 响应编号为`queryId`的查询.
//	  undo: 可选属性，如果不为nil，会在每次回溯时调用.
//
// 一共调用 **O(nlogn)** 次`mutate`，**O(n)** 次`copy` ,**O(q)** 次`query` 和 **O(nlogn)** 次`undo`.
func (o *SegmentTreeDivideAndConquerCopy) Run(
	initState *State,
	mutate func(state *State, mutationId int),
	copy func(state *State) *State,
	query func(state *State, queryId int),
	undo func(),
) {
	if len(o.operations) == 0 {
		return
	}
	o.initState = initState
	o.mutate, o.copy, o.query, o.undo = mutate, copy, query, undo
	times := make([]int, len(o.operations))
	for i := range o.operations {
		times[i] = o.operations[i].time
	}
	sort.Ints(times)
	uniqueInplace(&times)
	usedTimes := make([]bool, len(times)+1)
	usedTimes[0] = true
	for _, e := range o.mutations {
		usedTimes[lowerBound(times, e.start)] = true
		usedTimes[lowerBound(times, e.end)] = true
	}
	for i := 1; i < len(times); i++ {
		if !usedTimes[i] {
			times[i] = times[i-1]
		}
	}
	uniqueInplace(&times)

	n := len(times)
	offset := 1
	for offset < n {
		offset <<= 1
	}
	o.nodes = make([][]int, offset+n)
	for _, e := range o.mutations {
		left := offset + lowerBound(times, e.start)
		right := offset + lowerBound(times, e.end)
		eid := e.id << 1
		for left < right {
			if left&1 == 1 {
				o.nodes[left] = append(o.nodes[left], eid)
				left++
			}
			if right&1 == 1 {
				right--
				o.nodes[right] = append(o.nodes[right], eid)
			}
			left >>= 1
			right >>= 1
		}
	}

	for _, q := range o.operations {
		pos := offset + upperBound(times, q.time) - 1
		o.nodes[pos] = append(o.nodes[pos], (q.id<<1)|1)
	}

	o.dfs(1, o.initState)
}

func (o *SegmentTreeDivideAndConquerCopy) dfs(now int, state *State) {
	curNodes := o.nodes[now]
	for _, id := range curNodes {
		if id&1 == 1 {
			o.query(state, id>>1)
		} else {
			o.mutate(state, id>>1)
		}
	}

	if now<<1 < len(o.nodes) {
		o.dfs(now<<1, o.copy(state))
	}
	if (now<<1)|1 < len(o.nodes) {
		o.dfs((now<<1)|1, o.copy(state))
	}

	if o.undo != nil {
		for i := len(curNodes) - 1; i >= 0; i-- {
			if curNodes[i]&1 == 0 {
				o.undo()
			}
		}
	}
}

func uniqueInplace(sorted *[]int) {
	tmp := *sorted
	slow := 0
	for fast := 0; fast < len(tmp); fast++ {
		if tmp[fast] != tmp[slow] {
			slow++
			tmp[slow] = tmp[fast]
		}
	}
	*sorted = tmp[:slow+1]
}

func lowerBound(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := (left + right) >> 1
		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func upperBound(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := (left + right) >> 1
		if arr[mid] <= target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

type V = int
type Event = struct {
	start int
	end   int
	value V
}

// 将时间轴上单点的 add 和 remove 转化为区间上的 add.
// !不能加入相同的元素，删除时元素必须要存在。
// 如果 add 和 remove 按照时间顺序单增，那么可以使用 monotone = true 来加速。
type AddRemoveQuery struct {
	mp       *LinkedHashMap
	events   []Event
	adds     map[V][]int
	addKeys  []V
	removes  map[V][]int
	monotone bool
}

func NewAddRemoveQuery(monotone bool) *AddRemoveQuery {
	return &AddRemoveQuery{
		mp:       NewLinkedHashMap(16),
		events:   []Event{},
		adds:     map[V][]int{},
		removes:  map[V][]int{},
		monotone: monotone,
	}
}

func (adq *AddRemoveQuery) Add(time int, value V) {
	if adq.monotone {
		adq.addMonotone(time, value)
	} else {
		adq.adds[value] = append(adq.adds[value], time)
		adq.addKeys = append(adq.addKeys, value)
	}
}

func (adq *AddRemoveQuery) Remove(time int, value V) {
	if adq.monotone {
		adq.removeMonotone(time, value)
	} else {
		adq.removes[value] = append(adq.removes[value], time)
	}
}

// lastTime: 所有变更都结束的时间.例如INF.
func (adq *AddRemoveQuery) GetEvents(lastTime int) []Event {
	if adq.monotone {
		return adq.getMonotone(lastTime)
	}

	res := []Event{}
	// visited := make(map[V]struct{}, len(adq.adds))
	for value, addTimes := range adq.adds {
		// if _, v := visited[value]; v {
		// 	continue
		// }
		// visited[value] = struct{}{}
		// addTimes := adq.adds[value]
		removeTimes := []int{}
		if cand, ok := adq.removes[value]; ok {
			removeTimes = cand
			delete(adq.removes, value)
		}
		if len(removeTimes) < len(addTimes) {
			removeTimes = append(removeTimes, lastTime)
		}
		sort.Ints(addTimes)
		sort.Ints(removeTimes)
		for i, t := range addTimes {
			if t < removeTimes[i] {
				res = append(res, Event{t, removeTimes[i], value})
			}
		}
	}

	return res
}

func (adq *AddRemoveQuery) addMonotone(time int, value V) {
	if adq.mp.Has(value) {
		panic("can't add a value that already exists")
	}
	adq.mp.Set(value, time)
}

func (adq *AddRemoveQuery) removeMonotone(time int, value V) {
	if startTime, ok := adq.mp.Get(value); !ok {
		panic("can't remove a value that doesn't exist")
	} else {
		adq.mp.Delete(value)
		if startTime != time {
			adq.events = append(adq.events, Event{startTime.(int), time, value})
		}
	}
}

func (adq *AddRemoveQuery) getMonotone(lastTime int) []Event {
	adq.mp.ForEach(func(value interface{}, startTime interface{}) bool {
		if startTime.(int) != lastTime {
			adq.events = append(adq.events, Event{startTime.(int), lastTime, value.(int)})
		}
		return false
	})
	return adq.events
}

// VectorSpace，线性基空间.支持线性基合并.
type VectorSpace struct {
	bases  []int
	maxBit int
}

func NewVectorSpace(nums []int) *VectorSpace {
	res := &VectorSpace{}
	for _, num := range nums {
		res.Add(num)
	}
	return res
}

// 插入一个向量,如果插入成功(不能被表出)返回True,否则返回False.
func (lb *VectorSpace) Add(num int) bool {
	for _, base := range lb.bases {
		if base == 0 || num == 0 {
			break
		}
		num = min(num, num^base)
	}
	if num != 0 {
		lb.bases = append(lb.bases, num)
		lb.maxBit = max(lb.maxBit, num)
		return true
	}
	return false
}

// 求xor与所有向量异或的最大值.
func (lb *VectorSpace) Max(xor int) int {
	res := xor
	for _, base := range lb.bases {
		res = max(res, res^base)
	}
	return res
}

// 求xor与所有向量异或的最小值.
func (lb *VectorSpace) Min(xor int) int {
	res := xor
	for _, base := range lb.bases {
		res = min(res, res^base)
	}
	return res
}

func (lb *VectorSpace) Copy() *VectorSpace {
	res := &VectorSpace{}
	res.bases = append(res.bases, lb.bases...)
	res.maxBit = lb.maxBit
	return res
}

func (lb *VectorSpace) Len() int {
	return len(lb.bases)
}

func (lb *VectorSpace) ForEach(f func(base int)) {
	for _, base := range lb.bases {
		f(base)
	}
}

func (lb *VectorSpace) Has(v int) bool {
	for _, w := range lb.bases {
		if v == 0 {
			break
		}
		v = min(v, v^w)
	}
	return v == 0
}

// Merge.
func (lb *VectorSpace) Or(other *VectorSpace) *VectorSpace {
	v1, v2 := lb, other
	if v1.Len() < v2.Len() {
		v1, v2 = v2, v1
	}
	res := v1.Copy()
	for _, base := range v2.bases {
		res.Add(base)
	}
	return res
}

func (lb *VectorSpace) And(other *VectorSpace) *VectorSpace {
	maxDim := max(lb.maxBit, other.maxBit)
	x := lb.orthogonalSpace(maxDim)
	y := other.orthogonalSpace(maxDim)
	if x.Len() < y.Len() {
		x, y = y, x
	}
	for _, base := range y.bases {
		x.Add(base)
	}
	return x.orthogonalSpace(maxDim)
}

func (lb *VectorSpace) String() string {
	return fmt.Sprintf("%v", lb.bases)
}

// 正交空间.
func (lb *VectorSpace) orthogonalSpace(maxDim int) *VectorSpace {
	lb.normalize(true)
	m := maxDim
	tmp := make([]int, m)
	for _, base := range lb.bases {
		tmp[bits.Len(uint(base))-1] = base
	}
	tmp = Transpose(m, m, tmp, true)
	res := &VectorSpace{}
	for j, v := range tmp {
		if v>>j&1 == 1 {
			continue
		}
		res.Add(v | 1<<j)
	}
	return res
}

func (lb *VectorSpace) normalize(reverse bool) {
	for j, v := range lb.bases {
		for i := 0; i < j; i++ {
			lb.bases[i] = min(lb.bases[i], lb.bases[i]^v)
		}
	}
	if !reverse {
		sort.Ints(lb.bases)
	} else {
		sort.Sort(sort.Reverse(sort.IntSlice(lb.bases)))
	}
}

// 矩阵转置,O(n+m)log(n+m)
func Transpose(row, col int, grid []int, inPlace bool) []int {
	if len(grid) != row {
		panic("row not match")
	}
	if !inPlace {
		grid = append(grid[:0:0], grid...)
	}
	log := 0
	max_ := max(row, col)
	for 1<<log < max_ {
		log++
	}
	if len(grid) < 1<<log {
		*&grid = append(grid, make([]int, 1<<log-len(grid))...)
	}
	width := 1 << log
	mask := int(1)
	for i := 0; i < log; i++ {
		mask |= (mask << (1 << i))
	}
	for t := 0; t < log; t++ {
		width >>= 1
		mask ^= (mask >> width)
		for i := 0; i < 1<<t; i++ {
			for j := 0; j < width; j++ {
				x := &grid[width*(2*i)+j]
				y := &grid[width*(2*i+1)+j]
				*x = ((*y << width) & mask) ^ *x
				*y = ((*x & mask) >> width) ^ *y
				*x = ((*y << width) & mask) ^ *x
			}
		}
	}
	return grid[:col]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//
//

type T = int

func e() T        { return 0 }
func op(x, y T) T { return x ^ y }
func inv(x T) T   { return x }

// 维护到每个组根节点距离的可撤销并查集.
// 用于维护环的权值，树上的距离等.
type UnionFindWithDistAndUndo struct {
	data *RollbackArray
}

func NewUnionFindWithDistAndUndo(n int) *UnionFindWithDistAndUndo {
	return &UnionFindWithDistAndUndo{
		data: NewRollbackArray(n, func(index int) arrayItem { return arrayItem{parent: -1, dist: e()} }),
	}
}

// distToRoot(parent) + dist = distToRoot(child).
func (uf *UnionFindWithDistAndUndo) Union(parent int, child int, dist T) bool {
	v1, x1 := uf.Find(parent)
	v2, x2 := uf.Find(child)
	if v1 == v2 {
		uf.data.Set(v2, uf.data.Get(v2))
		uf.data.Set(v1, uf.data.Get(v1))
		return dist == op(x2, inv(x1))
	}
	s1, s2 := -uf.data.Get(v1).parent, -uf.data.Get(v2).parent
	if s1 < s2 {
		s1, s2 = s2, s1
		v1, v2 = v2, v1
		x1, x2 = x2, x1
		dist = inv(dist)
	}
	// v1 <- v2
	dist = op(x1, dist)
	dist = op(dist, inv(x2))
	uf.data.Set(v2, arrayItem{parent: v1, dist: dist})
	uf.data.Set(v1, arrayItem{parent: -(s1 + s2), dist: e()})
	return true
}

// 返回v所在组的根节点和到v到根节点的距离.
func (uf *UnionFindWithDistAndUndo) Find(v int) (root int, distToRoot T) {
	root, distToRoot = v, e()
	for {
		item := uf.data.Get(root)
		if item.parent < 0 {
			break
		}
		distToRoot = op(distToRoot, item.dist)
		root = item.parent
	}
	return
}

// Dist(x, y) = DistToRoot(x) - DistToRoot(y).
// 如果x和y不在同一个集合,抛出错误.
func (uf *UnionFindWithDistAndUndo) Dist(x int, y int) T {
	vx, dx := uf.Find(x)
	vy, dy := uf.Find(y)
	if vx != vy {
		panic("x and y are not in the same set")
	}
	return op(dx, inv(dy))
}

func (uf *UnionFindWithDistAndUndo) DistToRoot(x int) T {
	_, dx := uf.Find(x)
	return dx
}

func (uf *UnionFindWithDistAndUndo) GetTime() int {
	return uf.data.GetTime()
}

func (uf *UnionFindWithDistAndUndo) Rollback(time int) {
	uf.data.Rollback(time)
}

func (uf *UnionFindWithDistAndUndo) Undo() {
	uf.data.Undo()
	uf.data.Undo()
}

func (uf *UnionFindWithDistAndUndo) GetSize(x int) int {
	root, _ := uf.Find(x)
	return -uf.data.Get(root).parent
}

func (uf *UnionFindWithDistAndUndo) GetGroups() map[int][]int {
	res := make(map[int][]int)
	for i := 0; i < uf.data.Len(); i++ {
		root, _ := uf.Find(i)
		res[root] = append(res[root], i)
	}
	return res
}

type arrayItem struct {
	parent int
	dist   T
}

type historyItem struct {
	index int
	value arrayItem
}

type RollbackArray struct {
	n       int
	data    []arrayItem
	history []historyItem
}

func NewRollbackArray(n int, f func(index int) arrayItem) *RollbackArray {
	data := make([]arrayItem, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
	}
	return &RollbackArray{
		n:    n,
		data: data,
	}
}

func (r *RollbackArray) GetTime() int {
	return len(r.history)
}

func (r *RollbackArray) Rollback(time int) {
	for len(r.history) > time {
		pair := r.history[len(r.history)-1]
		r.history = r.history[:len(r.history)-1]
		r.data[pair.index] = pair.value
	}
}

func (r *RollbackArray) Undo() bool {
	if len(r.history) == 0 {
		return false
	}
	pair := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]
	r.data[pair.index] = pair.value
	return true
}

func (r *RollbackArray) Get(index int) arrayItem {
	return r.data[index]
}

func (r *RollbackArray) Set(index int, value arrayItem) {
	r.history = append(r.history, historyItem{index: index, value: r.data[index]})
	r.data[index] = value
}

func (r *RollbackArray) GetAll() []arrayItem {
	return append(r.data[:0:0], r.data...)
}

func (r *RollbackArray) Len() int {
	return r.n
}

type listItem = struct{ key, value interface{} }
type LinkedHashMap struct {
	mp   map[interface{}]*list.Element
	list *list.List
}

func NewLinkedHashMap(initCapacity int) *LinkedHashMap {
	return &LinkedHashMap{make(map[interface{}]*list.Element, initCapacity), list.New()}
}

func (s *LinkedHashMap) Get(key interface{}) (res interface{}, ok bool) {
	if ele, hit := s.mp[key]; hit {
		return ele.Value.(listItem).value, true
	}
	return
}

func (s *LinkedHashMap) Set(key, value interface{}) *LinkedHashMap {
	if ele, hit := s.mp[key]; hit {
		ele.Value = listItem{key, value}
	} else {
		s.mp[key] = s.list.PushBack(listItem{key, value})
	}
	return s
}

func (s *LinkedHashMap) Delete(key interface{}) bool {
	if ele, hit := s.mp[key]; hit {
		s.list.Remove(ele)
		delete(s.mp, key)
		return true
	}
	return false
}

func (s *LinkedHashMap) Has(key interface{}) bool {
	_, ok := s.mp[key]
	return ok
}

func (s *LinkedHashMap) Size() int {
	return len(s.mp)
}

// 按照插入顺序遍历哈希表中的元素
// 当 f 返回 true 时停止遍历
func (s *LinkedHashMap) ForEach(f func(key interface{}, value interface{}) bool) {
	for node := s.list.Front(); node != nil; node = node.Next() {
		if f(node.Value.(listItem).key, node.Value.(listItem).value) {
			break
		}
	}
}

func (s *LinkedHashMap) String() string {
	res := []string{"LinkedHashMap{"}
	content := []string{}
	s.ForEach(func(key interface{}, value interface{}) bool {
		content = append(content, fmt.Sprintf("%v: %v", key, value))
		return false
	})
	res = append(res, strings.Join(content, ", "))
	res = append(res, "}")
	return strings.Join(res, "")
}

func NewUnionFindArray(n int) *_UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_UnionFindArray{
		Part:   n,
		size:   n,
		rank:   rank,
		parent: parent,
	}
}

type _UnionFindArray struct {
	size   int
	Part   int
	rank   []int
	parent []int
}

func (ufa *_UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *_UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *_UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}
