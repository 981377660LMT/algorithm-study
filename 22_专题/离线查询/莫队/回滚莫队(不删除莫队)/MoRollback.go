// 不删除莫队,复杂度和普通莫队一样
// !删除操作很麻烦的时候使用
// 发明者:https://snuke.hatenablog.com/entry/2016/07/01/000000

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
	"sort"
)

func main() {
	StaticRangeInversionsQuery()
	// AT1219()
	// Luogu5906()
}

// Static Range Inversions Query - 静态区间逆序对查询
// https://judge.yosupo.jp/problem/static_range_inversions_query
func StaticRangeInversionsQuery() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	mo := NewMoRollback(n, q)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		mo.AddQuery(l, r)
	}

	// 离散化
	set := make(map[int]struct{})
	for _, v := range nums {
		set[v] = struct{}{}
	}
	allNums := make([]int, 0, len(set))
	for k := range set {
		allNums = append(allNums, k)
	}
	sort.Ints(allNums)
	mp := make(map[int]int, len(allNums))
	for i, v := range allNums {
		mp[v] = i
	}
	newNums := make([]int, n)
	for i, v := range nums {
		newNums[i] = mp[v]
	}

	bit := NewBitArray(len(allNums))
	inv, snap, snapInv := 0, 0, 0 // inv: 当前逆序对数, snap: 当前快照状态, snapInv: 当前快照逆序对数
	history := make([]int, 0, n)  // history: 当前操作历史便于undo
	res := make([]int, q)

	add := func(index, delta int) {
		if delta == 1 { // add_right ->
			x := newNums[index]
			inv += bit.ProdRange(x+1, len(allNums))
			bit.Apply(x, 1)
			history = append(history, x)
		} else { // add_left <-
			x := newNums[index]
			inv += bit.Prod(x)
			bit.Apply(x, 1)
			history = append(history, x)
		}
	}
	_move := func(state int) {
		for len(history) > state {
			x := history[len(history)-1]
			history = history[:len(history)-1]
			bit.Apply(x, -1)
		}
	}
	reset := func() {
		_move(0)
		inv = 0
	}
	snapshot := func() {
		snap = len(history)
		snapInv = inv
	}
	rollback := func() {
		_move(snap)
		inv = snapInv
	}
	query := func(qi int) {
		res[qi] = inv
	}

	mo.Run(add, reset, snapshot, rollback, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// 歴史の研究
// https://www.luogu.com.cn/problem/AT_joisc2014_c
// 给定一个数组nums和q个查询(l,r)
// 每次查询[l,r]区间内的`重要度`,一个数字num的重要度定义为`num乘以区间内num的个数`
func AT1219() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	mo := NewMoRollback(n, q)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		mo.AddQuery(l, r)
	}

	res := make([]int, q)
	cur := 0                     // 当前区间的答案
	snap, snapCur := 0, 0        // 当前快照状态,当前快照答案
	history := make([]int, 0, n) // x
	counter := make(map[int]int)

	add := func(index, _ int) { // TODO
		x := nums[index]
		counter[x] += 1
		cur = max(cur, x*counter[x])
		history = append(history, x)
	}
	_move := func(state int) {
		for len(history) > state {
			x := history[len(history)-1]
			history = history[:len(history)-1]
			counter[x]-- // TODO
		}
	}

	reset := func() {
		_move(0)
		cur = 0
	}
	snapshot := func() {
		snap = len(history)
		snapCur = cur
	}
	rollback := func() {
		_move(snap)
		cur = snapCur
	}
	query := func(qi int) {
		res[qi] = cur
	}

	mo.Run(add, reset, snapshot, rollback, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// https://www.luogu.com.cn/problem/P5906
// 给定一个序列，多次询问一段区间 [l,r]，求区间中相同的数的最远间隔距离。
// 如果区间内不存在两个数相同，则输出 0。
// 序列中两个元素的间隔距离指的是两个元素下标差的绝对值。
//
// !维护每个数在区间内索引的最大值和最小值.
func Luogu5906() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	// 离散化
	_pool := make(map[int]int)
	id := func(o int) int {
		if v, ok := _pool[o]; ok {
			return v
		}
		v := len(_pool)
		_pool[o] = v
		return v
	}
	for i, v := range nums {
		nums[i] = id(v)
	}

	var q int
	fmt.Fscan(in, &q)
	mo := NewMoRollback(n, q)
	for i := 0; i < q; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		l--
		mo.AddQuery(l, r)
	}

	res := make([]int, q)
	cur := 0                        // 当前区间的答案
	snap, snapCur := 0, 0           // 当前快照状态,当前快照答案
	history := make([][3]int, 0, n) // (x,preMinPos,preMaxPos)
	minPos := make([]int, n)
	maxPos := make([]int, n)
	for i := range minPos {
		minPos[i] = n
		maxPos[i] = -1
	}

	add := func(index, _ int) { // TODO
		x := nums[index]
		preMinPos := minPos[x]
		preMaxPos := maxPos[x]
		minPos[x] = min(preMinPos, index)
		maxPos[x] = max(preMaxPos, index)
		cur = max(cur, maxPos[x]-minPos[x])
		history = append(history, [3]int{x, preMinPos, preMaxPos})
	}
	_move := func(state int) {
		for len(history) > state {
			item := history[len(history)-1]
			x, preMinPos, preMaxPos := item[0], item[1], item[2]
			history = history[:len(history)-1]
			minPos[x] = preMinPos
			maxPos[x] = preMaxPos
		}
	}

	reset := func() {
		_move(0)
		cur = 0
	}
	snapshot := func() {
		snap = len(history)
		snapCur = cur
	}
	rollback := func() {
		_move(snap)
		cur = snapCur
	}
	query := func(qi int) {
		if cur > 0 {
			res[qi] = cur
		}
	}

	mo.Run(add, reset, snapshot, rollback, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type MoRollback struct {
	chunkSize          int
	left, right, order []int
}

func NewMoRollback(n, q int) *MoRollback {
	chunkSize := max(1, n/max(1, int(math.Sqrt(float64(q*2/3)))))
	order := make([]int, q)
	for i := range order {
		order[i] = i
	}
	return &MoRollback{chunkSize: chunkSize, order: order}
}

// 添加一个查询，查询范围为`左闭右开区间` [left, right).
//  0 <= left <= right <= n
func (mo *MoRollback) AddQuery(left, right int) {
	mo.left = append(mo.left, left)
	mo.right = append(mo.right, right)
}

// 返回每个查询的结果.
//  add: 将数据添加到窗口.
//  reset: 将窗口重置为初始状态.
//  snapshot: 保存当前窗口的状态.
//  rollback: 恢复窗口的状态.
//  query: 查询窗口的状态.
func (mo *MoRollback) Run(
	add func(index, delta int),
	reset func(),
	snapshot func(),
	rollback func(),
	query func(qi int),
) {
	left, right, order := mo.left, mo.right, mo.order
	chunkSize := mo.chunkSize
	sort.Slice(order, func(i, j int) bool {
		ii, jj := order[i], order[j]
		iblock, jblock := left[ii]/chunkSize, left[jj]/chunkSize
		if iblock != jblock {
			return iblock < jblock
		}
		return right[ii] < right[jj]
	})

	reset()
	for _, idx := range order {
		if right[idx]-left[idx] < chunkSize {
			for i := left[idx]; i < right[idx]; i++ {
				add(i, 1)
			}
			query(idx)
			rollback()
		}
	}

	nr := 0
	lastBlock := -1
	for _, idx := range order {
		if right[idx]-left[idx] < chunkSize {
			continue
		}
		block := left[idx] / chunkSize
		if lastBlock != block {
			reset()
			lastBlock = block
			nr = (block + 1) * chunkSize
		}
		for nr < right[idx] {
			add(nr, 1)
			nr++
		}
		snapshot()
		for j := (block+1)*chunkSize - 1; j >= left[idx]; j-- {
			add(j, -1)
		}
		query(idx)
		rollback()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type BitArray struct {
	n    int
	log  int
	data []int
}

// 長さ n の 0で初期化された配列で構築する.
func NewBitArray(n int) *BitArray {
	return &BitArray{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

// 配列で構築する.
func NewBitArrayFrom(arr []int) *BitArray {
	res := NewBitArray(len(arr))
	res.Build(arr)
	return res
}

func (b *BitArray) Build(arr []int) {
	if b.n != len(arr) {
		panic("len of arr is not equal to n")
	}
	for i := 1; i <= b.n; i++ {
		b.data[i] = arr[i-1]
	}
	for i := 1; i <= b.n; i++ {
		j := i + (i & -i)
		if j <= b.n {
			b.data[j] += b.data[i]
		}
	}
}

// 要素 i に値 v を加える.
func (b *BitArray) Apply(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r) の要素の総和を求める.
func (b *BitArray) Prod(r int) int {
	res := int(0)
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

// [l, r) の要素の総和を求める.
func (b *BitArray) ProdRange(l, r int) int {
	return b.Prod(r) - b.Prod(l)
}

// 区間[0,k]の総和がx以上となる最小のkを求める.数列が単調増加であることを要求する.
func (b *BitArray) LowerBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] < x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

// 区間[0,k]の総和がxを上回る最小のkを求める.数列が単調増加であることを要求する.
func (b *BitArray) UpperBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] <= x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}
