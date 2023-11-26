// https://www.luogu.com.cn/blog/yszs/xian-duan-shu-fen-zhi
// https://www.luogu.com.cn/problem/CF938G
// https://codeforces.com/contest/938/submission/184517092
// 每条无向边有边权和一段作用时间区间，查询两个点之间的异或最短路。
//
// !首先考虑怎么求两个点之间的异或最短路，其实就是P4151 最大XOR和路径
// 考虑如果一颗树的话，之间的答案就是路径上所有边异或起来的值。
// 那么推广到图上去，不难想到就是多了一些环，而且这个环走一圈是没贡献的（异或为零）。所以如果你要从环的一段走向另一端，实际上只有两种走法。

// 由于异或的优秀性质，可以将看上去复杂的图论问题转化成关于每条边的线性基上问题。
// 由于不支持撤销，我们在当前dfs栈中的每个节点上保存一份副本，供下一步的处理，最多保存logn个，可以接受。

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

const INF int = 1e18

// 给出一个连通带权无向图,边有边权,要求支持 q 个操作:
// 1 u v w : 在原图中加入一条 u 到 v 的边, 边权为 w.
// 2 u v : 把 u 到 v 的边删除.
// 3 u v : 询问 u 到 v 的异或最短路.
// 不会添加存在过的边.
// 保证任意时刻图是连通的.
// 1<=n,m,q<=2e5
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][3]int, m)
	for i := range edges {
		fmt.Fscan(in, &edges[i][0], &edges[i][1], &edges[i][2])
		edges[i][0]--
		edges[i][1]--
	}

	var q int
	fmt.Fscan(in, &q)
	operations := make([][4]int, q)
	for i := range operations {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var u, v, w int
			fmt.Fscan(in, &u, &v, &w)
			operations[i] = [4]int{op, u - 1, v - 1, w}
		} else {
			var u, v int
			fmt.Fscan(in, &u, &v)
			operations[i] = [4]int{op, u - 1, v - 1, 0}
		}
	}

	R := NewAddRemoveQuery(true)
	seg := NewSegmentTreeDivideAndConquerCopy()
	edgeId := make(map[int]int, n+q)
	for i := range edges {
		u, v := edges[i][0], edges[i][1]
		if u > v {
			u, v = v, u
		}
		edgeId[u*n+v] = i
		R.Add(-1, i) // 初始时刻的边
	}

	queryCount := 0
	for i, query := range operations {
		op, u, v, w := query[0], query[1], query[2], query[3]
		if u > v {
			u, v = v, u
		}
		edgeHash := u*n + v
		if op == 1 {
			edges = append(edges, [3]int{u, v, w})
			newId := len(edges) - 1
			edgeId[edgeHash] = newId
			R.Add(queryCount, newId)
		} else if op == 2 {
			id := edgeId[edgeHash]
			R.Remove(queryCount, id)
		} else {
			seg.AddQuery(queryCount, i)
			queryCount++
		}
	}
	events := R.GetEvents(INF)
	for _, e := range events {
		seg.AddMutation(e.start, e.end, e.value) // value为edges的下标
	}

	res := make([]int, q)
	for i := range res {
		res[i] = INF
	}

	uf := NewUnionFindWithDistAndUndo(n)
	initVectorSpace := NewVectorSpace(nil)
	seg.Run(
		&initVectorSpace,
		func(vectorSpace *State, edgeId int) {
			vs := *vectorSpace
			edge := edges[edgeId]
			u, v, w := edge[0], edge[1], edge[2]
			uf.SnapShot()

			root1, x1 := uf.Find(u)
			root2, x2 := uf.Find(v)
			if root1 != root2 {
				uf.Union(u, v, w)
			} else {
				cycleXor := x1 ^ x2 ^ w
				vs.Add(cycleXor)
			}
		},
		func(vectorSpace *State) *State {
			tmp := *vectorSpace
			res := tmp.Copy()
			return &res
		},
		func(vectorSpace *State, queryId int) {
			vs := *vectorSpace
			query := operations[queryId]
			u, v := query[1], query[2]
			dist := uf.Dist(u, v)
			res[queryId] = vs.Min(dist)
		},
		func() {
			uf.Rollback(-1)
		},
	)

	for _, v := range res {
		if v != INF {
			fmt.Fprintln(out, v)
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
	mp       map[V]int
	events   []Event
	adds     map[V][]int
	removes  map[V][]int
	monotone bool
}

func NewAddRemoveQuery(monotone bool) *AddRemoveQuery {
	return &AddRemoveQuery{
		mp:       map[V]int{},
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
	for value, addTimes := range adq.adds {
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
	if _, ok := adq.mp[value]; ok {
		panic("can't add a value that already exists")
	}
	adq.mp[value] = time
}

func (adq *AddRemoveQuery) removeMonotone(time int, value V) {
	if startTime, ok := adq.mp[value]; !ok {
		panic("can't remove a value that doesn't exist")
	} else {
		delete(adq.mp, value)
		if startTime != time {
			adq.events = append(adq.events, Event{startTime, time, value})
		}
	}
}

func (adq *AddRemoveQuery) getMonotone(lastTime int) []Event {
	for value, startTime := range adq.mp {
		if startTime == lastTime {
			continue
		}
		adq.events = append(adq.events, Event{startTime, lastTime, value})
	}
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
	data      *RollbackArray
	snapShots []int
}

func NewUnionFindWithDistAndUndo(n int) *UnionFindWithDistAndUndo {
	return &UnionFindWithDistAndUndo{
		data: NewRollbackArray(n, func(index int) arrayItem { return arrayItem{parent: -1, dist: e()} }),
	}
}

func (uf *UnionFindWithDistAndUndo) SnapShot() int {
	res := uf.data.GetTime()
	uf.snapShots = append(uf.snapShots, res)
	return res
}

func (uf *UnionFindWithDistAndUndo) GetTime() int {
	return uf.data.GetTime()
}

// time=-1表示回滚到上一个快照的时间，并删除该快照.
func (uf *UnionFindWithDistAndUndo) Rollback(time int) {
	if time != -1 {
		uf.data.Rollback(time)
	} else {
		if len(uf.snapShots) == 0 {
			return
		}
		time = uf.snapShots[len(uf.snapShots)-1]
		uf.snapShots = uf.snapShots[:len(uf.snapShots)-1]
		uf.data.Rollback(time)
	}
}

// distToRoot(parent) + dist = distToRoot(child).
func (uf *UnionFindWithDistAndUndo) Union(parent int, child int, dist T) bool {
	v1, x1 := uf.Find(parent)
	v2, x2 := uf.Find(child)
	if v1 == v2 {
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
