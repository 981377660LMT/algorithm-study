// https://ei1333.github.io/library/other/mo-tree.hpp
// https://oi-wiki.org/misc/mo-algo-on-tree/
// https://github.com/EndlessCheng/codeforces-go/blob/53262fb81ffea176cd5f039cec71e3bd266dce83/copypasta/mo.go#L301
// 处理树上的路径相关的离线查询.
// 一般的莫队只能处理线性问题，我们要把树强行压成序列。
// 通过欧拉序(括号序)转化成序列上的查询，然后用莫队解决。

package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

func main() {

	n := 8
	edges := [][]int{{1, 2, 6}, {1, 3, 4}, {2, 4, 6}, {2, 5, 3}, {3, 6, 6}, {3, 0, 8}, {7, 0, 2}}
	queries := [][]int{{4, 6}, {0, 4}, {6, 5}, {7, 4}}
	fmt.Println(minOperationsQueries(n, edges, queries))

}

// 100018. 边权重均等查询
// https://leetcode.cn/problems/minimum-edge-weight-equilibrium-queries-in-a-tree/description/
func minOperationsQueries(n int, edges [][]int, queries [][]int) []int {
	tree := make([][][2]int, n)
	for eid, edge := range edges {
		u, v := edge[0], edge[1]
		tree[u] = append(tree[u], [2]int{v, eid})
		tree[v] = append(tree[v], [2]int{u, eid})
	}

	mo := NewMoOnTreeEdge(tree, 0)
	for _, q := range queries {
		mo.AddQuery(q[0], q[1])
	}

	res := make([]int, len(queries))
	weightCounter := make(map[int]int) // 	weightCounter := [30]int{}
	sl := NewSortedListRangeBlock(n - 1)
	edgeCount := 0

	add := func(edgeId int) {
		weight := edges[edgeId][2]
		preCount := weightCounter[weight]
		weightCounter[weight] = preCount + 1
		sl.Discard(preCount)
		sl.Add(preCount + 1)
		edgeCount++
	}
	remove := func(edgeId int) {
		weight := edges[edgeId][2]
		preCount := weightCounter[weight]
		weightCounter[weight] = preCount - 1
		sl.Discard(preCount)
		if preCount > 1 {
			sl.Add(preCount - 1)
		}
		edgeCount--
	}
	query := func(qid int) {
		if sl.Len() > 0 {
			max_ := sl.Max()
			res[qid] = edgeCount - max_
		}
	}

	mo.Run(add, remove, query)
	return res
}

// 维护边权的树上莫队.
type MoOnTreeEdge struct {
	tree    [][][2]int // (next,eid)
	root    int
	queries [][2]int
}

func NewMoOnTreeEdge(tree [][][2]int, root int) *MoOnTreeEdge {
	return &MoOnTreeEdge{tree: tree, root: root}
}

// 添加从顶点from到顶点to的查询.
func (mo *MoOnTreeEdge) AddQuery(from, to int) {
	mo.queries = append(mo.queries, [2]int{from, to})
}

