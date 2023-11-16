// Lena and Queries
// https://www.luogu.com.cn/problem/CF678F
// https://codeforces.com/contest/678/submission/182385006
// 给定	q(q<=3e5)个操作，每个操作是以下三种之一：
//
// 1 k b : 将直线 y = kx + b 加入集合.
// 2 pos : 删除第 pos 个加入集合的直线.保证对应的直线存在.
// 3 x : 查询 x 处的最大值.-1e9<=x<=1e9.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var q int
	fmt.Fscan(in, &q)
	operations := make([][3]int, q)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var k, b int
			fmt.Fscan(in, &k, &b)
			operations[i] = [3]int{1, k, b}
		} else if op == 2 {
			var pos int
			fmt.Fscan(in, &pos)
			pos--
			operations[i] = [3]int{2, pos, 0}
		} else {
			var x int
			fmt.Fscan(in, &x)
			operations[i] = [3]int{3, x, 0}
		}
	}

	Q := NewAddRemoveQuery(true)
	seg := NewSegmentTreeDivideAndConquerCopy()
	queries := []int{} // 保存所有查询的x
	for i, operation := range operations {
		kind := operation[0]
		if kind == 1 {
			Q.Add(i, i)
		} else if kind == 2 {
			prePos := operation[1]
			Q.Remove(i, prePos)
		} else {
			x := operation[1]
			queries = append(queries, x)
			seg.AddQuery(i, len(queries)-1)
		}
	}
	mutations := Q.GetEvents(INF)
	for _, m := range mutations {
		seg.AddMutation(m.start, m.end, m.value) // start, end, pos
	}

	res := make([]int, len(queries))
	lichao := NewLiChaoTreeDynamic(-1e9, 1e9, false, true) // 可持久化
	initState := lichao.NewRoot()
	seg.Run(
		&initState,
		func(state *State, pos int) {
			operation := operations[pos]
			k, b := operation[1], operation[2]
			*state = lichao.AddLine(*state, Line{k: k, b: b})
		},
		func(state *State) *State {
			copy_ := lichao.Copy(*state)
			return &copy_
		},
		func(state *State, queryId int) {
			x := queries[queryId]
			res[queryId] = lichao.Query(*state, x).value
		},
	)

	for _, v := range res {
		if v == -INF {
			fmt.Fprintln(out, "EMPTY SET")
		} else {
			fmt.Fprintln(out, v)
		}
	}
}

type State = *LichaoNode // 注意修改state需要使用指针

type segMutation struct{ start, end, id int }
type segQuery struct{ time, id int }

// 线段树分治copy版.
// 如果修改操作难以撤销，可以在每个节点处保存一份副本.
// !调用O(n)次拷贝注意不要超出内存.
type SegmentTreeDivideAndConquerCopy struct {
	initState *State
	mutate    func(state *State, mutationId int)
	copy      func(state *State) *State
	query     func(state *State, queryId int)
	mutations []segMutation
	queries   []segQuery
	nodes     [][]int
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
	o.queries = append(o.queries, segQuery{time, id})
}

// dfs 遍历整棵线段树来得到每个时间点的答案.
//
//	 initState: 数据结构的初始状态.
//		mutate: 添加编号为`mutationId`的变更后产生的副作用.
//		copy: 拷贝一份数据结构的副本.
//		query: 响应编号为`queryId`的查询.
//		一共调用 **O(nlogn)** 次`mutate`，**O(n)** 次`copy` 和 **O(q)** 次`query`.
func (o *SegmentTreeDivideAndConquerCopy) Run(
	initState *State,
	mutate func(state *State, mutationId int),
	copy func(state *State) *State,
	query func(state *State, queryId int),
) {
	if len(o.queries) == 0 {
		return
	}
	o.initState = initState
	o.mutate, o.copy, o.query = mutate, copy, query
	times := make([]int, len(o.queries))
	for i := range o.queries {
		times[i] = o.queries[i].time
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

	for _, q := range o.queries {
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

type T = int

const INF T = 2e18

type Line struct{ k, b T } // y = k * x + b

// Evaluate を書き変えると、totally monotone な関数群にも使える
func Evaluate(line Line, x int) T {
	return line.k*x + line.b
}

type LichaoNode struct {
	lineId      int
	left, right *LichaoNode
}

type queryPair = struct {
	value  T
	lineId int
}

// 可持久化李超线段树.注意`添加线段`时空间消耗较大.
type LiChaoTreeDynamic struct {
	start, end int
	minimize   bool
	persistent bool
	lines      []Line
}

func NewLiChaoTreeDynamic(start, end int, minimize bool, persistent bool) *LiChaoTreeDynamic {
	end++
	return &LiChaoTreeDynamic{
		start: start, end: end,
		minimize:   minimize,
		persistent: persistent,
	}
}

func (tree *LiChaoTreeDynamic) NewRoot() *LichaoNode {
	return nil
}

// O(logn)
func (tree *LiChaoTreeDynamic) AddLine(root *LichaoNode, line Line) *LichaoNode {
	id := len(tree.lines)
	tree.lines = append(tree.lines, line)
	if root == nil {
		root = &LichaoNode{lineId: -1}
	}
	return tree._addLine(root, id, tree.start, tree.end)
}

// [start, end)
// O(log^2n)
func (tree *LiChaoTreeDynamic) AddSegment(root *LichaoNode, startX, endX int, line Line) *LichaoNode {
	if startX >= endX {
		return root
	}
	id := len(tree.lines)
	tree.lines = append(tree.lines, line)
	if root == nil {
		root = &LichaoNode{lineId: -1}
	}
	return tree._addSegment(root, startX, endX, id, tree.start, tree.end)
}

// O(logn)
func (tree *LiChaoTreeDynamic) Query(root *LichaoNode, x int) queryPair {
	if !(tree.start <= x && x < tree.end) {
		panic("x is out of range")
	}
	if root == nil {
		if tree.minimize {
			return queryPair{value: INF, lineId: -1}
		}
		return queryPair{value: -INF, lineId: -1}
	}
	return tree._query(root, x, tree.start, tree.end)
}

func (tree *LiChaoTreeDynamic) Clear() {
	tree.lines = tree.lines[:0]
}

func (tree *LiChaoTreeDynamic) Copy(node *LichaoNode) *LichaoNode {
	if node == nil || !tree.persistent {
		return node
	}
	return &LichaoNode{lineId: node.lineId, left: node.left, right: node.right}
}

func (tree *LiChaoTreeDynamic) _evaluateInner(fid int, x int) T {
	if fid == -1 {
		if tree.minimize {
			return INF
		}
		return -INF
	}
	return Evaluate(tree.lines[fid], x)
}

func (tree *LiChaoTreeDynamic) _addLine(node *LichaoNode, fid int, nodeL, nodeR int) *LichaoNode {
	gid := node.lineId
	fl := tree._evaluateInner(fid, nodeL)
	fr := tree._evaluateInner(fid, nodeR-1)
	gl := tree._evaluateInner(gid, nodeL)
	gr := tree._evaluateInner(gid, nodeR-1)
	var bl, br bool
	if tree.minimize {
		bl = fl < gl
		br = fr < gr
	} else {
		bl = fl > gl
		br = fr > gr
	}
	if bl && br {
		node = tree.Copy(node)
		node.lineId = fid
		return node
	}
	if !bl && !br {
		return node
	}
	node = tree.Copy(node)
	nodeM := (nodeL + nodeR) >> 1
	fm := tree._evaluateInner(fid, nodeM)
	gm := tree._evaluateInner(gid, nodeM)
	var bm bool
	if tree.minimize {
		bm = fm < gm
	} else {
		bm = fm > gm
	}
	if bm {
		node.lineId = fid
		if bl {
			if node.right == nil {
				node.right = &LichaoNode{lineId: -1}
			}
			node.right = tree._addLine(node.right, gid, nodeM, nodeR)
		} else {
			if node.left == nil {
				node.left = &LichaoNode{lineId: -1}
			}
			node.left = tree._addLine(node.left, gid, nodeL, nodeM)
		}
	} else {
		if !bl {
			if node.right == nil {
				node.right = &LichaoNode{lineId: -1}
			}
			node.right = tree._addLine(node.right, fid, nodeM, nodeR)
		} else {
			if node.left == nil {
				node.left = &LichaoNode{lineId: -1}
			}
			node.left = tree._addLine(node.left, fid, nodeL, nodeM)
		}
	}
	return node
}

func (tree *LiChaoTreeDynamic) _addSegment(node *LichaoNode, xl, xr int, fid int, nodeL, nodeR int) *LichaoNode {
	if nodeL > xl {
		xl = nodeL
	}
	if nodeR < xr {
		xr = nodeR
	}
	if xl >= xr {
		return node
	}
	if nodeL < xl || xr < nodeR {
		node = tree.Copy(node)
		nodeM := (nodeL + nodeR) >> 1
		if node.left == nil {
			node.left = &LichaoNode{lineId: -1}
		}
		if node.right == nil {
			node.right = &LichaoNode{lineId: -1}
		}
		node.left = tree._addSegment(node.left, xl, xr, fid, nodeL, nodeM)
		node.right = tree._addSegment(node.right, xl, xr, fid, nodeM, nodeR)
		return node
	}
	return tree._addLine(node, fid, nodeL, nodeR)
}

func (tree *LiChaoTreeDynamic) _query(node *LichaoNode, x int, nodeL, nodeR int) queryPair {
	fid := node.lineId
	res := queryPair{value: tree._evaluateInner(fid, x), lineId: fid}
	nodeM := (nodeL + nodeR) >> 1
	if x < nodeM && node.left != nil {
		cand := tree._query(node.left, x, nodeL, nodeM)
		if tree.minimize {
			if cand.value < res.value {
				res = cand
			}
		} else {
			if cand.value > res.value {
				res = cand
			}
		}
	}
	if x >= nodeM && node.right != nil {
		cand := tree._query(node.right, x, nodeM, nodeR)
		if tree.minimize {
			if cand.value < res.value {
				res = cand
			}
		} else {
			if cand.value > res.value {
				res = cand
			}
		}
	}
	return res
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
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