// 处理每个查询.
//
//	add: 将边添加到窗口.
//	remove: 将边从窗口移除.
//	query: 查询窗口内的数据.
func (mo *MoOnTreeEdge) Run(
	add func(edgeId int),
	remove func(edgeId int),
	query func(qid int),
) {
	if len(mo.queries) == 0 {
		return
	}

	n := len(mo.tree)
	dfn := 0
	dfnToEdge := make([]int, 2*n)
	ins := make([]int, n)
	outs := make([]int, n)

	var dfs func(cur, pre int)
	dfs = func(cur, pre int) {
		ins[cur] = dfn
		for _, e := range mo.tree[cur] {
			to, eid := e[0], e[1]
			if to != pre {
				dfnToEdge[dfn] = eid
				dfn++
				dfs(to, cur)
				dfnToEdge[dfn] = eid
				dfn++
			}
		}
		outs[cur] = dfn
	}
	dfs(mo.root, -1)

	lca := _offlineLCA(mo.tree, mo.queries, mo.root)
	blockSize := int(math.Ceil(float64(2*n) / math.Sqrt(float64(len(mo.queries)))))
	type Q struct{ bid, l, r, qid int }
	qs := make([]Q, len(mo.queries))
	for qi := range qs {
		v, w := mo.queries[qi][0], mo.queries[qi][1]
		if ins[v] > ins[w] {
			v, w = w, v
		}
		if lca_ := lca[qi]; lca_ != v {
			qs[qi] = Q{outs[v] / blockSize, outs[v], ins[w], qi}
		} else {
			qs[qi] = Q{ins[v] / blockSize, ins[v], ins[w], qi}
		}
	}

	sort.Slice(qs, func(i, j int) bool {
		a, b := qs[i], qs[j]
		if a.bid != b.bid {
			return a.bid < b.bid
		}
		if a.bid&1 == 0 {
			return a.r < b.r
		}
		return a.r > b.r
	})

	flip := make([]bool, n)
	f := func(u int) {
		flip[u] = !flip[u]
		if flip[u] {
			add(u)
		} else {
			remove(u)
		}
	}

	l, r := 0, 0
	for _, q := range qs {
		for ; r < q.r; r++ {
			f(dfnToEdge[r])
		}
		for ; l < q.l; l++ {
			f(dfnToEdge[l])
		}
		for l > q.l {
			l--
			f(dfnToEdge[l])
		}
		for r > q.r {
			r--
			f(dfnToEdge[r])
		}
		query(q.qid)
	}
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

// LCA离线.
func _offlineLCA(tree [][][2]int, queries [][2]int, root int) []int {
	n := len(tree)
	ufa := NewUnionFindArray(n)
	st, mark, ptr, res := make([]int, n), make([]int, n), make([]int, n), make([]int, len(queries))
	for i := 0; i < len(queries); i++ {
		res[i] = -1
	}
	top := 0
	st[top] = root
	for _, q := range queries {
		mark[q[0]]++
		mark[q[1]]++
	}
	q := make([][][2]int, n)
	for i := 0; i < n; i++ {
		q[i] = make([][2]int, 0, mark[i])
		mark[i] = -1
		ptr[i] = len(tree[i])
	}
	for i := range queries {
		u, v := queries[i][0], queries[i][1]
		q[u] = append(q[u], [2]int{v, i})
		q[v] = append(q[v], [2]int{u, i})
	}
	run := func(u int) bool {
		nexts := tree[u]
		for ptr[u] != 0 {
			v := nexts[ptr[u]-1][0]
			ptr[u]--
			if mark[v] == -1 {
				top++
				st[top] = v
				return true
			}
		}
		return false
	}

	for top != -1 {
		u := st[top]
		nexts := tree[u]
		if mark[u] == -1 {
			mark[u] = u
		} else {
			ufa.Union(u, nexts[ptr[u]][0])
			mark[ufa.Find(u)] = u
		}

		if !run(u) {
			for _, v := range q[u] {
				if mark[v[0]] != -1 && res[v[1]] == -1 {
					res[v[1]] = mark[ufa.Find(v[0])]
				}
			}
			top--
		}
	}

	return res
}

type _unionFindArray struct {
	data []int
}

func NewUnionFindArray(n int) *_unionFindArray {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &_unionFindArray{data: data}
}

func (ufa *_unionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	return true
}

func (ufa *_unionFindArray) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

const INF int = 1e18

type SortedListRangeBlock struct {
	_blockSize  int   // 每个块的大小.
	_counter    []int // 每个数出现的次数.
	_blockCount []int // 每个块的数的个数.
	_blockSum   []int // 每个块的和.
	_belong     []int // 每个数所在的块.
	_len        int   // 所有数的个数.
}

// 值域分块模拟SortedList.
// `O(1)`add/remove，`O(sqrt(n))`查询.
// 一般配合莫队算法使用.
//
//	max:值域的最大值.0 <= max <= 1e6.
//	iterable:初始值.
func NewSortedListRangeBlock(max int, nums ...int) *SortedListRangeBlock {
	max += 5
	size := int(math.Sqrt(float64(max)))
	count := 1 + (max / size)
	sl := &SortedListRangeBlock{
		_blockSize:  size,
		_counter:    make([]int, max),
		_blockCount: make([]int, count),
		_blockSum:   make([]int, count),
		_belong:     make([]int, max),
	}
	for i := 0; i < max; i++ {
		sl._belong[i] = i / size
	}
	if len(nums) > 0 {
		sl.Update(nums...)
	}
	return sl
}

// O(1).
func (sl *SortedListRangeBlock) Add(value int) {
	sl._counter[value]++
	pos := sl._belong[value]
	sl._blockCount[pos]++
	sl._blockSum[pos] += value
	sl._len++
}

// O(1).
func (sl *SortedListRangeBlock) Remove(value int) {
	sl._counter[value]--
	pos := sl._belong[value]
	sl._blockCount[pos]--
	sl._blockSum[pos] -= value
	sl._len--
}

// O(1).
func (sl *SortedListRangeBlock) Discard(value int) bool {
	if !sl.Has(value) {
		return false
	}
	sl.Remove(value)
	return true
}

// O(1).
func (sl *SortedListRangeBlock) Has(value int) bool {
	return sl._counter[value] > 0
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) At(index int) int {
	if index < 0 {
		index += sl._len
	}
	if index < 0 || index >= sl._len {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	for i := 0; i < len(sl._blockCount); i++ {
		count := sl._blockCount[i]
		if index < count {
			num := i * sl._blockSize
			for {
				numCount := sl._counter[num]
				if index < numCount {
					return num
				}
				index -= numCount
				num++
			}
		}
		index -= count
	}
	panic("unreachable")
}

// 严格小于 value 的元素个数.
// 也即第一个大于等于 value 的元素的下标.
// O(sqrt(n)).
func (sl *SortedListRangeBlock) BisectLeft(value int) int {
	pos := sl._belong[value]
	res := 0
	for i := 0; i < pos; i++ {
		res += sl._blockCount[i]
	}
	for v := pos * sl._blockSize; v < value; v++ {
		res += sl._counter[v]
	}
	return res
}

// 小于等于 value 的元素个数.
// 也即第一个大于 value 的元素的下标.
// O(sqrt(n)).
func (sl *SortedListRangeBlock) BisectRight(value int) int {
	return sl.BisectLeft(value + 1)
}

func (sl *SortedListRangeBlock) Count(value int) int {
	return sl._counter[value]
}

// 返回范围 `[min, max]` 内数的个数.
// O(sqrt(n)).
func (sl *SortedListRangeBlock) CountRange(min, max int) int {
	if min > max {
		return 0
	}

	minPos := sl._belong[min]
	maxPos := sl._belong[max]
	if minPos == maxPos {
		res := 0
		for i := min; i <= max; i++ {
			res += sl._counter[i]
		}
		return res
	}

	res := 0
	minEnd := (minPos + 1) * sl._blockSize
	for v := min; v < minEnd; v++ {
		res += sl._counter[v]
	}
	for i := minPos + 1; i < maxPos; i++ {
		res += sl._blockCount[i]
	}
	maxStart := maxPos * sl._blockSize
	for v := maxStart; v <= max; v++ {
		res += sl._counter[v]
	}
	return res
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) Lower(value int) (res int, ok bool) {
	pos := sl._belong[value]
	start := pos * sl._blockSize
	for v := value - 1; v >= start; v-- {
		if sl._counter[v] > 0 {
			return v, true
		}
	}

	for i := pos - 1; i >= 0; i-- {
		if sl._blockCount[i] == 0 {
			continue
		}
		num := (i + 1) * sl._blockSize
		for {
			if sl._counter[num] > 0 {
				return num, true
			}
			num--
		}
	}

	return
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) Higher(value int) (res int, ok bool) {
	pos := sl._belong[value]
	end := (pos + 1) * sl._blockSize
	for v := value + 1; v < end; v++ {
		if sl._counter[v] > 0 {
			return v, true
		}
	}

	for i := pos + 1; i < len(sl._blockCount); i++ {
		if sl._blockCount[i] == 0 {
			continue
		}
		num := i * sl._blockSize
		for {
			if sl._counter[num] > 0 {
				return num, true
			}
			num++
		}
	}

	return
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) Floor(value int) (res int, ok bool) {
	if sl.Has(value) {
		return value, true
	}
	return sl.Lower(value)
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) Ceiling(value int) (res int, ok bool) {
	if sl.Has(value) {
		return value, true
	}
	return sl.Higher(value)
}

// 返回区间 `[start, end)` 的和.
// O(sqrt(n)).
func (sl *SortedListRangeBlock) SumSlice(start, end int) int {
	if start < 0 {
		start += sl._len
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += sl._len
	}
	if end > sl._len {
		end = sl._len
	}
	if start >= end {
		return 0
	}

	res := 0
	remain := end - start
	cur, index := sl._findKth(start)
	sufCount := sl._counter[cur] - index
	if sufCount >= remain {
		return remain * cur
	}

	res += sufCount * cur
	remain -= sufCount
	cur++

	// 当前块内的和
	blockEnd := (sl._belong[cur] + 1) * sl._blockSize
	for remain > 0 && cur < blockEnd {
		count := sl._counter[cur]
		real := count
		if real > remain {
			real = remain
		}
		res += real * cur
		remain -= real
		cur++
	}

	// 以块为单位消耗remain
	pos := sl._belong[cur]
	for pos < len(sl._blockCount) && remain >= sl._blockCount[pos] {
		res += sl._blockSum[pos]
		remain -= sl._blockCount[pos]
		pos++
		cur += sl._blockSize
	}

	// 剩余的
	for remain > 0 {
		count := sl._counter[cur]
		real := count
		if real > remain {
			real = remain
		}
		res += real * cur
		remain -= real
		cur++
	}

	return res
}

// 返回范围 `[min, max]` 的和.
// O(sqrt(n)).
func (sl *SortedListRangeBlock) SumRange(min, max int) int {
	minPos := sl._belong[min]
	maxPos := sl._belong[max]
	if minPos == maxPos {
		res := 0
		for i := min; i <= max; i++ {
			res += sl._counter[i] * i
		}
		return res
	}

	res := 0
	minEnd := (minPos + 1) * sl._blockSize
	for v := min; v < minEnd; v++ {
		res += sl._counter[v] * v
	}
	for i := minPos + 1; i < maxPos; i++ {
		res += sl._blockSum[i]
	}
	maxStart := maxPos * sl._blockSize
	for v := maxStart; v <= max; v++ {
		res += sl._counter[v] * v
	}
	return res
}

func (sl *SortedListRangeBlock) ForEach(f func(value, index int), reverse bool) {
	if reverse {
		ptr := 0
		for i := len(sl._counter) - 1; i >= 0; i-- {
			count := sl._counter[i]
			for j := 0; j < count; j++ {
				f(i, ptr)
				ptr++
			}
		}
	} else {
		ptr := 0
		for i := 0; i < len(sl._counter); i++ {
			count := sl._counter[i]
			for j := 0; j < count; j++ {
				f(i, ptr)
				ptr++
			}
		}
	}
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) Pop(index int) int {
	if index < 0 {
		index += sl._len
	}
	if index < 0 || index >= sl._len {
		panic(fmt.Sprintf("index out of range: %d", index))
	}
	value := sl.At(index)
	sl.Remove(value)
	return value
}

func (sl *SortedListRangeBlock) Slice(start, end int) []int {
	if start < 0 {
		start += sl._len
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end += sl._len
	}
	if end > sl._len {
		end = sl._len
	}
	if start >= end {
		return nil
	}

	res := make([]int, end-start)
	count := 0
	sl.Enumerate(start, end, func(value int) {
		res[count] = value
		count++
	}, false)

	return res
}

// O(sqrt(n)).
func (sl *SortedListRangeBlock) Erase(start, end int) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SortedListRangeBlock) Enumerate(start, end int, f func(value int), erase bool) {
	if start < 0 {
		start = 0
	}
	if end > sl._len {
		end = sl._len
	}
	if start >= end {
		return
	}

	remain := end - start
	cur, index := sl._findKth(start)
	sufCount := sl._counter[cur] - index
	real := sufCount
	if real > remain {
		real = remain
	}
	if f != nil {
		for i := 0; i < real; i++ {
			f(cur)
		}
	}
	if erase {
		for i := 0; i < real; i++ {
			sl.Remove(cur)
		}
	}
	remain -= sufCount
	if remain == 0 {
		return
	}
	cur++

	// 当前块内
	blockEnd := (sl._belong[cur] + 1) * sl._blockSize
	for remain > 0 && cur < blockEnd {
		count := sl._counter[cur]
		real := count
		if real > remain {
			real = remain
		}
		remain -= real
		if f != nil {
			for i := 0; i < real; i++ {
				f(cur)
			}
		}
		if erase {
			for i := 0; i < real; i++ {
				sl.Remove(cur)
			}
		}
		cur++
	}

	// 以块为单位消耗remain
	pos := sl._belong[cur]
	for count := sl._blockCount[pos]; remain >= count; {
		remain -= count
		if f != nil {
			for v := cur; v < cur+sl._blockSize; v++ {
				c := sl._counter[v]
				for i := 0; i < c; i++ {
					f(v)
				}
			}
		}
		if erase {
			for v := cur; v < cur+sl._blockSize; v++ {
				sl._counter[v] = 0
			}
			sl._len -= count
			sl._blockCount[pos] = 0
			sl._blockSum[pos] = 0
		}
		pos++
		cur += sl._blockSize
	}

	// 剩余的
	for remain > 0 {
		count := sl._counter[cur]
		real := count
		if real > remain {
			real = remain
		}
		remain -= real
		if f != nil {
			for i := 0; i < real; i++ {
				f(cur)
			}
		}
		if erase {
			for i := 0; i < real; i++ {
				sl.Remove(cur)
			}
		}
		cur++
	}
}

func (sl *SortedListRangeBlock) Clear() {
	for i := range sl._counter {
		sl._counter[i] = 0
	}
	for i := range sl._blockCount {
		sl._blockCount[i] = 0
	}
	for i := range sl._blockSum {
		sl._blockSum[i] = 0
	}
	sl._len = 0
}
func (sl *SortedListRangeBlock) Update(values ...int) {
	for _, value := range values {
		sl.Add(value)
	}
}

func (sl *SortedListRangeBlock) Merge(other *SortedListRangeBlock) {
	other.ForEach(func(value, _ int) {
		sl.Add(value)
	}, false)
}

func (sl *SortedListRangeBlock) String() string {
	sb := make([]string, 0, sl._len)
	sl.ForEach(func(value, _ int) {
		sb = append(sb, fmt.Sprintf("%d", value))
	}, false)
	return fmt.Sprintf("SortedListRangeBlock{%s}", strings.Join(sb, ", "))
}

func (sl *SortedListRangeBlock) Len() int {
	return sl._len
}

func (sl *SortedListRangeBlock) Min() int {
	return sl.At(0)
}

func (sl *SortedListRangeBlock) Max() int {
	if sl._len == 0 {
		panic("empty")
	}

	for i := len(sl._blockCount) - 1; i >= 0; i-- {
		if sl._blockCount[i] == 0 {
			continue
		}
		num := (i+1)*sl._blockSize - 1
		for {
			if sl._counter[num] > 0 {
				return num
			}
			num--
		}
	}

	panic("unreachable")
}

// 返回索引在`kth`处的元素的`value`,以及该元素是`value`中的第几个(`index`).
func (sl *SortedListRangeBlock) _findKth(kth int) (value, index int) {
	for i := 0; i < len(sl._blockCount); i++ {
		count := sl._blockCount[i]
		if kth < count {
			num := i * sl._blockSize
			for {
				numCount := sl._counter[num]
				if kth < numCount {
					return num, kth
				}
				kth -= numCount
				num++
			}
		}
		kth -= count
	}

	panic("unreachable")
}
